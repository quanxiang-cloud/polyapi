package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/config"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"

	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// local debug test
const testEnv = "local"

type WorkflowSuiteStrutor struct {
	suite.Suite
	ctx       context.Context
	conf      *config.Config
	engin     *gin.Engine
	rawAPI    RawAPI
	polyAPI   PolyAPI
	idsPermit struct {
		raw         []string
		poly        string
		update      bool
		requestPoly bool
	}
	// idsForm struct {
	// 	raw    []string
	// 	poly   string
	// 	update bool
	// }
}

func _TestWorkflowStructor(t *testing.T) {
	start, end := 3, 3
	if init := false; init {
		start, end = 0, 2
	}
	for i := start; i <= end; i++ {
		s := &WorkflowSuiteStrutor{}
		switch testEnv {
		case "local", "debug", "test":
			s.idsPermit.raw = []string{
				"/system/form/base_pergroup_create.r", //create
				"/system/form/base_pergroup_update.r", //update
			}
			s.idsPermit.poly = "/system/poly/permissionInit.p"
			s.idsPermit.update = false
			s.idsPermit.requestPoly = true
		}

		if step := i; true {
			if step <= 0 {
				s.idsPermit.raw = []string{}
			}
			if step <= 1 {
				s.idsPermit.poly = ""
			}
			if step <= 2 {
				s.idsPermit.requestPoly = false
			}
		}

		suite.Run(t, s)
	}

}

func (s *WorkflowSuiteStrutor) SetupSuite() {
	var err error
	s.ctx = context.TODO()
	s.conf, err = config.NewConfig(fmt.Sprintf("./testdata/%s/polyapi.yaml", testEnv))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.conf)
	s.conf.Log.Level = 1 // warn
	logger.Logger = logger.New(&s.conf.Log)
	assert.Nil(s.T(), err)
	if s.rawAPI, err = CreateRaw(s.conf); err != nil {
		panic(err)
	}
	if s.polyAPI, err = CreatePoly(s.conf); err != nil {
		panic(err)
	}
	_, _ = CreateNamespaceOper(s.conf) // NOTE: init namespace oper
}

func (s *WorkflowSuiteStrutor) TestPermissionAll() {
	if len(s.idsPermit.raw) == 0 { // create raw
		s.testCreateRaw()
		return
	}
	s.testGetRaw()
	if s.idsPermit.poly == "" { // create poly
		s.testCreatePoly()

		return
	}
	if s.idsPermit.update {
		s.testUpdatePoly()
		s.testUpdatePolyScript(`"script V1"`, "")
	}

	s.testBuildPoly()
	//s.testRequestPoly()
}

func (s *WorkflowSuiteStrutor) testCreateRaw() {
	if true {
		file, _ := os.Open(fmt.Sprintf("./testdata/%s/swaggerapi_structor_mini.json", testEnv))
		all, _ := io.ReadAll(file)
		req := &RegReq{
			Service:   "",
			Swagger:   string(all),
			Namespace: "/system/form",
			Owner:     consts.SystemName,
			OwnerName: consts.SystemTitle,
			Host:      "structor",
			AuthType:  "none",
			Schema:    "http",
		}
		resp, err := s.rawAPI.RegSwagger(s.ctx, req)
		fmt.Println(resp)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), resp)
	}
	if false {
		file, _ := os.Open(fmt.Sprintf("./testdata/%s/swagger_form.json", testEnv))
		all, _ := io.ReadAll(file)
		req := &RegReq{
			Swagger: string(all),
		}
		resp, err := s.rawAPI.RegSwagger(s.ctx, req)
		fmt.Println(resp)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), resp)
	}
	if false {
		file, _ := os.Open(fmt.Sprintf("./testdata/%s/swagger_kd.json", testEnv))
		all, _ := io.ReadAll(file)
		req := &RegReq{
			Swagger: string(all),
		}
		resp, err := s.rawAPI.RegSwagger(s.ctx, req)
		fmt.Println(resp)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), resp)
	}
}

func (s *WorkflowSuiteStrutor) testGetRaw() {
	req := &QueryReq{
		APIPath: s.idsPermit.raw[0],
	}
	resp, err := s.rawAPI.Query(s.ctx, req)
	if err != nil {
		fmt.Printf("get raw %+v %v\n", resp, err)
		return
	}

	fmt.Printf("%+v\n", resp.Content)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteStrutor) testCreatePoly() {
	req := &PolyCreateReq{
		Owner:     consts.SystemName,
		OwnerName: consts.SystemTitle,
		Namespace: "/system/poly",
		Name:      "permissionInit",
		Title:     "应用初始化",
		//Access:    []string{"read", "execute"},
		Method: "POST",
	}
	resp, err := s.polyAPI.Create(s.ctx, req)
	fmt.Printf("create %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	if err == nil {
		s.testGetPoly(resp.APIPath)
	}
}

func (s *WorkflowSuiteStrutor) testUpdatePoly() {
	req := &PolyUpdateArrangeReq{
		APIPath: s.idsPermit.poly,
		Arrange: "arrange V1",
	}
	resp, err := s.polyAPI.UpdateArrange(s.ctx, req)
	fmt.Printf("update V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPoly("")
}

func (s *WorkflowSuiteStrutor) testGetPoly(id string) {
	if id == "" {
		id = s.idsPermit.poly
	}
	req := &PolyGetArrangeReq{APIPath: id}
	resp, err := s.polyAPI.GetArrange(s.ctx, req)
	fmt.Printf("get %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteStrutor) testUpdatePolyScript(script, doc string) {
	req := &PolyUpdateScriptReq{
		APIPath: s.idsPermit.poly,
		Script:  script,
		Doc:     doc,
	}
	resp, err := s.polyAPI.UpdateScript(s.ctx, req)
	fmt.Printf("update script V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPolyScript()
	s.testRequestPoly()
}

func (s *WorkflowSuiteStrutor) setPolyActive(active uint) error {
	req := &PolyActiveReq{
		APIPath: s.idsPermit.poly,
		Active:  active,
	}
	_, err := s.polyAPI.Active(s.ctx, req)
	assert.Nil(s.T(), err)
	return err
}

func (s *WorkflowSuiteStrutor) testBuildPoly() {
	if err := s.setPolyActive(rule.ActiveDisable); err != nil {
		return
	}
	req := &PolyBuildReq{
		APIPath: s.idsPermit.poly,
		Arrange: s.genSampleArrange(),
	}
	resp, err := s.polyAPI.Build(s.ctx, req)
	fmt.Printf("build V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)

	reqUA := &PolyUpdateArrangeReq{
		APIPath: s.idsPermit.poly,
		Arrange: "{}",
	}
	respUA, err := s.polyAPI.UpdateArrange(s.ctx, reqUA)
	fmt.Println(respUA, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), respUA)

	if err := s.setPolyActive(rule.ActiveEnable); err != nil {
		return
	}
	if script := s.testGetPolyScript(); script != "" {
		s.testUpdatePolyScript(script, "{}")
	}

	s.testRequestPoly()
}

func (s *WorkflowSuiteStrutor) testGetPolyScript() string {
	req := &PolyGetScriptReq{APIPath: s.idsPermit.poly}
	resp, err := s.polyAPI.GetScript(s.ctx, req)
	fmt.Printf("get script %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	if err != nil {
		return ""
	}
	return resp.Script
}

func (s *WorkflowSuiteStrutor) testRequestPoly() {
	if !s.idsPermit.requestPoly {
		return
	}

	f, err := os.Open(fmt.Sprintf("./testdata/%s/polytest_structor.testdata", testEnv))
	if err != nil {
		s.T().Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		s.T().Fatal(err)
	}

	req := &apiprovider.RequestReq{
		APIPath: s.idsPermit.poly,
		Method:  consts.MethodPost,
		Body:    b,
	}
	resp, err := s.polyAPI.Request(s.ctx, req)
	if err == nil {
		body := string(resp.Response)
		resp.Response = nil
		fmt.Printf("request %+v\n%s\n", resp, body)
	}

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

/*
// qyTmpScript_/system/poly/permissionInit_permissionInit_2021-11-20T09:27:02CST
var _tmp = function(){
  var d = { "__input": __input, } // qyAllLocalData

  d.start = __input.body

  d.start.header = d.start.header || {}
  d.start._ = [
    pdCreateNS('/system/app',d.start.appID,'应用'),
    pdCreateNS('/system/app/'+d.start.appID,'poly','API编排'),
    pdCreateNS('/system/app/'+d.start.appID,'raw','原生API'),
    pdCreateNS('/system/app/'+d.start.appID+'/raw','customer','代理第三方API'),
    pdCreateNS('/system/app/'+d.start.appID+'/raw','inner','平台API'),
    pdCreateNS('/system/app/'+d.start.appID+'/raw/inner','form','表单模型API'),
    pdCreateNS('/system/app/'+d.start.appID+'/raw/inner/form','form','表单API'),
    pdCreateNS('/system/app/'+d.start.appID+'/raw/inner/form','custom','模型API'),
  ]
  if (true) { // req1, create
    var _apiPath = format("http://structor/api/v1/structor/%v/base/permission/perGroup/create" ,d.start.appID)
    var _t = {
        "name": d.start.name,
        "description": d.start.description,
        "types": d.start.types,
      }
    var _th = pdNewHttpHeader()
    pdAddHttpHeader(_th, "Content-Type", "application/json")

    var _tk = '';
    var _tb = pdAppendAuth(_tk, 'none', _th, pdToJson(_t))
    d.req1 = pdToJsobj("json", pdHttpRequest(_apiPath, "POST", _tb, _th, pdQueryUser(true)))
  }
  d.cond1 = { y: false, }
  if (d.req1.code==0) {
    d.cond1.y = true
    if (true) { // req2, update
      var _apiPath = format("http://structor/api/v1/structor/%v/base/permission/perGroup/update" ,d.start.appID)
      var _t = {
          "id": d.req1.data.id,
          "scopes": d.start.scopes,
        }
      var _th = pdNewHttpHeader()
      pdAddHttpHeader(_th, "Content-Type", "application/json")

      var _tk = '';
      var _tb = pdAppendAuth(_tk, 'none', _th, pdToJson(_t))
      d.req2 = pdToJsobj("json", pdHttpRequest(_apiPath, "POST", _tb, _th, pdQueryUser(true)))
    }
  }

  d.end = {
    "createNamespaces": d.start._,
    "req1": d.req1,
    "req2": sel(d.cond1.y,d.req2,undefined),
  }
  return pdToJsonP(d.end)
}; _tmp();
*/

func (s *WorkflowSuiteStrutor) genCreateAppNameSpaceArrange() expr.ValArray {
	const (
		appIDPlaceholder = "$$"
		realAppID        = "start.appID"
	)
	_split := func(path string) (string, string) {
		p, n := apipath.Split(path)
		if n == appIDPlaceholder {
			return fmt.Sprintf(`'%s'`, p), realAppID
		}
		ps := strings.Split(p, appIDPlaceholder)
		root, sub := ps[0], ""
		if len(ps) >= 2 {
			sub = ps[1]
		}
		if sub == "" {
			return fmt.Sprintf(`'%s'+%s`, root, realAppID), fmt.Sprintf(`'%s'`, n)
		}
		return fmt.Sprintf(`'%s'+%s+'%s'`, root, realAppID, sub), fmt.Sprintf(`'%s'`, n)
	}
	createPaths := app.GetCreateAppPaths(appIDPlaceholder)
	ret := expr.ValArray{}
	for _, v := range createPaths {
		ns, n := _split(v.Path)
		val := expr.Value{
			Type: expr.Enum(exprx.ExprTypeDirectExpr),
			Data: exprx.FlexJSONObject{
				D: exprx.ValDirectExpr(fmt.Sprintf("//##%s(%s,%s,'%s')",
					consts.PDCreateNS, ns, n, v.Title)),
			},
		}
		ret = append(ret, val)
	}
	return ret
}

func (s *WorkflowSuiteStrutor) genSampleArrange() string {
	a := &arrange.Arrange{
		Info: &arrange.APIInfo{
			Namespace: "/system/poly",
			Name:      "permissionInit",
			Title:     "应用初始化",
			Desc:      "structor permissionInit",
			Version:   "v1.0.0",
			Method:    "POST",
		},
		Nodes: []arrange.Node{
			arrange.Node{
				Name:      "start",
				Type:      "input",
				Desc:      "the only entrance of the arrange",
				NextNodes: []string{"req1"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 1, Y: 10}},
				Detail: arrange.FlexJSONObject{
					D: arrange.InputNodeDetail{
						Consts: exprx.ValueSet{
							expr.InputValue{
								Name: consts.PolyDummyName,
								Type: expr.ValTypeArray,
								Data: expr.FlexJSONObject{ // create app namespace array
									D: s.genCreateAppNameSpaceArrange(),
								},
							},
						},
						Inputs: []exprx.ValueDefine{
							exprx.ValueDefine{
								InputValue: expr.InputValue{
									Type:     "string",
									Name:     "appID",
									Title:    "应用ID",
									In:       "body",
									Required: true,
								},
							},
							exprx.ValueDefine{
								InputValue: expr.InputValue{
									Type:     "string",
									Name:     "name",
									Title:    "应用名",
									In:       "body",
									Required: true,
								},
							},
							exprx.ValueDefine{
								InputValue: expr.InputValue{
									Type:     "string",
									Name:     "description",
									Title:    "应用描述",
									In:       "body",
									Required: false,
								},
							},
							exprx.ValueDefine{
								InputValue: expr.InputValue{
									Type:     "string",
									Name:     "types",
									Title:    "应用描述",
									In:       "body",
									Required: false,
								},
							},
							exprx.ValueDefine{
								InputValue: expr.InputValue{
									Type:     "array",
									Name:     "scopes",
									Title:    "",
									In:       "body",
									Required: true,
									Data: expr.FlexJSONObject{
										D: expr.ValArray{
											expr.Value{
												Type:  "object",
												Name:  "type",
												Title: "",
												Data: expr.FlexJSONObject{
													D: expr.ValObject{
														expr.Value{
															Type:  "number",
															Name:  "type",
															Title: "",
														},
														expr.Value{
															Type:  "string",
															Name:  "id",
															Title: "",
														},
														expr.Value{
															Type:  "string",
															Name:  "name",
															Title: "",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			arrange.Node{
				Name:      "req1",
				Type:      "request",
				Desc:      "create",
				NextNodes: []string{"cond1", "end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 10}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawPath: s.idsPermit.raw[0],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Name: "Content-Type",
								Type: "string",
								Data: expr.FlexJSONObject{
									D: expr.NewStringer(consts.MIMEJSON, expr.ValTypeString),
								},
								In: "header",
							},
							arrange.InputValue{
								Name:  "appID",
								Type:  "path",
								Field: "start.appID",
								In:    "path",
							},

							arrange.InputValue{
								Name:  "name",
								Type:  "field",
								Field: "start.name",
							},
							arrange.InputValue{
								Name:  "description",
								Type:  "field",
								Field: "start.description",
							},
							arrange.InputValue{
								Name:  "types",
								Type:  "field",
								Field: "start.types",
							},
						},
					},
				},
			},

			arrange.Node{
				Name:      "cond1",
				Type:      "if",
				Desc:      "req1 is ok",
				NextNodes: []string{},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.IfNodeDetail{
						Yes: "req2",
						No:  "",
						Cond: arrange.CondExpr{
							Op: "",
							Value: arrange.Value{
								Desc: "",
								Type: "exprcmp",
								Data: arrange.FlexJSONObject{
									D: arrange.ValExprCmp{
										LValue: arrange.Value{
											Type:  "number",
											Field: "req1.code",
										},
										Cmp: "eq",
										RValue: arrange.Value{
											Type: "number",
											Data: arrange.FlexJSONObject{
												D: "0",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			arrange.Node{
				Name:      "req2",
				Type:      "request",
				Desc:      "update",
				NextNodes: []string{"end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawPath: s.idsPermit.raw[1],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Name: "Content-Type",
								Type: "string",
								Data: expr.FlexJSONObject{
									D: expr.NewStringer(consts.MIMEJSON, expr.ValTypeString),
								},
								In: "header",
							},
							arrange.InputValue{
								Name:  "appID",
								Type:  "path",
								Field: "start.appID",
								In:    "path",
							},

							arrange.InputValue{
								Name:  "id",
								Type:  "field",
								Field: "req1.data.id",
							},
							arrange.InputValue{
								Name:  "scopes",
								Type:  "field",
								Field: "start.scopes",
							},
						},
					},
				},
			},
			arrange.Node{
				Name: "end",
				Type: "output",
				Desc: "the only end of this arrange",
				//Disp: arrange.FlexJSONObject{D: arrange.DispDemo{X: 15, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.OutputNodeDetail{
						Body: arrange.Value{
							Type: "object",
							Data: arrange.FlexJSONObject{
								D: arrange.ValObject{
									arrange.Value{
										Type:  "field",
										Name:  "createNamespaces",
										Field: "start._",
									},
									arrange.Value{
										Type:  "field",
										Field: "req1",
									},
									arrange.Value{
										Type:  "field",
										Field: "req2",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	b, err := json.MarshalIndent(a, "", "  ")
	//b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	ss := string(b)
	fmt.Println("arrange", ss)
	return ss
}
