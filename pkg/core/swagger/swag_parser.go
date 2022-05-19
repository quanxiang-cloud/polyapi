package swagger

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/version"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"

	"github.com/quanxiang-cloud/cabin/logger"
)

type swagDefine map[string]json.RawMessage

func parseDefine(define json.RawMessage) (swagDefine, error) {
	out := swagDefine{}
	if len(define) > 0 {
		if err := json.Unmarshal(define, &out); err != nil {
			return nil, err
		}
	}

	return out, nil
}

//------------------------------------------------------------------------------

// SwagParser prase swagger
type SwagParser struct {
	swag      *SwagDoc
	api       *adaptor.RawAPIFull
	define    swagDefine
	swagAPI   *SwagAPI
	withRefer bool
}

// ParseSwagger parse swagger bytes
func (p *SwagParser) ParseSwagger(content []byte, cfg *APIServiceConfig) ([]*adaptor.RawAPIFull, error) {
	p.swag = &SwagDoc{}
	doc := p.swag
	if err := json.Unmarshal(content, doc); err != nil {
		return nil, err
	}
	if err := doc.Consts.DelayedJSONDecode(); err != nil {
		return nil, err
	}

	host := selPriority(cfg.Host, doc.Host) // use config.host firstly
	if host == "" {
		return nil, errors.New("missing host")
	}
	if err := rule.ValidateHost(host); err != nil {
		return nil, err
	}

	schema := selScheme(append([]string{cfg.Schema}, doc.Schemes...)) //doc.Schemes
	if schema == "" {
		return nil, errors.New("missing schema")
	}

	list := make([]*adaptor.RawAPIFull, 0, len(doc.Paths))
	for key, apiData := range doc.Paths {
		if err := apipath.ValidateAPIPath(&key); err != nil {
			logger.Logger.Warn("invalid-api-path", err.Error())
			return nil, errcode.ErrAPIPath.FmtError(err.Error())
		}

		ks := strings.Split(key, "?") // "apiPath?action"
		apiPath, action := ks[0], ""
		if len(ks) >= 2 {
			action = ks[1]
		}
		for method, v := range apiData {
			_method := strings.ToUpper(method)
			if !expr.MethodEnum.Verify(_method) {
				logger.Logger.Warnf("unsupported method %s in swagger", _method)
				continue
			}

			apiName := cfg.selectAPIName(v.Name, len(doc.Paths)*len(apiData))
			if err := rule.ValidateName(apiName, rule.MaxNameLength-2, true); err != nil {
				return nil, err
			}

			p.swagAPI = v

			relPath := apipath.Join(doc.BasePath, apiPath)
			if strings.HasSuffix(apiPath, "/") && !strings.HasSuffix(relPath, "/") {
				relPath += "/"
			}
			p.api = &adaptor.RawAPIFull{
				ID:        "?",
				Owner:     "?",
				Namespace: cfg.Namespace,
				Name:      apiName, //unique name
				Service:   cfg.Service,
				Version:   selPriority(cfg.Version, doc.Info.Version),
				Path:      relPath,
				URL:       apipath.MakeRequestURL(schema, host, relPath),
				Host:      host,

				Title:    v.Summary,
				Action:   action,
				Method:   _method,
				Desc:     v.Desc,
				Schema:   schema,
				Access:   0,
				Active:   rule.ActiveEnable,
				AuthType: cfg.AuthType,
			}
			api := p.api

			if err := rule.CheckCharSet(api.Title, api.Desc); err != nil {
				return nil, err
			}
			if err := rule.CheckDescLength(api.Desc); err != nil {
				return nil, err
			}

			// BUG: do not DelayedJSONDecode here, maybe missing type
			// if err := v.Consts.DelayedJSONDecode(); err != nil {
			// 	return nil, err
			// }

			predefVals, err := mergePredefValue(v.Consts, doc.Consts)
			if err != nil {
				return nil, err
			}

			api.Content = &adaptor.RawAPIContent{
				ID:      api.ID,
				Path:    relPath,
				Action:  action,
				Summary: v.Summary,
				Desc:    v.Desc,

				EncodingIn:  selEncoding(v.EncodingsIn),
				EncodingOut: selEncoding(v.EncodingsOut),
				Method:      _method,
				Consts:      predefVals,

				BasePath: "/",
			}
			if err := api.Content.Check(); err != nil {
				return nil, err
			}
			api.Doc = &adaptor.APIDoc{
				// for swagger
				Defines:             doc.Defines,
				Parameters:          v.Parameters,
				Responses:           v.Responses,
				Security:            v.Security,
				SecurityDefinitions: doc.SecurityDefinitions,
			}

			if _, err := p.GenerateSwagger(p.api); err != nil {
				return nil, err
			}

			// remove from doc
			api.Doc.Defines = nil
			api.Doc.Parameters = nil
			api.Doc.Responses = nil
			api.Doc.Security = nil
			api.Doc.SecurityDefinitions = nil

			list = append(list, api)
		}
	}
	return list, nil
}

func (p *SwagParser) parseInOut() error {
	p.api.Doc.FmtInOut.Method = p.api.Method
	p.api.Doc.FmtInOut.SetAccessURL(apipath.Join(p.api.Namespace, p.api.Name))

	define, err := parseDefine(p.api.Doc.Defines)
	if err != nil {
		return err
	}
	p.define = define
	if err := p.parseSwagInput(); err != nil {
		return fmt.Errorf("parse parameters fail: %s", err.Error())
	}
	if err := p.parseSwagOutput(); err != nil {
		return fmt.Errorf("parse responses fail: %s", err.Error())
	}
	if err := p.genSample(); err != nil {
		return err
	}

	//DEBUG:
	// b, err := GetRawAPIDoc(p.api, "http://xxx.com", "raw", false)
	// if err != nil {
	// 	panic(err)
	// }
	// println(string(b))

	return nil
}

func (p *SwagParser) genSample() error {
	p.api.Doc.Version = version.FullVersion()
	if err := p.genSampleInput(false); err != nil {
		return err
	}
	if err := p.genSampleOutput(false); err != nil {
		return err
	}

	if err := p.genSampleInput(true); err != nil {
		return err
	}
	if err := p.genSampleOutput(true); err != nil {
		return err
	}
	return nil
}

func (p *SwagParser) genSampleInput(titleFirst bool) error {
	f := &p.api.Doc.FmtInOut
	s := &f.SampleIn[GetSampleIndex(titleFirst)]

	in := &f.Input
	hideObj := value.Object{}
	var body value.JSONValue
	for i := 0; i < len(in.Inputs); i++ {
		x := &in.Inputs[i]
		switch {
		case x.In.IsBody():
			body = x.CreateSampleData(nil, titleFirst)
		case x.In.IsHeader():
			hName := x.GetName(titleFirst)
			val := x.Mock
			if val == "" {
				val = value.RandString("", -1)
			}
			addHTTPHeader(&s.Header, hName, val)
		case x.In.IsPath():
			var val value.JSONValue
			if x.Mock != "" {
				v := value.String(x.Mock)
				val = &v
			} else {
				val = x.CreateSampleData(nil, titleFirst)
			}
			hideObj.AddElement(x.GetName(titleFirst), val)
		}
	}

	// NOTE: Add $polyapi_hide$ and x_polyapi_signature
	if body == nil {
		body = &value.Object{}
	}
	if b, ok := body.(*value.Object); ok {
		if err := b.AddElement(inputBodyPolySignatureValue.GetName(titleFirst),
			inputBodyPolySignatureValue.CreateSampleData(nil, titleFirst)); err != nil {
			return err
		}
		if len(hideObj) > 0 {
			if err := b.AddElement(inputBodyPolyReservedValue.GetName(titleFirst), hideObj); err != nil {
				return err
			}
		}
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	s.Body = b
	return nil
}

func (p *SwagParser) genSampleOutput(titleFirst bool) error {
	f := &p.api.Doc.FmtInOut
	s := &f.SampleOut[GetSampleIndex(titleFirst)]
	out := &f.Output
	var body value.JSONValue
	for i := 0; i < len(out.Doc); i++ {
		x := &out.Doc[i]
		switch {
		case x.In.IsBody():
			body = x.CreateSampleData(nil, titleFirst)
		case x.In.IsHeader():
			//TODO output header?
		}
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	s.Resp = b
	return nil
}

// nameFilter verify name conflict
type nameFilter map[string]struct{}

func (f nameFilter) Filt(name string) error {
	if _, ok := f[name]; !ok {
		f[name] = struct{}{}
	} else {
		return errcode.ErrDuplicateName.FmtError(name)
	}
	return nil
}

// swagAPI.parameters -> p.api.Content.FmtInOut
func (p *SwagParser) parseSwagInput() error {
	if len(p.api.Doc.Parameters) == 0 {
		return nil
	}

	p.withRefer = false
	var inputs []SwagInputValue
	if err := json.Unmarshal(p.api.Doc.Parameters, &inputs); err != nil {
		return err
	}
	f := &p.api.Doc.FmtInOut
	f.Input.Inputs = append(inputSignatureDefine, f.Input.Inputs...)

	filter := nameFilter{}
	hasPathVal := false
	bodyIndex := -1
	for i := 0; i < len(inputs); i++ {
		v := &inputs[i]
		switch v.In {
		case "header", "path", "formData", "query":
			if !hasPathVal && v.In == "path" {
				hasPathVal = true
			}

			if err := filter.Filt(fmt.Sprintf("%s@%s", v.Name, v.In)); err != nil {
				return err
			}

			if v.Ref != "" {
				if err := p.mergeRefer(v, v.Ref); err != nil {
					return fmt.Errorf("merge refer %s error: %s", v.Ref, err.Error())
				}
				v.Ref = ""
			}
			d := expr.ValueDefine{
				InputValue: expr.InputValue{
					Name:     v.Name,
					Type:     expr.ValTypeString,
					In:       expr.Enum(v.In),
					Desc:     v.Desc,
					Required: v.Required,
				},
			}
			f.Input.Inputs = append(f.Input.Inputs, d)

		case "body":
			if err := filter.Filt("<root>@body"); err != nil {
				return err
			}

			val, err := p.parseSchema(v.Schema, 0, true)
			if err != nil {
				return fmt.Errorf("parse body error: %s", err.Error())
			}

			// NOTE: body is not object, change to {"$body$":originBody}
			if val.Type != expr.ValTypeObject {
				val.Name = polysign.XPolyCustomerBodyRoot
				obj := expr.Value{
					Type:  expr.ValTypeObject,
					Name:  v.Name,
					Title: "",
					Data: expr.FlexJSONObject{
						D: &expr.ValObject{
							*val,
						},
					},
				}
				val = &obj // change val to *Object
			}

			d := expr.ValueDefine{
				InputValue: expr.InputValue{
					Name:     v.Name,
					Type:     val.Type,
					Title:    val.Title,
					In:       expr.ParaTypeBody,
					Desc:     v.Desc,
					Data:     val.Data,
					Required: v.Required,
				},
			}
			f.Input.Inputs = append(f.Input.Inputs, d)
			bodyIndex = len(f.Input.Inputs) - 1
		default:
			return errcode.ErrParameterIn.FmtError(v.In)
		}
	}

	// NOTE: add $polyapi_hide$ field to body
	if bodyIndex < 0 {
		d := expr.ValueDefine{
			InputValue: expr.InputValue{
				Name: "",
				Type: expr.ValTypeObject,
				In:   expr.ParaTypeBody,
				Desc: "body root",
				Data: expr.FlexJSONObject{
					D: &expr.ValObject{},
				},
			},
		}
		f.Input.Inputs = append(f.Input.Inputs, d)
		bodyIndex = len(f.Input.Inputs) - 1
	}
	body := &f.Input.Inputs[bodyIndex]
	if obj, ok := body.Data.D.(*expr.ValObject); ok {
		*obj = append(*obj, inputBodyPolySignatureValue)
		if hasPathVal {
			*obj = append(*obj, inputBodyPolyReservedValue)
		}
	}

	if p.withRefer { // rewrite inputs with refer
		p.withRefer = false
		b, err := json.Marshal(inputs)
		if err != nil {
			return err
		}
		p.api.Doc.Parameters = b
	}
	return nil
}

func addHTTPHeader(h *http.Header, k, v string) {
	if *h == nil {
		*h = http.Header{}
	}
	h.Add(k, v)
}

func (p *SwagParser) parseSwagOutput() error {
	if len(p.api.Doc.Responses) == 0 {
		return nil
	}

	p.withRefer = false
	var outputs SwagResponsesSchema
	if err := json.Unmarshal(p.api.Doc.Responses, &outputs); err != nil {
		return err
	}

	f := &p.api.Doc.FmtInOut
	if suc := outputs.Success; suc != nil {
		if s := suc.Schama; s != nil {
			if p.GetType(s) == "" {
				//NOTE: fix compatible error
				logger.Logger.Warnf("Responses of %s %s missing type", p.api.Path, string(p.api.Doc.Responses))
				s.Type = "object" //BUG: maybe missing type
			}
			val, err := p.parseSchema(s, 0, false)
			if err != nil {
				return fmt.Errorf("parse response error: %s", err.Error())
			}
			d := expr.ValueDefine{
				InputValue: expr.InputValue{
					Name:  "",
					Type:  val.Type,
					In:    expr.ParaTypeBody,
					Desc:  suc.Desc,
					Data:  val.Data,
					Title: val.Title,
				},
				//Required: v.Required,
			}
			f.Output.Doc = append(f.Output.Doc, d)
		}
		for k, v := range suc.Header {
			d := expr.InputValue{
				Name: k,
				Type: expr.ValTypeString,
				In:   expr.ParaTypeHeader,
				Desc: v.Desc,
			}
			f.Output.Header = append(f.Output.Header, d)
		}
	}

	if p.withRefer { // rewrite inputs with refer
		p.withRefer = false
		b, err := json.Marshal(outputs)
		if err != nil {
			return err
		}
		p.api.Doc.Responses = b
	}
	return nil
}

func (p *SwagParser) parseSchema(s *SwagSchema, depth int, input bool) (val *expr.Value, err error) {

	if s.Ref != "" {
		if err := p.mergeRefer(s, s.Ref); err != nil {
			return nil, err
		}
		s.Ref = ""
	}

	if depth == 0 && len(s.AllOf) > 0 {
		if err := p.mergeAllOf(s, depth, input); err != nil {
			return nil, err
		}
	}

	if s == nil {
		return nil, fmt.Errorf("missing schema field")
	}

	ty := p.GetType(s)
	switch ty {
	case "string", "number", "boolean":
		return p.parseSchemaSingleValue(ty, s, depth)
	case "integer":
		return p.parseSchemaSingleValue("number", s, depth)
	case "null", "":
		//NOTE: fix compatible error
		logger.Logger.Warnf("found null/empty type at api %s", p.api.Path)
		s.Type = "object"
		p.withRefer = true
		fallthrough
	case "object":
		return p.parseSchemaObject(s, depth, input)
	case "array":
		return p.parseSchemaArray(s, depth, input)
	}
	return nil, fmt.Errorf("unsupport schema type [%s] at %s: %#v", ty, p.api.Path, s)
}

func (p *SwagParser) mergeAllOf(s *SwagSchema, depth int, input bool) error {
	allOf := s.AllOf
	s.AllOf = nil
	for _, b := range allOf {
		if err := json.Unmarshal(b, s); err != nil {
			return err
		}
		p.withRefer = true

		if s.Ref != "" {
			if err := p.mergeRefer(s, s.Ref); err != nil {
				return err
			}
			s.Ref = ""
		}
	}
	return nil
}

func (p *SwagParser) parseSchemaSingleValue(kind string, v *SwagSchema, depth int) (*expr.Value, error) {
	x := &expr.Value{
		Type:  expr.Enum(kind),
		Desc:  v.Desc,
		Title: v.Title,
	}
	return x, nil
}

func (p *SwagParser) parseSchemaObject(s *SwagSchema, depth int, input bool) (*expr.Value, error) {
	xx := expr.ValObject{}
	x := &expr.Value{
		Type:  expr.ValTypeObject,
		Desc:  s.Desc,
		Title: s.Title,
		Data: expr.FlexJSONObject{
			D: &xx,
		},
	}
	// if input && depth == 0 {
	// 	xx = append(xx, inputBodyPolyReservedValue)
	// }
	if s.Properties != nil {
		//NOTE: with $body$ change body type, dont allow other fields
		if _, ok := s.Properties[polysign.XPolyCustomerBodyRoot]; ok {
			if len(s.Properties) != 1 {
				return nil, errcode.ErrUniqueCustomerRootField.FmtError(polysign.XPolyCustomerBodyRoot)
			}
		}

		for k, v := range s.Properties {
			val, err := p.parseSchema(v, depth+1, input)
			if err != nil {
				return nil, err
			}
			val.Name = k
			xx = append(xx, *val)
		}
	}

	return x, nil
}
func (p *SwagParser) parseSchemaArray(s *SwagSchema, depth int, input bool) (*expr.Value, error) {
	xx := expr.ValArray{}
	x := &expr.Value{
		Type:  expr.ValTypeArray,
		Desc:  s.Desc,
		Title: s.Title,
		Data: expr.FlexJSONObject{
			D: &xx,
		},
	}
	if s.Items != nil {
		val, err := p.parseSchema(s.Items, depth+1, input)
		if err != nil {
			return nil, err
		}
		xx = append(xx, *val)
	} else {
		return nil, fmt.Errorf("schema array without items")
	}

	return x, nil
}

func (p *SwagParser) mergeRefer(obj interface{}, ref string) error {
	if ref != "" {
		if d := p.getRefer(ref); d != nil {
			p.withRefer = true
			// TODO: keep the old data?
			if err := json.Unmarshal(d, obj); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s not found", ref)
		}
	}

	return nil
}

func (p *SwagParser) getRefer(ref string) json.RawMessage {
	if idx := strings.LastIndex(ref, "/"); idx >= 0 {
		name := ref[idx+1:]
		if v, ok := p.define[name]; ok {
			return v
		}
	}
	return nil
}

// GetType Type maybe "array" or ["array", "null"]
// so get without null
func (p *SwagParser) GetType(s *SwagSchema) string {
	switch t := s.Type.(type) {
	case string:
		return t
	case []interface{}:
		for _, _tt := range t {
			if tt, ok := _tt.(string); ok {
				if tt != "null" {
					p.withRefer = true
					s.Type = tt
					return tt
				}
			}
		}
	}
	return ""
}
