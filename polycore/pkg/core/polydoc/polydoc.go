package polydoc

import (
	"encoding/json"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

func init() {
	// init the swagger generator function, reset the adaptor
	arrange.PolyDocGennerator = genPolyDoc
}

// genPolySwagger generate and set api.Doc
func genPolyDoc(poly *arrange.Arrange, in *inputNodeDetail, out *outputNodeDetail) (*adaptor.APIDoc, error) {
	if poly.Info == nil {
		return nil, fmt.Errorf("missing arrange.Info")
	}

	service := "polyapi"
	method := poly.Info.Method
	api := &adaptor.RawAPIFull{
		Namespace: poly.Info.Namespace,
		Name:      apipath.GenerateAPIName(poly.Info.Name, ".p"),
		Version:   poly.Info.Version,
		Service:   service,
		Desc:      poly.Info.Desc,
		Method:    method,
		Schema:    consts.SchemaHTTP,
		Host:      service,
		Content: &adaptor.RawAPIContent{
			Summary:     poly.Info.Title,
			Desc:        poly.Info.Desc,
			BasePath:    "/",
			EncodingIn:  consts.PolyEncoding,
			EncodingOut: consts.PolyEncoding,
			Consts:      expr.ValueSet(in.Consts),
			Path:        fmt.Sprintf("/api/v1/polyapi/request/%s", poly.Info.APIPath()),
			Method:      method,
		},
		Doc: &adaptor.APIDoc{},
	}

	if err := genSwaggerInput(api, in); err != nil {
		return nil, err
	}
	if err := genSwaggerOutput(api, out); err != nil {
		return nil, err
	}

	var parser swagger.SwagParser
	if _, err := parser.GenerateSwagger(api); err != nil {
		return nil, err
	}

	// remove from doc
	api.Doc.Defines = nil
	api.Doc.Parameters = nil
	api.Doc.Responses = nil
	api.Doc.Security = nil
	api.Doc.SecurityDefinitions = nil

	return api.Doc, nil
}

//------------------------------------------------------------------------------

func genSwaggerInput(api *adaptor.RawAPIFull, in *inputNodeDetail) error {
	if in != nil {
		inputs := in.Inputs
		resp := make(swagParameters, 0, len(inputs))

		resp = append(resp, genBodySchema(inputs, true))
		for _, v := range inputs {
			if v.In.IsBody() {
				continue
			}
			e := swagValue{
				"name": v.GetName(false),
				"in":   v.In.String(),
			}
			if v.Desc != "" {
				e["description"] = v.Desc
			}
			if v.Required {
				e["required"] = v.Required
			}
			switch exprx.Enum(v.Type) {
			case exprx.ValTypeString:
				e["type"] = "string"
			case exprx.ValTypeBoolean:
				e["type"] = "boolean"
			case exprx.ValTypeNumber:
				e["type"] = "number"
			case exprx.ValTypeObject:
				e["type"] = "object"
			case exprx.ValTypeArray, exprx.ValTypeArrayString,
				/*exprx.ValTypeArrayElem,*/ exprx.XValTypeArrayStringElem:
				e["type"] = "array"
			default:
				continue
			}
			resp = append(resp, e)
		}
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return err
		}
		api.Doc.Parameters = b
	}
	return nil
}

func genBodySchema(in []exprx.ValueDefine, input bool) swagValue {
	required := []string{}
	schema := swagValue{
		"type": "object",
	}
	r := swagValue{
		"schema": schema,
	}
	if input {
		r["name"] = "root"
		r["description"] = "body inputs"
		r["in"] = exprx.ParaTypeBody.String()
	} else {
		r = schema
		r["description"] = "body response"
	}

	properties := swagObjectPropperties{}
	schema["properties"] = properties
	for _, v := range in {
		if !v.In.IsBody() {
			continue
		}
		e := swagValue{}
		if v.Desc != "" {
			e["description"] = v.Desc
		}
		if v.Title != "" {
			e["tile"] = v.Title
		}
		if input {
			if v.Default != "" {
				e["default"] = v.Default
			}
			if v.Mock != "" {
				e["mock"] = swagValue{
					"mock": v.Mock,
				}
			}
			if len(v.Enums) > 0 {
				e["enum"] = v.Enums
			}
		}
		switch exprx.Enum(v.Type) {
		case exprx.ExprTypeDirectExpr:
			e["type"] = "string"
		case exprx.ValTypeString:
			e["type"] = "string"
		case exprx.ValTypeBoolean, exprx.ValTypeNumber, exprx.ValTypeTimestamp:
			e["type"] = v.Type.String()
		case exprx.ValTypeObject:
			e["type"] = "object"
			if d, ok := v.Data.D.(*exprx.ValObject); ok {
				e = genObjectSchema(e, *d, 0)
			} else {
				// TODO: log?
				continue
			}

		case exprx.ValTypeArray:
			if d, ok := v.Data.D.(*exprx.ValArray); ok {
				e["type"] = "array"
				e["items"] = genArrayItems(e, *d, 0)
			} else {
				// TODO: log?
				continue
			}
		case exprx.ValTypeArrayString, exprx.XValTypeArrayStringElem:
			e["type"] = "array"
			e["items"] = swagValue{
				"type": "string",
			}
		default:
			// TODO: log?
			continue
		}
		if v.Required {
			required = append(required, v.GetName(false))
		}

		properties[v.GetName(false)] = e
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	return r
}

func genObjectSchema(v swagValue, in []exprx.Value, depth int) swagValue {
	required := []string{}
	properties := swagObjectPropperties{}
	schema := v

	schema["properties"] = properties
	for _, v := range in {
		e := swagValue{}
		if v.Desc != "" {
			e["description"] = v.Desc
		}
		switch v.Type {
		case exprx.ExprTypeDirectExpr:
			e["type"] = "string"
		case exprx.ValTypeString, exprx.ValTypeTimestamp:
			e["type"] = "string"
		case exprx.ValTypeBoolean, exprx.ValTypeNumber:
			e["type"] = v.Type.String()

		case exprx.ValTypeObject:
			e["type"] = "object"
			if d, ok := v.Data.D.(*exprx.ValObject); ok {
				e = genObjectSchema(e, *d, depth+1)
			} else {
				// TODO: log?
				continue
			}

		case exprx.ValTypeArray:
			if d, ok := v.Data.D.(*exprx.ValArray); ok {
				e["type"] = "array"
				e["items"] = genArrayItems(e, *d, depth+1)
			}
		// case exprx.ValTypeArrayElem:
		// 	if d, ok := v.Data.D.(*exprx.ValArrayElem); ok {
		// 		e["type"] = "array"
		// 		e["items"] = genArrayItems(e, d.Array, depth+1)
		// 	}
		case exprx.ValTypeArrayString, exprx.XValTypeArrayStringElem:
			e["type"] = "array"
			e["items"] = swagValue{
				"type": "string",
			}
		default:
			// TODO: log?
			continue
		}

		properties[v.GetName(false)] = e
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func genArrayItems(v swagValue, in []exprx.Value, depth int) swagValue {
	items := swagValue{
		"type": "",
	}
	if len(in) > 0 {
		v := in[0]
		e := items
		if v.Desc != "" {
			e["description"] = v.Desc
		}
		switch v.Type {
		case exprx.ExprTypeDirectExpr:
			e["type"] = "string"
		case exprx.ValTypeString:
			e["type"] = "string"
		case exprx.ValTypeBoolean, exprx.ValTypeNumber, exprx.ValTypeTimestamp:
			e["type"] = v.Type.String()
		case exprx.ValTypeObject:
			e["type"] = "object"
			if d, ok := v.Data.D.(*exprx.ValObject); ok {
				e = genObjectSchema(e, *d, depth+1)
			} else {
				// TODO: log?
			}
		case exprx.ValTypeArray:
			if d, ok := v.Data.D.(*exprx.ValArray); ok {
				e["type"] = "array"
				e["items"] = genArrayItems(e, *d, depth+1)
			}
		// case exprx.ValTypeArrayElem:
		// 	if d, ok := v.Data.D.(*exprx.ValArrayElem); ok {
		// 		e["type"] = "array"
		// 		e["items"] = genArrayItems(e, d.Array, depth+1)
		// 	}
		case exprx.ValTypeArrayString, exprx.XValTypeArrayStringElem:
			e["type"] = "array"
			e["items"] = swagValue{
				"type": "string",
			}
		default:
			// TODO: log?
		}
	}

	return items
}

func genSwaggerOutput(api *adaptor.RawAPIFull, out *outputNodeDetail) error {
	if out != nil {
		if err := genOutputDoc(out); err != nil {
			return err
		}

		resp := swagResponses{
			"200": &swagResponseObject{
				Desc:   "",
				Header: genOutputHeader(out),
				Schama: genBodySchema(out.Doc, false),
			},
		}
		b, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return err
		}
		api.Doc.Responses = b
	}

	return nil
}

func genOutputDoc(out *outputNodeDetail) error {
	return out.GenerateDoc()
}

func genOutputHeader(out *outputNodeDetail) swagResponseHeader {
	r := swagResponseHeader{}
	for _, v := range out.Header {
		if !v.In.IsHeader() {
			continue
		}
		h := &swagResponseHeaderItem{
			Type: exprx.ValTypeString.String(),
			Desc: v.Desc,
		}
		r[v.GetName(false)] = h
	}
	return r
}
