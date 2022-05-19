package service

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/quanxiang-cloud/polyapi/internal/arrange"
	"github.com/quanxiang-cloud/polyapi/pkg/misc/config"
	"testing"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkflowSuiteQingCloud struct {
	suite.Suite
	ctx     context.Context
	conf    *config.Config
	engin   *gin.Engine
	rawAPI  RawAPI
	polyAPI PolyAPI
	ids     struct {
		raw     []string
		poly    string
		update  bool
		request bool
	}
}

func _TestWorkflowQingCloud(t *testing.T) {
	s := &WorkflowSuiteQingCloud{}
	switch testEnv {
	case "local":
		// show stop start create
		s.ids.raw = []string{
			"BbJ5bHzVtU9HjuGqluGW8L3eM2V52Sxm0UiFMncVEqk", // show
			"GzFWtDHH3kC28FTHArz7DRJg0UzxPNPjl787yWgMinc", // stop
			"otyHwCFAHyhjne9OcqJsW8f3M_0riB6SJe1hywoLwbQ", // start
			"3Fe78hFqi380S1lzfyn2xdzG0i-TqGXzVeEziJeVuYA", // create
		}
		s.ids.poly = "2cxYzbBQF0gwua-MTNhK3e7d1A-2-08aihT8tI0tijM"
		s.ids.update = false
		s.ids.request = true
	case "debug":
	case "test":
	}

	suite.Run(t, s)
}

func (s *WorkflowSuiteQingCloud) SetupSuite() {
	var err error
	s.ctx = logger.GenRequestID(context.TODO())
	s.conf, err = config.NewConfig(fmt.Sprintf("./testdata/%s/polyapi.yaml", testEnv))
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.conf)
	s.conf.Log.Level = 1 // warn
	err = logger.New(&s.conf.Log)
	assert.Nil(s.T(), err)
	s.rawAPI, err = CreateRaw(s.conf)
	s.polyAPI, err = CreatePoly(s.conf)
}

func (s *WorkflowSuiteQingCloud) TestAll() {
	if len(s.ids.raw) == 0 { // create raw
		s.testCreateRaw()
		return
	}
	s.testGetRaw()
	if s.ids.poly == "" { // create poly
		s.testCreatePoly()

		return
	}
	if s.ids.update {
		s.testUpdatePoly()
		s.testUpdatePolyScript()
	}

	s.testBuildPoly()
	//s.testRequestPoly()
}

func (s *WorkflowSuiteQingCloud) testCreateRaw() {
	file, _ := os.Open(fmt.Sprintf("./testdata/%s/swagger_qingcloud.json", testEnv))
	all, _ := io.ReadAll(file)
	req := &RegReq{
		Swagger: string(all),
		Host:    "api.qingcloud.com.unuse",
		Version: "1.0",
	}
	resp, err := s.rawAPI.RegSwagger(s.ctx, req)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteQingCloud) testGetRaw() {
	req := &QueryReq{
		ID: s.ids.raw[0],
	}
	resp, err := s.rawAPI.Query(s.ctx, req)
	fmt.Printf("get raw %+v %v\n", resp, err)
	fmt.Printf("%+v\n", resp.Content)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteQingCloud) testCreatePoly() {
	req := &PolyCreateReq{
		Parent: "testParent",
		Name:   "powerOn_Off",
	}
	resp, err := s.polyAPI.Create(s.ctx, req)
	fmt.Printf("create %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPoly(resp.ID)
}

func (s *WorkflowSuiteQingCloud) testUpdatePoly() {
	req := &PolyUpdateArrangeReq{
		ID:      s.ids.poly,
		Arrange: "arrange V1",
	}
	resp, err := s.polyAPI.UpdateArrange(s.ctx, req)
	fmt.Printf("update V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPoly("")
}

func (s *WorkflowSuiteQingCloud) testGetPoly(id string) {
	if id == "" {
		id = s.ids.poly
	}
	req := &PolyGetArrangeReq{ID: id}
	resp, err := s.polyAPI.GetArrange(s.ctx, req)
	fmt.Printf("get %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteQingCloud) testUpdatePolyScript() {
	req := &PolyUpdateScriptReq{
		ID:     s.ids.poly,
		Script: `"script V1"`,
	}
	resp, err := s.polyAPI.UpdateScript(s.ctx, req)
	fmt.Printf("update script V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPolyScript()
	s.testRequestPoly()
}

func (s *WorkflowSuiteQingCloud) testBuildPoly() {
	req := &PolyBuildReq{
		ID:      s.ids.poly,
		Arrange: s.genSampleArrange(),
	}
	resp, err := s.polyAPI.Build(s.ctx, req)
	fmt.Printf("build V1 %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
	s.testGetPolyScript()
	if s.ids.request {
		s.testRequestPoly()
	}
}

func (s *WorkflowSuiteQingCloud) testGetPolyScript() {
	req := &PolyGetScriptReq{ID: s.ids.poly}
	resp, err := s.polyAPI.GetScript(s.ctx, req)
	fmt.Printf("get script %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteQingCloud) testRequestPoly() {
	f, err := os.Open(fmt.Sprintf("./testdata/%s/polytest_qc.testdata", testEnv))
	if err != nil {
		s.T().Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		s.T().Fatal(err)
	}

	req := &PolyRequestReq{
		ID:   s.ids.poly,
		Body: b,
		Header: http.Header{
			"User-Id": []string{"893ca81d-f571-4a6f-8088-673e8775ff64"},
		},
	}
	resp, err := s.polyAPI.Request(s.ctx, req)
	fmt.Printf("request %+v %v\n", resp, err)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), resp)
}

func (s *WorkflowSuiteQingCloud) genSampleArrange() string {
	a := &arrange.Arrange{
		ID:      s.ids.poly,
		Name:    "polyApi",
		Desc:    "polyApi sample",
		Version: "v0.0.3",
		Nodes: []arrange.Node{
			arrange.Node{
				Name:      "start",
				Type:      "input",
				Desc:      "the only entrance of the arrange",
				NextNodes: []string{"req1"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 1, Y: 10}},
				Detail: arrange.FlexJSONObject{
					D: arrange.InputNodeDetail{
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Type: "array_string",
								Name: "test_array_string",
								Data: arrange.FlexJSONObject{
									D: arrange.ValArrayString{"foo", "bar"},
								},
								In: "hide",
							},
						},
						Defines: []arrange.ValueDefine{
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									Name: "h1In",
									In:   "header",
									Type: "string",
									Desc: "h1 input from header",
								},
								Default: "xx",
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									Name: "p1",
									In:   "path",
									Type: "string",
									Desc: "p1 input from path",
								},
								Default: "xx",
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "a",
									Type: "string",
								},
								Default:  "xx",
								Enums:    []string{"x", "y"},
								Mock:     "aa",
								Required: true,
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "b",
									Type: "number",
								},
								Default:  "xx",
								Enums:    []string{"1", "2"},
								Mock:     "0",
								Required: true,
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "c",
									Type: "object",
									Data: arrange.FlexJSONObject{
										D: &arrange.ValObject{
											arrange.Value{
												Name: "a",
												Type: "string",
											},
											arrange.Value{
												Name: "b",
												Type: "number",
											},
											arrange.Value{
												Name: "b",
												Type: "array",
												Data: arrange.FlexJSONObject{
													D: &arrange.ValArray{
														arrange.Value{
															Type: "object",
															Data: arrange.FlexJSONObject{
																D: &arrange.ValObject{
																	arrange.Value{
																		Name: "a",
																		Type: "string",
																	},
																	arrange.Value{
																		Name: "b",
																		Type: "number",
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
					},
				},
			},
			arrange.Node{
				Name:      "req1",
				Type:      "request",
				Desc:      "show vm list",
				NextNodes: []string{"cond1", "end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 10}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawID: s.ids.raw[0],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Type:  "skey",
								Field: "start.security_key",
								In:    "hide",
							},
							arrange.InputValue{
								Name:  "zone",
								Type:  "string",
								Field: "start.zone",
							},
							arrange.InputValue{
								Name:  "access_key_id",
								Type:  "string",
								Field: "start.access_key_id",
							},

							arrange.InputValue{
								Type: "array_string_elem",
								Data: arrange.FlexJSONObject{
									D: arrange.ValArrayStringElem{
										Name: "status",
										Array: arrange.ValArrayString{
											"pending", "running", "stopped",
										},
									},
								},
							},
							arrange.InputValue{
								Name: "verbose",
								Type: "number",
								Data: arrange.FlexJSONObject{
									D: "1",
								},
							},
						},
					},
				},
			},

			arrange.Node{
				Name:      "cond1",
				Type:      "if",
				Desc:      "check if has vms",
				NextNodes: []string{},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.IfNodeDetail{
						Yes: "cond2",
						No:  "req4",
						Cond: arrange.CondExpr{
							Op: "",
							Value: arrange.Value{
								Desc: "req1.total_count /gt 0",
								Type: "exprcmp",
								Data: arrange.FlexJSONObject{
									D: arrange.ValExprCmp{
										LValue: arrange.Value{
											Type:  "number",
											Field: "req1.total_count",
										},
										Cmp: "gt",
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
				Name:      "cond2",
				Type:      "if",
				Desc:      "check if vm[0] is running",
				NextNodes: []string{},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.IfNodeDetail{
						Yes: "req2",
						No:  "cond3",
						Cond: arrange.CondExpr{
							Op: "",
							Value: arrange.Value{
								Desc: "req1.instance_set[0].status /eq running",
								Type: "exprcmp",
								Data: arrange.FlexJSONObject{
									D: arrange.ValExprCmp{
										LValue: arrange.Value{
											Type:  "string",
											Field: "req1.instance_set[0].status",
										},
										Cmp: "eq",
										RValue: arrange.Value{
											Type: "string",
											Data: arrange.FlexJSONObject{
												D: "running",
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
				Name:      "cond3",
				Type:      "if",
				Desc:      "check if vm[0] is stopped",
				NextNodes: []string{},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.IfNodeDetail{
						Yes: "req3",
						No:  "",
						Cond: arrange.CondExpr{
							Op: "",
							Value: arrange.Value{
								Desc: "req1.instance_set[0].status /eq stopped",
								Type: "exprcmp",
								Data: arrange.FlexJSONObject{
									D: arrange.ValExprCmp{
										LValue: arrange.Value{
											Type:  "string",
											Field: "req1.instance_set[0].status",
										},
										Cmp: "eq",
										RValue: arrange.Value{
											Type: "string",
											Data: arrange.FlexJSONObject{
												D: "stopped",
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
				Desc:      "stop instance",
				NextNodes: []string{"end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawID: s.ids.raw[1],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Type:  "skey",
								Field: "start.security_key",
								In:    "hide",
							},
							arrange.InputValue{
								Name:  "zone",
								Type:  "string",
								Field: "start.zone",
							},
							arrange.InputValue{
								Name:  "access_key_id",
								Type:  "string",
								Field: "start.access_key_id",
							},
							arrange.InputValue{
								Name:  "instances.1",
								Type:  "string",
								Field: "req1.instance_set[0].instance_id",
							},

							arrange.InputValue{
								Name: "force",
								Type: "number",
								Data: arrange.FlexJSONObject{
									D: "1",
								},
							},
						},
					},
				},
			},
			arrange.Node{
				Name:      "req3",
				Type:      "request",
				Desc:      "start instance",
				NextNodes: []string{"end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawID: s.ids.raw[2],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Type:  "skey",
								Field: "start.security_key",
								In:    "hide",
							},
							arrange.InputValue{
								Name:  "zone",
								Type:  "string",
								Field: "start.zone",
							},
							arrange.InputValue{
								Name:  "access_key_id",
								Type:  "string",
								Field: "start.access_key_id",
							},

							arrange.InputValue{
								Name:  "instances.1",
								Type:  "string",
								Field: "req1.instance_set[0].instance_id",
							},
						},
					},
				},
			},
			arrange.Node{
				Name:      "req4",
				Type:      "request",
				Desc:      "create instance",
				NextNodes: []string{"end"},
				//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
				Detail: arrange.FlexJSONObject{
					D: arrange.RequestNodeDetail{
						RawID: s.ids.raw[3],
						Inputs: arrange.ValueSet{
							arrange.InputValue{
								Type:  "skey",
								Field: "start.security_key",
								In:    "hide",
							},
							arrange.InputValue{
								Name:  "zone",
								Type:  "field",
								Field: "start.zone",
							},
							arrange.InputValue{
								Name:  "access_key_id",
								Type:  "field",
								Field: "start.access_key_id",
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
										Field: "req1",
									},
									arrange.Value{
										Type:  "field",
										Field: "req2",
									},
									arrange.Value{
										Type:  "field",
										Name:  "req3x",
										Field: "req3",
									},
									arrange.Value{
										Type:  "field",
										Name:  "req4x",
										Field: "req4",
									},
									arrange.Value{
										Name: "merged",
										Type: "mergeobj",
										Data: arrange.FlexJSONObject{
											D: arrange.ValMergeObj{
												arrange.Value{
													Type: "filter",
													Data: arrange.FlexJSONObject{
														D: arrange.ValFilter{
															Source: "req1",
															White: arrange.FieldMap{
																"total_count": "InstNummmmmmmmmmmm",
																"ret_code":    "CodeReq11111111111111",
															},
															Black: arrange.FieldMap{},
														},
													},
												},
												arrange.Value{
													Type: "filter",
													Data: arrange.FlexJSONObject{
														D: arrange.ValFilter{
															Source: "req2",
															White:  arrange.FieldMap{},
															Black: arrange.FieldMap{
																"None": "",
															},
														},
													},
												},
												arrange.Value{
													Type: "filter",
													Data: arrange.FlexJSONObject{
														D: arrange.ValFilter{
															Source: "req3",
															White:  arrange.FieldMap{},
															Black: arrange.FieldMap{
																"None": "",
															},
														},
													},
												},
												arrange.Value{
													Type: "filter",
													Data: arrange.FlexJSONObject{
														D: arrange.ValFilter{
															Source: "req4",
															White:  arrange.FieldMap{},
															Black: arrange.FieldMap{
																"None": "",
															},
														},
													},
												},
												arrange.Value{
													Type: "exprsel",
													Data: arrange.FlexJSONObject{
														D: arrange.ValExprSel{
															Cond: arrange.CondExpr{
																Op: "",
																Value: arrange.Value{
																	Desc: "req1.code /eq 0",
																	Type: "exprcmp",
																	Data: arrange.FlexJSONObject{
																		D: arrange.ValExprCmp{
																			LValue: arrange.Value{
																				Type:  "field",
																				Field: "cond1.yes",
																			},
																			Cmp: "eq",
																			RValue: arrange.Value{
																				Type: "boolean",
																				Data: arrange.FlexJSONObject{
																					D: "true",
																				},
																			},
																		},
																	},
																},
															},
															Yes: arrange.Value{
																Type:  "field",
																Field: "req2",
															},
															No: arrange.Value{
																Type: "object",
																Data: arrange.FlexJSONObject{
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
						Doc: []arrange.ValueDefine{
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "header",
									Type: "string",
								},
								Default: "xx",
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "a",
									Type: "string",
								},
								Default:  "xx",
								Enums:    []string{"x", "y"},
								Mock:     "aa",
								Required: true,
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "b",
									Type: "number",
								},
								Default:  "xx",
								Enums:    []string{"1", "2"},
								Mock:     "0",
								Required: true,
							},
							arrange.ValueDefine{
								InputValue: arrange.InputValue{
									In:   "body",
									Name: "c",
									Type: "object",
									Data: arrange.FlexJSONObject{
										D: &arrange.ValObject{
											arrange.Value{
												Name: "a",
												Type: "string",
											},
											arrange.Value{
												Name: "b",
												Type: "number",
											},
											arrange.Value{
												Name: "b",
												Type: "array",
												Data: arrange.FlexJSONObject{
													D: &arrange.ValArray{
														arrange.Value{
															Type: "object",
															Data: arrange.FlexJSONObject{
																D: &arrange.ValObject{
																	arrange.Value{
																		Name: "a",
																		Type: "string",
																	},
																	arrange.Value{
																		Name: "b",
																		Type: "number",
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
						Header: arrange.ValueSet{
							arrange.InputValue{
								In:   "header",
								Name: "h1",
								Type: "header",
								Desc: "header h1",
							},
							arrange.InputValue{
								In:   "header",
								Name: "h2",
								Type: "string",
								Desc: "header h2",
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
