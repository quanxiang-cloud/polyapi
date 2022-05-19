package swagger

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// GenerateSwagger generate swagger from api.Doc
func (p *SwagParser) GenerateSwagger(api *adaptor.RawAPIFull) ([]byte, error) {
	p.api = api // init the api object

	// parse api.doc.FmtDoc from api.Doc.Parameter
	if err := p.parseInOut(); err != nil {
		err = fmt.Errorf("parse inout error for api %s: %s", p.api.Path, err.Error())
		return nil, err
	}

	doc, err := p.genSwagDoc()
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	p.api.Doc.Swagger = b // update swagger field
	return b, nil
}

func (p *SwagParser) genSwagDoc() (*SwagDoc, error) {
	api := p.api
	encIn, err := consts.ToMIME(api.Content.EncodingIn)
	if err != nil {
		return nil, err
	}
	encOut, err := consts.ToMIME(api.Content.EncodingOut)
	if err != nil {
		return nil, err
	}
	swagapi := &SwagAPI{
		Consts: api.Content.Consts,
		//Security:     api.Doc.Security,
		EncodingsIn:  []string{encIn},
		EncodingsOut: []string{encOut},
		Tags:         nil,
		Name:         apipath.BaseName(api.Name),
		Summary:      api.Content.Summary,
		Desc:         api.Content.Desc,
		Deprecated:   false,
		Parameters:   api.Doc.Parameters,
		Responses:    api.Doc.Responses,
	}
	doc := &SwagDoc{
		Consts: nil,
		//Defines:             api.Doc.Defines,
		//SecurityDefinitions: api.Doc.SecurityDefinitions,
		Version: SwagVersion,
		Host:    api.Host,
		Paths: map[string]SwagMethods{
			getPathWithAction(api.Content.Path, api.Action): SwagMethods{
				strings.ToLower(api.Content.Method): swagapi,
			},
		},
		Info: SwagInfo{
			Title:   api.Service,
			Version: api.Version,
			Desc:    fmt.Sprintf("auto generated"),
		},
		Schemes:      []string{api.Schema},
		BasePath:     api.Content.BasePath,
		EncodingsIn:  nil,
		EncodingsOut: nil,
	}

	return doc, nil
}
