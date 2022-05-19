package service

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"github.com/quanxiang-cloud/polyapi/internal/arrange"
	"github.com/quanxiang-cloud/polyapi/pkg/misc/config"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkflowSuite struct {
	suite.Suite
	ctx     context.Context
	conf    *config.Config
	r       *gin.Engine
	rawAPI  Raw
	polyApi PolyApi
	ids     struct {
		raw  []string
		poly string
	}
}

func _TestWorkflow(t *testing.T) {
	s := &WorkflowSuite{}
	s.ids.raw = []string{""}
	s.ids.poly = ""

	suite.Run(t, s)
}

func (s *WorkflowSuite) SetupSuite() {
	var err error
	s.ctx = logger.GenRequestID(context.TODO())
	s.conf, err = config.NewConfig("./testdata/polyapi.yaml")
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.conf)
	s.conf.Log.Level = 1 // warn
	err = logger.New(&s.conf.Log)
	assert.Nil(s.T(), err)
	s.rawAPI, err = CreateRaw(s.conf)
	s.polyApi, err = CreatePoly(s.conf)
}

func (s *WorkflowSuite) TestAll() {
	if len(s.ids.raw) == 0 { // create raw
		s.testCreateRaw()
		return
	}
	s.testGetRaw()
	if s.ids.poly == "" { // create poly
		s.testCreatePoly()
		return
	}

	s.testUpdatePoly()
	s.testUpdatePolyScript()
	s.testBuildPoly()
	//s.testRequestPoly()
}

func (s *WorkflowSuite) testCreateRaw() {
	file, _ := os.Open("./testdata/swagger_local.json")
	all, _ := io.ReadAll(file)
	req := &RegReq{
		File: all,
		Host: "127.0.0.1:9000",
	}
	resp, err := s.rawAPI.Reg(s.ctx, req)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}
func (s *WorkflowSuite) testGetRaw() {
	req := &QueryReq{
		ID: s.ids.raw[0],
	}
	resp, err := s.rawAPI.Query(s.ctx, req)
	fmt.Printf("get raw %+v %v\n", resp, err)
	fmt.Printf("%+v\n", resp.Content)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuite) testCreatePoly() {
	req := &PolyCreateReq{}
	resp, err := s.polyApi.Create(s.ctx, req)
	fmt.Printf("create %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPoly(resp.ID)
}

func (s *WorkflowSuite) testUpdatePoly() {
	req := &PolyUpdateArrangeReq{
		ID:      s.ids.poly,
		Arrange: "arrange V1",
	}
	resp, err := s.polyApi.UpdateArrange(s.ctx, req)
	fmt.Printf("update V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPoly("")
}

func (s *WorkflowSuite) testGetPoly(id string) {
	if id == "" {
		id = s.ids.poly
	}
	req := &PolyGetArrangeReq{ID: id}
	resp, err := s.polyApi.GetArrange(s.ctx, req)
	fmt.Printf("get %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuite) testUpdatePolyScript() {
	req := &PolyUpdateScriptReq{
		ID:     s.ids.poly,
		Script: `"script V1"`,
	}
	resp, err := s.polyApi.UpdateScript(s.ctx, req)
	fmt.Printf("update script V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPolyScript()
	s.testRequestPoly()
}

func (s *WorkflowSuite) testBuildPoly() {
	req := &PolyBuildReq{
		ID:      s.ids.poly,
		Arrange: s.genSampleArrange(),
	}
	resp, err := s.polyApi.Build(s.ctx, req)
	fmt.Printf("build V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPolyScript()
	s.testRequestPoly()

}

func (s *WorkflowSuite) testGetPolyScript() {
	req := &PolyGetScriptReq{ID: s.ids.poly}
	resp, err := s.polyApi.GetScript(s.ctx, req)
	fmt.Printf("get script %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuite) testRequestPoly() {
	req := &PolyRequestReq{
		ID:    s.ids.poly,
		Input: `{"id": 0}`,
	}
	resp, err := s.polyApi.Request(s.ctx, req)
	fmt.Printf("request %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuite) genSampleArrange() string {
	a := &arrange.Arrange{
		Uid:  s.ids.poly,
		Name: "polyApi",
		Desc: "polyApi sample",
		Nodes: []arrange.Node{
			arrange.Node{
				Name:      "start",
				Type:      "input",
				Desc:      "the only entrance of the arrange",
				NextNodes: []string{"req1"},
				Disp:      arrange.FlexJsonObject{D: arrange.DispDemo{X: 1, Y: 10}},
				Detail: arrange.FlexJsonObject{
					D: arrange.InputNodeDetail{},
				},
			},
			arrange.Node{
				Name:      "req1",
				Type:      "request",
				Desc:      "request1",
				NextNodes: []string{"cond1", "end"},
				Disp:      arrange.FlexJsonObject{D: arrange.DispDemo{X: 5, Y: 10}},
				Detail: arrange.FlexJsonObject{
					D: arrange.RequestNodeDetail{
						RawId:  s.ids.raw[0],
						Inputs: []arrange.ValueAccess{},
					},
				},
			},
			arrange.Node{
				Name:      "cond1",
				Type:      "if",
				Desc:      "check if need call request2",
				NextNodes: []string{},
				Disp:      arrange.FlexJsonObject{D: arrange.DispDemo{X: 5, Y: 20}},
				Detail: arrange.FlexJsonObject{
					D: arrange.IfNodeDetail{
						Yes: "req2",
						No:  "",
						Cond: arrange.CondExpr{
							Op: "",
							Value: arrange.Value{
								Desc: "req1.a /eq aaa",
								Type: "exprcmp",
								Value: arrange.FlexJsonObject{
									D: arrange.ExprCmp{
										LValue: arrange.Value{
											Type: "field",
											Value: arrange.FlexJsonObject{
												D: "req1.a",
											},
										},
										Cmp: "eq",
										RValue: arrange.Value{
											Type: "string",
											Value: arrange.FlexJsonObject{
												D: "aaa",
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
				Desc:      "request2",
				NextNodes: []string{"end"},
				Disp:      arrange.FlexJsonObject{D: arrange.DispDemo{X: 10, Y: 20}},
				Detail: arrange.FlexJsonObject{
					D: arrange.RequestNodeDetail{
						RawId:  s.ids.raw[1],
						Inputs: []arrange.ValueAccess{},
					},
				},
			},
			arrange.Node{
				Name: "end",
				Type: "output",
				Desc: "the only end of this arrange",
				Disp: arrange.FlexJsonObject{D: arrange.DispDemo{X: 15, Y: 20}},
				Detail: arrange.FlexJsonObject{
					D: arrange.OutputNodeDetail{
						AsHeader: false,
						Value: arrange.Value{
							Type: "mergeobj",
							Value: arrange.FlexJsonObject{
								D: arrange.ValMergeObj{
									arrange.Value{
										Type: "filter",
										Value: arrange.FlexJsonObject{
											D: arrange.ValFilter{
												Source: "req1",
												White: arrange.FieldMap{
													"a": "A",
													"b": "",
												},
											},
										},
									},
									arrange.Value{
										Type: "exprsel",
										Value: arrange.FlexJsonObject{
											D: arrange.ExprSel{
												Cond: arrange.CondExpr{
													Op: "",
													Value: arrange.Value{
														Desc: "req1.code /eq 0",
														Type: "exprcmp",
														Value: arrange.FlexJsonObject{
															D: arrange.ExprCmp{
																LValue: arrange.Value{
																	Type: "field",
																	Value: arrange.FlexJsonObject{
																		D: "cond1.yes",
																	},
																},
																Cmp: "eq",
																RValue: arrange.Value{
																	Type: "boolean",
																	Value: arrange.FlexJsonObject{
																		D: "true",
																	},
																},
															},
														},
													},
												},
												Yes: arrange.Value{
													Type: "field",
													Value: arrange.FlexJsonObject{
														D: "req2",
													},
												},
												No: arrange.Value{
													Type: "object",
													Value: arrange.FlexJsonObject{
														D: arrange.ValObject{},
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
*/
