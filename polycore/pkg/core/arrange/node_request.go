// api request

package arrange

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/auth"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

// SwagConstValue is the predefined value
type SwagConstValue = exprx.SwagConstValue

// RequestNodeDetail represents detail of an request node
type RequestNodeDetail struct {
	RawPath    string         `json:"rawPath"`    // Id of the request raw API
	APIKeyID   string         `json:"apiKeyID"`   // fixed api key id
	DynamicKey bool           `json:"dynamicKey"` // use dynamic key for requester
	Inputs     exprx.ValueSet `json:"inputs"`     // input from header, path, body or uri(GET)
}

// TypeName returns name of the type
func (d RequestNodeDetail) TypeName() string { return NodeTypeRequest.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (d *RequestNodeDetail) DelayedJSONDecode() error {
	if err := d.Inputs.DelayedJSONDecode(); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------
func toJsParas(paras []string) string {
	var buf = bytes.NewBuffer(make([]byte, 0, exprx.DefaultBufLen))
	for _, v := range paras {
		buf.WriteString(",")
		buf.WriteString(v)
	}
	return buf.String()
}

func toJsFmtScript(_fmt string, paras []string) string {
	return fmt.Sprintf(`format("%s" %s)`, _fmt, toJsParas(paras))
}

func (d *RequestNodeDetail) getEncodingParaFunc(encoding string, method exprx.Enum) string {
	enc := exprx.Enum(encoding)
	switch {
	case method == exprx.MethodGet || enc == exprx.EncodingJSON:
		return consts.PDToJSON
	case enc == exprx.EncodingXML:
		return consts.PDToXML
	case enc == exprx.EncodingYAML:
		return consts.PDToYAML
	}
	return ""
}

func (d *RequestNodeDetail) buildBodyScript(depth int, b *jsBuilder) (string, error) {
	s, err := d.Inputs.ToScript(depth+1, b.evaler)
	if err != nil {
		return "", err
	}
	return b.resolveEmptyNodeFieldRefer(s, b.evaler)
}

// BuildJsScript build and write script text into buffer
func (d *RequestNodeDetail) BuildJsScript(n *Node, b *jsBuilder, depth int) error {
	if d.RawPath == "" {
		return errcode.ErrBuildRequestNodeMissingRawAPI.NewError()
	}

	if err := d.checkRequiredField(n, b, depth); err != nil {
		return err
	}

	raw, err := adaptor.GetRawAPIOper().Query(context.TODO(), &adaptor.QueryRawAPIReq{APIPath: d.RawPath})
	if err != nil {
		return err
	}

	// permission verification
	if err := auth.ValidateAPIKeyPermission(d.APIKeyID, raw.Service, b.owner); err != nil {
		return err
	}

	if err := d.handleInputsFromRaw(raw, n); err != nil {
		return err
	}

	b.putRawAPI(d.RawPath) // record the raw apipath

	// repeace path parameter
	const apiPathVarName = "_apiPath"
	reqURL, paras, err := d.Inputs.ResolvePathArgs(raw.URL, n.Name)
	if err != nil {
		return fmt.Errorf("raw api [%s] path parameter error: %s", d.RawPath, err.Error())
	}

	lineHead := exprx.GenLinehead(depth)
	b.buf.WriteString(fmt.Sprintf("%sif (true) { // %s, %s\n", lineHead, n.Name, n.Desc))
	b.buf.WriteString(fmt.Sprintf("%s  var %s = %s\n", lineHead, apiPathVarName, toJsFmtScript(reqURL, paras)))
	body, err := d.buildBodyScript(depth+1, b)
	if err != nil {
		return errcode.ErrBuildNodeFail.FmtError(n.getNameAlias(), err.Error())
	}
	b.buf.WriteString(fmt.Sprintf(`%s  var %s = %s`,
		lineHead, polyRequestTempVar, body))
	b.writeLn()

	const tmpHeaderVarName = "_th"
	if err := d.buildHeaderScript(tmpHeaderVarName, n, b, depth); err != nil {
		return err
	}

	const tmpBodyName = "_tb"
	const tmpKeyName = "_tk"
	var keyID = "''"
	if d.APIKeyID == "" {
		if useDynamicKey := false; !useDynamicKey { // TODO: useDynamicKey := d.DynamicKey
			id, err := auth.QueryGrantedAPIKeyWithError(b.owner, raw.Service, raw.AuthType, "")
			if err != nil {
				return err
			}
			keyID = fmt.Sprintf("'%s'", id)
		} else { // dynamic key script
			keyID = fmt.Sprintf("%s(%s(%v), '%s', '%s', '')", consts.PDQueryGrantedAPIKey,
				consts.PDQueryUser, true, raw.Service, raw.AuthType)
		}
	} else {
		keyID = fmt.Sprintf("'%s'", d.APIKeyID)
	}

	b.buf.WriteString(fmt.Sprintf(`%s  var %s = %s;`,
		lineHead, tmpKeyName, keyID))
	b.writeLn()

	method := exprx.Enum(raw.Content.Method)
	authType := auth.GetAsAuthType(raw.AuthType).String()
	//  var _tb = pdAppendAuth('', 'none', _th, pdToJson(_t))
	b.buf.WriteString(fmt.Sprintf(`%s  var %s = %s(%s, '%s', %s, %s(%s))`,
		lineHead, tmpBodyName, consts.PDAppendAuth, tmpKeyName, authType, tmpHeaderVarName,
		d.getEncodingParaFunc(raw.Content.EncodingIn, method), polyRequestTempVar))
	b.writeLn()

	switch exprx.Enum(raw.Schema) {
	case exprx.SchemaHTTP, exprx.SchemaHTTPS:
		// d.req1 = pdToJsobj("json", pdHttpRequest(_apiPath, "POST", pdToJson(_t), _th, pdQueryUser(true)))
		b.buf.WriteString(fmt.Sprintf(`%s  %s = %s("%s", %s(%s, "%s", %s, %s, %s(%v)))`,
			lineHead, exprx.FullVarName(n.Name), consts.PDToJsobj,
			raw.Content.EncodingOut, consts.PDHttpRequest, apiPathVarName,
			raw.Content.Method, tmpBodyName, tmpHeaderVarName, consts.PDQueryUser, true))

	default:
		return fmt.Errorf("%s: unsupport schema(%s)", raw.URL, raw.Schema)
	}

	b.writeLn()
	b.buf.WriteString(fmt.Sprintf("%s}\n", lineHead))
	return nil
}

func (d *RequestNodeDetail) checkRequiredField(n *Node, b *jsBuilder, depth int) error {
	doc, err := adaptor.GetRawAPIOper().QueryDoc(context.TODO(),
		&apiprovider.QueryDocReq{
			APIPath: d.RawPath,
			DocType: apiprovider.DocTypeRaw.String(),
		})
	if err != nil {
		return err
	}

	var docInput InputNodeDetail
	var apiDoc = apiprovider.APIDoc{
		Input:  &docInput,
		Output: new(json.RawMessage),
	}
	if err := json.Unmarshal(doc.Doc, &apiDoc); err != nil {
		return err
	}
	if err := docInput.checkRequired(d.Inputs); err != nil {
		return err
	}

	return nil
}

func (d *RequestNodeDetail) buildHeaderScript(varName string, n *Node, b *jsBuilder, depth int) error {
	lineHead := exprx.GenLinehead(depth)
	b.buf.WriteString(fmt.Sprintf(`%s  var %s = %s()`,
		lineHead, varName, consts.PDNewHTTPHeader))
	b.writeLn()
	for i := 0; i < len(d.Inputs); i++ {
		if p := (*exprx.XInputValue)(&d.Inputs[i]); p.In.IsHeader() /*&& p.Type.IsHeaderAcceptable()*/ {
			s, err := p.ToScript(depth, b.evaler)
			if err != nil {
				return err
			}
			b.buf.WriteString(fmt.Sprintf(`%s  %s(%s, "%s", %s)`,
				lineHead, consts.PDAddHTTPHeader, varName, p.GetName(false), s))
			b.writeLn()
		}
	}
	b.buf.WriteString(fmt.Sprintf(`%s  %s(%s, "%s")`,
		lineHead, consts.PDUpdateReferPath, varName, d.RawPath))
	b.writeLn()
	b.writeLn()
	return nil
}

func (d *RequestNodeDetail) handleInputsFromRaw(raw *adaptor.QueryRawAPIResp, n *Node) error {
	if err := d.checkRawAPI(raw, n); err != nil {
		return err
	}

	var action *expr.SwagConstValue
	inputs := raw.Content.Consts
	for i := 0; i < len(inputs); i++ {
		p := &inputs[i]
		switch t := exprx.Enum(p.Type); {
		case t.IsAction():
			if action != nil {
				return fmt.Errorf("dunplicate action parmeter %s and %s",
					action.Name, p.Name)
			}
			action = p
		case t.IsPredefineable():
			switch in := exprx.Enum(p.In); in {
			case "", exprx.ParaTypeBody, exprx.ParaTypeHeader, exprx.ParaTypePath, exprx.ParaTypeQuery:
				if err := d.Inputs.AddKV(p.Name, p.GetAsString(), t, in); err != nil {
					return err
				}
			}
		}
	}

	enc := exprx.Enum(raw.Content.EncodingIn)

	// BUG: https://home.yunify.com/distributor.action?serviceName=clogin
	// response 500 with "Content-Type" header
	if !httputil.IsQueryMethod(raw.Method) {
		contenType, err := enc.EncodingToMIME()
		if err != nil {
			fmt.Errorf("rawapi error: %s, %s", raw.ID, err.Error())
		}
		err = d.Inputs.AddKV(consts.HeaderContentType, contenType, exprx.ValTypeString, exprx.ParaTypeHeader)
		if err != nil {
			return err
		}
	}

	if raw.Content.Action != "" {
		if action == nil {
			return fmt.Errorf("missing action parameter for node[%s]", n.Name)
		}

		// add a string KV
		p := action
		err := d.Inputs.AddKV(p.Name, raw.Content.Action, exprx.ValTypeString, exprx.Enum(p.In))
		if err != nil {
			return err
		}
	}

	return nil
}

// checkRawApi check input via raw API
func (d *RequestNodeDetail) checkRawAPI(raw *adaptor.QueryRawAPIResp, n *Node) error {
	if raw.ID == "" {
		return fmt.Errorf("raw api [%s] not exists for node [%s]", d.RawPath, n.Name)
	}
	if !exprx.EncodingEnum.Verify(raw.Content.EncodingIn) {
		return fmt.Errorf("raw api [%s] refer from node %s (encoding-in %s) unsupported, Valid: %v",
			d.RawPath, n.Name, raw.Content.EncodingIn, exprx.EncodingEnum.GetAll())
	}
	if !exprx.EncodingEnum.Verify(raw.Content.EncodingOut) {
		return fmt.Errorf("raw api [%s] refer from node %s (encoding-in %s) unsupported, Valid: %v",
			d.RawPath, n.Name, raw.Content.EncodingOut, exprx.EncodingEnum.GetAll())
	}

	if err := d.checkInputs(raw); err != nil {
		return err
	}

	return nil
}

// checkInputs verify the inputs of the request
func (d *RequestNodeDetail) checkInputs(raw *adaptor.QueryRawAPIResp) error {
	// TODO: check inputs parameters
	return nil
}
