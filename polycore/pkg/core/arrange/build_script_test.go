package arrange_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"

	// NOTE: init poly doc generator
	_ "github.com/quanxiang-cloud/polyapi/polycore/pkg/core/polydoc"
)

const testCfgGenArrange = false
const testCfgArrange = false

func TestBuildScript(t *testing.T) {
	old := adaptor.SetRawAPIOper(dummyRawAPI(0))
	defer func() { adaptor.SetRawAPIOper(old) }()

	testGenerateArrange(t)
	testBuildArrange(t)
}

func testBuildArrange(t *testing.T) {
	for _, v := range sampleArranges[len(sampleArranges)-1:] {
		script, doc, _, err := arrange.BuildJsScript(nil, v, "")
		if err != nil {
			t.Errorf("testBuildArrange error:%s", err)
		}
		if !testCfgGenArrange && testCfgArrange {
			fmt.Println(script)
			fmt.Println(doc)
		}
	}
}

func _testGenArrangeJson(d interface{}) (string, error) {
	b, err := json.MarshalIndent(d, "", "  ")
	//b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func testGenerateArrange(t *testing.T) string {
	rawPath := "/system/raw/sample"

	arg := &arrange.Arrange{}
	if _, err := arrange.InitArrange("", arg); err != nil {
		t.Fatalf("InitArrange fail: %s", err)
	}

	var exports = []interface{}{
		&arrange.Arrange{
			Info: &arrange.APIInfo{
				Name: "polyApi",
				Desc: "polyApi sample",
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
							Inputs: []arrange.ValueDefine{
								arrange.ValueDefine{
									InputValue: arrange.InputValue{
										Name: "Access-Token",
										Type: "string",
										Desc: "cookie token",
										Data: arrange.FlexJSONObject{
											D: "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5",
										},
										In:       "header",
										Required: true,
									},
									Mock: "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5",
								},
								arrange.ValueDefine{
									InputValue: arrange.InputValue{
										Name: "Content-Type",
										Type: "string",
										Desc: "encoding method",
										Data: arrange.FlexJSONObject{
											D: "application/json",
										},
										In:       "header",
										Required: true,
									},
									Mock: "",
								},
								arrange.ValueDefine{
									InputValue: arrange.InputValue{
										Name: "userName",
										Type: "string",
										Desc: "encoding",
										Data: arrange.FlexJSONObject{
											D: "",
										},
										In:       "body",
										Required: true,
									},
									Mock: "testUser",
								},
								arrange.ValueDefine{
									InputValue: arrange.InputValue{
										Name: "xx",
										Type: "string",
										Desc: "xxxx",
										Data: arrange.FlexJSONObject{
											D: "foo",
										},
										In:       "body",
										Required: true,
									},
									Mock: "xxx",
								},
							},
							Consts: arrange.ValueSet{
								arrange.InputValue{
									Name: "x",
									Type: "string",
									Desc: "xxxx",
									Data: arrange.FlexJSONObject{
										D: "foo",
									},
									In: "body",
								},
								arrange.InputValue{
									Name: "y",
									Type: "string",
									Desc: "xxxx",
									Data: arrange.FlexJSONObject{
										D: "foo",
									},
									In: "header",
								},
								arrange.InputValue{
									Name: "z",
									Type: "string",
									Desc: "xxxx",
									Data: arrange.FlexJSONObject{
										D: "foo",
									},
									In: "path",
								},
							},
						},
					},
				},
				arrange.Node{
					Name:      "req1",
					Type:      "request",
					Desc:      "request1",
					NextNodes: []string{"cond1", "end"},
					//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 10}},
					Detail: arrange.FlexJSONObject{
						D: arrange.RequestNodeDetail{
							RawPath: rawPath,
							Inputs: arrange.ValueSet{
								arrange.InputValue{
									Name: "Access-Token",
									Type: "field",
									Data: arrange.FlexJSONObject{D: "start.Access-Token"},
									In:   "header",
								},
								arrange.InputValue{
									Name: "userName",
									Type: "field",
									Data: arrange.FlexJSONObject{D: "start.userName"},
									In:   "body",
								},
							},
						},
					},
				},
				arrange.Node{
					Name:      "cond1",
					Type:      "if",
					Desc:      "check if need call request2",
					NextNodes: []string{},
					//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 5, Y: 20}},
					Detail: arrange.FlexJSONObject{
						D: arrange.IfNodeDetail{
							Yes: "req2",
							No:  "req3",
							Cond: arrange.CondExpr{
								Op: "",
								Value: arrange.Value{
									Type: "exprgroup",
									Desc: `(req1.data.id /lt 10) /and ((req1.data.x+req1.data.y) /eq 50)`,
									Data: arrange.FlexJSONObject{
										D: arrange.CondExprGroup{
											arrange.CondExpr{
												Op: "",
												Value: arrange.Value{
													Desc: "req1.data.id /lt 10",
													Type: "exprcmp",
													Data: arrange.FlexJSONObject{
														D: arrange.ValExprCmp{
															LValue: arrange.Value{
																Type: "field",
																Data: arrange.FlexJSONObject{
																	D: "req1.data.id",
																},
															},
															Cmp: "lt",
															RValue: arrange.Value{
																Type: "number",
																Data: arrange.FlexJSONObject{
																	D: "10",
																},
															},
														},
													},
												},
											},
											arrange.CondExpr{
												Op: "and",
												Value: arrange.Value{
													Desc: "(req1.data.x+req1.data.y) /eq 50",
													Type: "exprcmp",
													Data: arrange.FlexJSONObject{
														D: arrange.ValExprCmp{
															LValue: arrange.Value{
																Type: "exprgroup",
																Data: arrange.FlexJSONObject{
																	D: arrange.ValExprGroup{
																		arrange.ValExpr{
																			Op: "",
																			Value: arrange.Value{
																				Type: "field",
																				Data: arrange.FlexJSONObject{
																					D: "req1.data.x",
																				},
																			},
																		},
																		arrange.ValExpr{
																			Op: "add",
																			Value: arrange.Value{
																				Type: "field",
																				Data: arrange.FlexJSONObject{
																					D: "req1.data.y",
																				},
																			},
																		},
																	},
																},
															},
															Cmp: "eq",
															RValue: arrange.Value{
																Type: "number",
																Data: arrange.FlexJSONObject{
																	D: "50",
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
					Name:      "req2",
					Type:      "request",
					Desc:      "request2",
					NextNodes: []string{"end"},
					//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
					Detail: arrange.FlexJSONObject{
						D: arrange.RequestNodeDetail{
							RawPath: rawPath,
							Inputs: []arrange.InputValue{
								arrange.InputValue{
									Name: "Access-Token",
									Type: "field",
									Data: arrange.FlexJSONObject{D: "start.Access-Token"},
									In:   "header",
								},
								arrange.InputValue{
									Name: "userName",
									Type: "field",
									Data: arrange.FlexJSONObject{D: "start.userName"},
									In:   "header",
								},
							},
						},
					},
				},
				arrange.Node{
					Name:      "req3",
					Type:      "request",
					Desc:      "request2",
					NextNodes: []string{"end"},
					//Disp:      arrange.FlexJSONObject{D: arrange.DispDemo{X: 10, Y: 20}},
					Detail: arrange.FlexJSONObject{
						D: arrange.RequestNodeDetail{
							RawPath: rawPath,
							Inputs:  []arrange.InputValue{},
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
											Name: "code",
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
																		Type: "field",
																		Data: arrange.FlexJSONObject{
																			D: "req1.code",
																		},
																	},
																	Cmp: "lt",
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
													Yes: arrange.Value{
														Type: "number",
														Data: arrange.FlexJSONObject{
															D: "0",
														},
													},
													No: arrange.Value{
														Type: "number",
														Data: arrange.FlexJSONObject{
															D: "1",
														},
													},
												},
											},
										},
										arrange.Value{
											Name: "data",
											Type: "mergeobj",
											Data: arrange.FlexJSONObject{
												D: arrange.ValMergeObj{
													arrange.Value{
														Type: "filter",
														Data: arrange.FlexJSONObject{
															D: arrange.ValFiltObj{
																Source: "req1.data",
																White: arrange.FieldMap{
																	"field1": "field1_n",
																	"field2": "",
																},
															},
														},
													},
													arrange.Value{
														Type: "object",
														Data: arrange.FlexJSONObject{
															D: arrange.ValObject{
																arrange.Value{
																	Name: "field3",
																	Type: "field",
																	Data: arrange.FlexJSONObject{
																		D: "req2.data.x",
																	},
																},
																arrange.Value{
																	Name: "time_stamp",
																	Type: "exprfunc0",
																	Data: arrange.FlexJSONObject{
																		D: arrange.ValExprFunc{
																			Func: "timestamp",
																		},
																	},
																},
																arrange.Value{
																	Name: "field5",
																	Type: "filter",
																	Data: arrange.FlexJSONObject{
																		D: arrange.ValFiltObj{
																			Source: "req2.data.aryField5",
																			White: arrange.FieldMap{
																				"field51": "field51_n",
																				"field52": "",
																			},
																			Filter: arrange.CondExpr{
																				Value: arrange.Value{
																					Desc: "id /lt 10",
																					Type: "exprcmp",
																					Data: arrange.FlexJSONObject{
																						D: arrange.ValExprCmp{
																							LValue: arrange.Value{
																								Type: "field",
																								Data: arrange.FlexJSONObject{
																									D: "id",
																								},
																							},
																							Cmp: "lt",
																							RValue: arrange.Value{
																								Type: "number",
																								Data: arrange.FlexJSONObject{
																									D: "10",
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
									},
								},
							},
							Header: arrange.ValueSet{
								arrange.InputValue{
									Name: "cookie1",
									Type: "field",
									Data: arrange.FlexJSONObject{
										D: "req1.cookie1",
									},
									In: "header",
								},
							},
						},
					},
				},
			},
		},
	}

	for i := 0; i < len(exports); i++ {
		txt, err := _testGenArrangeJson(exports[i])
		if err != nil {
			t.Errorf("case %d error:%s", i+1, err)
		}
		if txt != sampleArranges[0] {
			fmt.Printf("case%d:---------------\n", i+1)
			fmt.Printf("@%q@\n", sampleArranges[0])
			fmt.Printf("@%q@\n", txt)
			fmt.Printf("@%s@\n", txt)
			t.Errorf("case %d build result mismatch", i+1)
			continue
		}
		if testCfgArrange {
			fmt.Printf("case%d:---------------\n", i+1)
			fmt.Printf("@%s@\n", txt)
		}
		return txt
	}
	return ""
}

var errNotSupport = errors.New("not support")

type dummyRawAPI int

func (r dummyRawAPI) Query(c context.Context, req *adaptor.QueryRawAPIReq) (*adaptor.QueryRawAPIResp, error) {
	return &adaptor.QueryRawAPIResp{
		ID:     "xx",
		URL:    "http://sample.com/api/v1/request",
		Schema: "http",
		Method: "GET",
		Content: &adaptor.RawAPIContent{
			EncodingIn:  "json",
			EncodingOut: "json",

			Method: "GET",
		},
	}, nil
}
func (r dummyRawAPI) List(c context.Context, req *adaptor.RawListReq) (*adaptor.RawListResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) ListInService(c context.Context, req *adaptor.ListInServiceReq) (*adaptor.ListInServiceResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) InnerUpdateRawInBatch(ctx context.Context, req *adaptor.InnerUpdateRawInBatchReq) (*adaptor.InnerUpdateRawInBatchResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) InnerDel(c context.Context, req *adaptor.DelReq) (*adaptor.DelResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) ValidInBatches(ctx context.Context, req *adaptor.RawValidInBatchesReq) (*adaptor.RawValidInBatchesResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) ValidByPrefixPath(ctx context.Context, req *adaptor.RawValidByPrefixPathReq) (*adaptor.RawValidByPrefixPathResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) QueryInBatches(c context.Context, req *adaptor.QueryRawAPIInBatchesReq) (*adaptor.QueryRawAPIInBatchesResp, error) {
	return nil, errNotSupport
}

func (r dummyRawAPI) ListByPrefixPath(ctx context.Context, req *adaptor.ListRawByPrefixPathReq) (*adaptor.ListRawByPrefixPathResp, error) {
	return nil, errNotSupport
}

func (r dummyRawAPI) InnerImport(ctx context.Context, req *adaptor.InnerImportRawReq) (*adaptor.InnerImportRawResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) InnerDelByPrefixPath(ctx context.Context, req *adaptor.InnerDelRawByPrefixPathReq) (*adaptor.InnerDelRawByPrefixPathResp, error) {
	return nil, errNotSupport
}
func (r dummyRawAPI) QueryDoc(ctx context.Context, req *apiprovider.QueryDocReq) (*apiprovider.QueryDocResp, error) {
	return &apiprovider.QueryDocResp{
		Doc: json.RawMessage(`{}`),
	}, nil
}
func (r dummyRawAPI) QuerySwagger(ctx context.Context, req *adaptor.QueryRawSwaggerReq) (*adaptor.QueryRawSwaggerResp, error) {
	return nil, errNotSupport
}

var sampleArranges = []string{`{
  "info": {
    "name": "polyApi",
    "desc": "polyApi sample"
  },
  "nodes": [
    {
      "name": "start",
      "title": "",
      "desc": "the only entrance of the arrange",
      "type": "input",
      "nextNodes": [
        "req1"
      ],
      "detail": {
        "inputs": [
          {
            "type": "string",
            "name": "Access-Token",
            "desc": "cookie token",
            "required": true,
            "data": "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5",
            "in": "header",
            "mock": "Bear M2MWZDKZNJQTYTM5YY01MJLMLWE4MGQTMJE3MTLKYMUYNMI5"
          },
          {
            "type": "string",
            "name": "Content-Type",
            "desc": "encoding method",
            "required": true,
            "data": "application/json",
            "in": "header"
          },
          {
            "type": "string",
            "name": "userName",
            "desc": "encoding",
            "required": true,
            "data": "",
            "in": "body",
            "mock": "testUser"
          },
          {
            "type": "string",
            "name": "xx",
            "desc": "xxxx",
            "required": true,
            "data": "foo",
            "in": "body",
            "mock": "xxx"
          }
        ],
        "consts": [
          {
            "type": "string",
            "name": "x",
            "desc": "xxxx",
            "data": "foo",
            "in": "body"
          },
          {
            "type": "string",
            "name": "y",
            "desc": "xxxx",
            "data": "foo",
            "in": "header"
          },
          {
            "type": "string",
            "name": "z",
            "desc": "xxxx",
            "data": "foo",
            "in": "path"
          }
        ]
      }
    },
    {
      "name": "req1",
      "title": "",
      "desc": "request1",
      "type": "request",
      "nextNodes": [
        "cond1",
        "end"
      ],
      "detail": {
        "rawPath": "/system/raw/sample",
        "apiKeyID": "",
        "dynamicKey": false,
        "inputs": [
          {
            "type": "field",
            "name": "Access-Token",
            "data": "start.Access-Token",
            "in": "header"
          },
          {
            "type": "field",
            "name": "userName",
            "data": "start.userName",
            "in": "body"
          }
        ]
      }
    },
    {
      "name": "cond1",
      "title": "",
      "desc": "check if need call request2",
      "type": "if",
      "nextNodes": [],
      "detail": {
        "cond": {
          "op": "",
          "type": "exprgroup",
          "name": "",
          "desc": "(req1.data.id /lt 10) /and ((req1.data.x+req1.data.y) /eq 50)",
          "data": [
            {
              "op": "",
              "type": "exprcmp",
              "name": "",
              "desc": "req1.data.id /lt 10",
              "data": {
                "lvalue": {
                  "type": "field",
                  "name": "",
                  "data": "req1.data.id"
                },
                "cmp": "lt",
                "rvalue": {
                  "type": "number",
                  "name": "",
                  "data": "10"
                }
              }
            },
            {
              "op": "and",
              "type": "exprcmp",
              "name": "",
              "desc": "(req1.data.x+req1.data.y) /eq 50",
              "data": {
                "lvalue": {
                  "type": "exprgroup",
                  "name": "",
                  "data": [
                    {
                      "op": "",
                      "type": "field",
                      "name": "",
                      "data": "req1.data.x"
                    },
                    {
                      "op": "add",
                      "type": "field",
                      "name": "",
                      "data": "req1.data.y"
                    }
                  ]
                },
                "cmp": "eq",
                "rvalue": {
                  "type": "number",
                  "name": "",
                  "data": "50"
                }
              }
            }
          ]
        },
        "yes": "req2",
        "no": "req3"
      }
    },
    {
      "name": "req2",
      "title": "",
      "desc": "request2",
      "type": "request",
      "nextNodes": [
        "end"
      ],
      "detail": {
        "rawPath": "/system/raw/sample",
        "apiKeyID": "",
        "dynamicKey": false,
        "inputs": [
          {
            "type": "field",
            "name": "Access-Token",
            "data": "start.Access-Token",
            "in": "header"
          },
          {
            "type": "field",
            "name": "userName",
            "data": "start.userName",
            "in": "header"
          }
        ]
      }
    },
    {
      "name": "req3",
      "title": "",
      "desc": "request2",
      "type": "request",
      "nextNodes": [
        "end"
      ],
      "detail": {
        "rawPath": "/system/raw/sample",
        "apiKeyID": "",
        "dynamicKey": false,
        "inputs": []
      }
    },
    {
      "name": "end",
      "title": "",
      "desc": "the only end of this arrange",
      "type": "output",
      "nextNodes": null,
      "detail": {
        "header": [
          {
            "type": "field",
            "name": "cookie1",
            "data": "req1.cookie1",
            "in": "header"
          }
        ],
        "body": {
          "type": "object",
          "name": "",
          "data": [
            {
              "type": "exprsel",
              "name": "code",
              "data": {
                "cond": {
                  "op": "",
                  "type": "exprcmp",
                  "name": "",
                  "desc": "req1.code /eq 0",
                  "data": {
                    "lvalue": {
                      "type": "field",
                      "name": "",
                      "data": "req1.code"
                    },
                    "cmp": "lt",
                    "rvalue": {
                      "type": "number",
                      "name": "",
                      "data": "0"
                    }
                  }
                },
                "yes": {
                  "type": "number",
                  "name": "",
                  "data": "0"
                },
                "no": {
                  "type": "number",
                  "name": "",
                  "data": "1"
                }
              }
            },
            {
              "type": "mergeobj",
              "name": "data",
              "data": [
                {
                  "type": "filter",
                  "name": "",
                  "data": {
                    "source": "req1.data",
                    "white": {
                      "field1": "field1_n",
                      "field2": ""
                    },
                    "black": null,
                    "filter": {
                      "op": "",
                      "type": "",
                      "name": "",
                      "data": null
                    }
                  }
                },
                {
                  "type": "object",
                  "name": "",
                  "data": [
                    {
                      "type": "field",
                      "name": "field3",
                      "data": "req2.data.x"
                    },
                    {
                      "type": "exprfunc0",
                      "name": "time_stamp",
                      "data": {
                        "func": "timestamp",
                        "paras": null
                      }
                    },
                    {
                      "type": "filter",
                      "name": "field5",
                      "data": {
                        "source": "req2.data.aryField5",
                        "white": {
                          "field51": "field51_n",
                          "field52": ""
                        },
                        "black": null,
                        "filter": {
                          "op": "",
                          "type": "exprcmp",
                          "name": "",
                          "desc": "id /lt 10",
                          "data": {
                            "lvalue": {
                              "type": "field",
                              "name": "",
                              "data": "id"
                            },
                            "cmp": "lt",
                            "rvalue": {
                              "type": "number",
                              "name": "",
                              "data": "10"
                            }
                          }
                        }
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      }
    }
  ]
}`,
	`
{
    "nodes": [
        {
            "title": "",
            "name": "start",
            "type": "input",
            "nextNodes": [
                "J8TGwtOX",
                "KwLhv8xp"
            ],
            "detail": {
                "inputs": [
                    {
                        "type": "object",
                        "name": "body_ip",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "ip",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "descPath": "开始节点.body_ip.ip"
                            },
                            {
                                "type": "string",
                                "name": "key",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "descPath": "开始节点.body_ip.key"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "descPath": "开始节点.body_ip"
                    },
                    {
                        "type": "object",
                        "name": "body_tianqi",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "key",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "descPath": "开始节点.body_tianqi.key"
                            },
                            {
                                "type": "string",
                                "name": "city",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "descPath": "开始节点.body_tianqi.city"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "descPath": "开始节点.body_tianqi"
                    },
                    {
                        "type": "object",
                        "name": "body_fanyi",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "q",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "descPath": "开始节点.body_fanyi.q"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "descPath": "开始节点.body_fanyi"
                    }
                ],
                "consts": []
            }
        },
        {
            "title": "ip",
            "name": "J8TGwtOX",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/rtftw/raw/customer/zhuji/ip",
                "apiName": "ip",
                "inputs": [
                    {
                        "type": "string",
                        "name": "Content-Type",
                        "title": "数据格式",
                        "desc": "application/json",
                        "data": "application/json",
                        "in": "header",
                        "required": true,
                        "mock": "application/json"
                    },
                    {
                        "type": "direct_expr",
                        "name": "ip",
                        "data": "start.body_ip.ip ",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "direct_expr",
                        "name": "key",
                        "data": "start.body_ip.key ",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "object",
                        "data": [
                        ],
                        "in": "body"
                    }
                ],
                "outputs": []
            }
        },
        {
            "title": "",
            "name": "end",
            "type": "output",
            "nextNodes": [],
            "detail": {
                "body": {
                    "type": "object",
                    "data": [
                        {
                            "type": "object",
                            "name": "code",
                            "desc": "",
                            "data": [
                                {
                                    "type": "direct_expr",
                                    "name": "data1",
                                    "desc": "",
                                    "data": "$KwLhv8xp",
                                    "in": "body",
                                    "required": false
                                },
                                {
                                    "type": "direct_expr",
                                    "name": "data2",
                                    "desc": "",
                                    "data": "$DyI3fAv7",
                                    "in": "body",
                                    "required": false
                                },
                                {
                                    "type": "direct_expr",
                                    "name": "data3",
                                    "desc": "",
                                    "data": "$J8TGwtOX",
                                    "in": "body",
                                    "required": false
                                }
                            ],
                            "in": "body",
                            "required": false
                        }
                    ]
                }
            }
        },
        {
            "title": "fanyi",
            "name": "KwLhv8xp",
            "type": "request",
            "nextNodes": [
                "DyI3fAv7"
            ],
            "detail": {
                "rawPath": "/system/app/rtftw/raw/customer/demo/kkk",
                "apiName": "kkk",
                "inputs": [
                    {
                        "type": "string",
                        "name": "Content-Type",
                        "title": "数据格式",
                        "desc": "application/json",
                        "data": "application/json",
                        "in": "header",
                        "required": true,
                        "mock": "application/json"
                    },
                    {
                        "type": "direct_expr",
                        "name": "q",
                        "data": "start.body_fanyi.q ",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "object",
                        "data": [
                        ],
                        "in": "body"
                    }
                ],
                "outputs": []
            }
        },
        {
            "title": "tianqi",
            "name": "DyI3fAv7",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/rtftw/raw/customer/demo_001/test",
                "apiName": "test",
                "inputs": [
                    {
                        "type": "string",
                        "name": "Content-Type",
                        "title": "数据格式",
                        "desc": "application/json",
                        "data": "application/json",
                        "in": "header",
                        "required": true,
                        "mock": "application/json"
                    },
                    {
                        "type": "direct_expr",
                        "name": "city",
                        "data": "start.body_tianqi.city ",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "direct_expr",
                        "name": "key",
                        "data": "start.body_tianqi.key ",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "name": "Access-Token",
                        "data": "",
                        "in": "header"
                    },
                    {
                        "type": "object",
                        "data": [
                        ],
                        "in": "body"
                    }
                ],
                "outputs": []
            }
        }
    ]
}
`,
	`
	{
    "nodes": [
        {
            "title": "",
            "name": "start",
            "type": "input",
            "nextNodes": [
                "DYqIdUzn"
            ],
            "detail": {
                "inputs": [
                    {
                        "type": "string",
                        "name": "serviceName",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": true,
                        "id": "fsGFZu9V-hfloC3POloPn",
                        "descPath": "开始节点.serviceName"
                    },
                    {
                        "type": "string",
                        "name": "objectApiName",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": true,
                        "id": "EmMgJpX-dtqQSZhgDM24L",
                        "descPath": "开始节点.objectApiName"
                    },
                    {
                        "type": "string",
                        "name": "userName",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": true,
                        "id": "ffWkiiZ_AUIS1zNAnNT2H",
                        "descPath": "开始节点.userName"
                    },
                    {
                        "type": "string",
                        "name": "password",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": true,
                        "id": "eeWV3yFEvo9R9UJg6Kzy7",
                        "descPath": "开始节点.password"
                    }
                ],
                "consts": []
            },
            "handles": {
                "right": "start__right"
            }
        },
        {
            "title": "登录",
            "name": "DYqIdUzn",
            "type": "request",
            "nextNodes": [
                "pryLZHnL"
            ],
            "detail": {
                "rawPath": "/system/app/xkh8v/raw/customer/cloudcc/distributor.r",
                "inputs": [
                    {
                        "type": "string",
                        "name": "serviceName",
                        "data": null,
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "userName",
                        "data": "\"clogin\"",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "password",
                        "data": "$start.userName ",
                        "in": "query"
                    }
                ],
                "outputs": [
                    {
                        "type": "string",
                        "name": "returnCode",
                        "data": null
                    },
                    {
                        "type": "object",
                        "name": "userInfo",
                        "data": [
                            {
                                "type": "string",
                                "name": "dbType",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "phone",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "appMainPage",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "timeZone",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "userName",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "isfirstlogin",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "department",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "email",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "profileId",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "loginName",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "roleId",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "mobilePhone",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "isBoundMfa",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "showColleague",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "userId",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "profileName",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "orgName",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "managerId",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "orgId",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "isdisablum",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "roleName",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "isEnableAblumWatermark",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "isUsingGotpFlag",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "pageEffectiveTime",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "countryCode",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "lang",
                                "data": null
                            },
                            {
                                "type": "string",
                                "name": "version",
                                "data": null
                            }
                        ]
                    },
                    {
                        "type": "string",
                        "name": "result",
                        "data": null
                    },
                    {
                        "type": "string",
                        "name": "returnInfo",
                        "data": null
                    },
                    {
                        "type": "string",
                        "name": "binding",
                        "data": null
                    }
                ]
            },
            "handles": {
                "top": "DYqIdUzn__top",
                "left": "DYqIdUzn__left",
                "right": "DYqIdUzn__right",
                "bottom": "DYqIdUzn__bottom"
            }
        },
        {
            "title": "",
            "name": "end",
            "type": "output",
            "nextNodes": [],
            "detail": {
                "body": {
                    "type": "object",
                    "data": [
                        {
                            "type": "direct_expr",
                            "name": "result",
                            "desc": "",
                            "data": "$pryLZHnL.result",
                            "in": "body",
                            "required": false,
                            "id": "n2NiVDFOFD53AhVYwGFsH"
                        },
                        {
                            "type": "direct_expr",
                            "name": "data",
                            "desc": "",
                            "data": "$pryLZHnL.data",
                            "in": "body",
                            "required": false,
                            "id": "qtQpBHvN49MmFkee6o3rN"
                        }
                    ]
                }
            },
            "handles": {
                "left": "end__left"
            }
        },
        {
            "title": "查询",
            "name": "pryLZHnL",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/xkh8v/raw/customer/cloudcc/search.r",
                "inputs": [
                    {
                        "type": "string",
                        "name": "serviceName",
                        "data": null,
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "objectApiName",
                        "data": "$start.serviceName ",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "expressions",
                        "data": "$start.objectApiName ",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "binding",
                        "data": "//##\"ownerid='\"+$DYqIdUzn.userInfo.userId+\"'\"",
                        "in": "query"
                    }
                ],
                "outputs": [
                    {
                        "type": "string",
                        "name": "returnInfo",
                        "data": null,
                        "descPath": "查询.returnInfo"
                    },
                    {
                        "type": "string",
                        "name": "returnCode",
                        "data": null,
                        "descPath": "查询.returnCode"
                    },
                    {
                        "type": "array",
                        "name": "data",
                        "data": [
                            {
                                "type": "object",
                                "name": "",
                                "data": [],
                                "descPath": "查询.data.[0]"
                            }
                        ],
                        "descPath": "查询.data"
                    },
                    {
                        "type": "string",
                        "name": "result",
                        "data": null,
                        "descPath": "查询.result"
                    }
                ]
            },
            "handles": {
                "top": "pryLZHnL__top",
                "left": "pryLZHnL__left",
                "right": "pryLZHnL__right",
                "bottom": "pryLZHnL__bottom"
            }
        }
    ],
    "uis": {
        "edges": [
            {
                "id": "estart-DYqIdUzn-fH2b7I3b",
                "source": "start",
                "target": "DYqIdUzn",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "start__right",
                "targetHandle": "DYqIdUzn__left"
            },
            {
                "id": "eDYqIdUzn-pryLZHnL-NFCknjvE",
                "source": "DYqIdUzn",
                "target": "pryLZHnL",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "DYqIdUzn__right",
                "targetHandle": "pryLZHnL__left"
            },
            {
                "id": "epryLZHnL-end-bfcaw0cG",
                "source": "pryLZHnL",
                "target": "end",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "pryLZHnL__right",
                "targetHandle": "end__left"
            }
        ],
        "metas": [
            {
                "id": "start",
                "type": "input",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "DYqIdUzn",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "end",
                "type": "output",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "pryLZHnL",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            }
        ]
    }
}`,
	`
{
    "nodes": [
        {
            "title": "",
            "name": "start",
            "type": "input",
            "nextNodes": [
                "oKhoaMPd"
            ],
            "detail": {
                "inputs": [
                    {
                        "type": "string",
                        "name": "sk",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": true,
                        "id": "6zesiI-ZFsq6YkIKSs6Di",
                        "descPath": "开始节点.sk"
                    }
                ],
                "consts": []
            },
            "handles": {
                "right": "start__right"
            }
        },
        {
            "title": "获取产品目录",
            "name": "oKhoaMPd",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/nnhqn/raw/customer/newbilling/list_catalogs.r",
                "inputs": [
                    {
                        "type": "direct_expr",
                        "name": "cookie",
                        "required": true,
                        "data": "'sk='+$start.sk",
                        "in": "header"
                    }
                ],
                "outputs": [
                    {
                        "type": "array",
                        "name": "product_set",
                        "data": [
                            {
                                "type": "object",
                                "name": "",
                                "data": [],
                                "descPath": "获取产品目录.product_set.[0]"
                            }
                        ],
                        "descPath": "获取产品目录.product_set"
                    },
                    {
                        "type": "number",
                        "name": "total",
                        "data": null,
                        "descPath": "获取产品目录.total"
                    }
                ]
            },
            "handles": {
                "top": "oKhoaMPd__top",
                "left": "oKhoaMPd__left",
                "right": "oKhoaMPd__right",
                "bottom": "oKhoaMPd__bottom"
            }
        },
        {
            "title": "",
            "name": "end",
            "type": "output",
            "nextNodes": [],
            "detail": {
                "body": {
                    "type": "object",
                    "data": [
                        {
                            "type": "direct_expr",
                            "name": "total",
                            "desc": "",
                            "data": "$oKhoaMPd.total",
                            "in": "body",
                            "required": false,
                            "id": "7LCTePQXmZu0BrwztDGLe"
                        },
                        {
                            "type": "direct_expr",
                            "name": "product_set",
                            "desc": "",
                            "data": "$oKhoaMPd.product_set",
                            "in": "body",
                            "required": false,
                            "id": "-4qAPv-xjq_UAjCj44X3S"
                        }
                    ]
                }
            },
            "handles": {
                "left": "end__left"
            }
        }
    ],
    "uis": {
        "edges": [
            {
                "id": "estart-oKhoaMPd-E9CB69U3",
                "source": "start",
                "target": "oKhoaMPd",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "start__right",
                "targetHandle": "oKhoaMPd__left"
            },
            {
                "id": "eoKhoaMPd-end-jCvR5KGb",
                "source": "oKhoaMPd",
                "target": "end",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "oKhoaMPd__right",
                "targetHandle": "end__left"
            }
        ],
        "metas": [
            {
                "id": "start",
                "type": "input",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "oKhoaMPd",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "end",
                "type": "output",
                "position": {
                    "x": 0,
                    "y": 0
                }
            }
        ]
    }
}
`,
	`
{
    "nodes": [
        {
            "title": "",
            "name": "start",
            "type": "input",
            "nextNodes": [
                "KxJ4NYhA",
                "FooT5996",
                "VH54FFwZ"
            ],
            "detail": {
                "inputs": [
                    {
                        "type": "object",
                        "name": "districts",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "keywords",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "id": "ecYFaMsCQSe8xrAcV4ACe",
                                "descPath": "开始节点.districts.keywords"
                            },
                            {
                                "type": "number",
                                "name": "subdistrict",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "id": "cJ8jpPL3obBo8TB9CJTCt",
                                "descPath": "开始节点.districts.subdistrict"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "id": "EjqAANg3L9KJtUBRmCjgZ",
                        "descPath": "开始节点.districts"
                    },
                    {
                        "type": "object",
                        "name": "weather",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "city",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "id": "NPJbjUnHZSOz3Fuf6gXMj",
                                "descPath": "开始节点.weather.city"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "id": "hHxwSPq4PPsBYD_q7N-Ba",
                        "descPath": "开始节点.weather"
                    },
                    {
                        "type": "object",
                        "name": "weather_new",
                        "desc": "",
                        "data": [
                            {
                                "type": "string",
                                "name": "city",
                                "desc": "",
                                "data": "",
                                "in": "body",
                                "required": false,
                                "id": "09Sjb_xc4cU5RbEkdTjxG",
                                "descPath": "开始节点.weather_new.city"
                            }
                        ],
                        "in": "body",
                        "required": false,
                        "id": "h48PUwL4pUCuIyVdXdAxp",
                        "descPath": "开始节点.weather_new"
                    },
                    {
                        "type": "string",
                        "name": "key",
                        "desc": "",
                        "data": "",
                        "in": "body",
                        "required": false,
                        "id": "1om_If2luyNJNbAHt56yH",
                        "descPath": "开始节点.key"
                    }
                ],
                "consts": []
            },
            "handles": {
                "right": "start__right"
            }
        },
        {
            "title": "行政",
            "name": "KxJ4NYhA",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/dlrg4/raw/customer/default/test05.r",
                "inputs": [
                    {
                        "type": "direct_expr",
                        "name": "key",
                        "required": true,
                        "data": "$start.key ",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "keywords",
                        "data": "$start.districts.keywords ",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "subdistrict",
                        "data": "$start.districts.subdistrict ",
                        "in": "query"
                    }
                ],
                "outputs": [
                    {
                        "type": "string",
                        "name": "status",
                        "desc": "返回结果状态值  值为0或1，0表示失败；1表示成功",
                        "data": null,
                        "descPath": "行政.返回结果状态值值为0或1，0表示失败；1表示成功"
                    },
                    {
                        "type": "string",
                        "name": "info",
                        "desc": "返回状态说明  返回状态说明，status为0时，info返回错",
                        "data": null,
                        "descPath": "行政.返回状态说明返回状态说明，status为0时，info返回错"
                    },
                    {
                        "type": "string",
                        "name": "infocode",
                        "desc": "状态码  返回状态说明，10000代表正确，详情参阅info状态",
                        "data": null,
                        "descPath": "行政.状态码返回状态说明，10000代表正确，详情参阅info状态"
                    }
                ]
            },
            "handles": {
                "top": "KxJ4NYhA__top",
                "left": "KxJ4NYhA__left",
                "right": "KxJ4NYhA__right",
                "bottom": "KxJ4NYhA__bottom"
            }
        },
        {
            "title": "",
            "name": "end",
            "type": "output",
            "nextNodes": [],
            "detail": {
                "body": {
                    "type": "object",
                    "data": [
                        {
                            "type": "object",
                            "name": "data",
                            "desc": "",
                            "data": [
                                {
                                    "type": "object",
                                    "name": "weather",
                                    "desc": "",
                                    "data": [
                                        {
                                            "type": "direct_expr",
                                            "name": "ststus",
                                            "desc": "",
                                            "data": "$FooT5996.status",
                                            "in": "body",
                                            "required": false,
                                            "id": "YzdESzU5BWYiSfCDipoIh"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "count",
                                            "desc": "",
                                            "data": "$FooT5996.count",
                                            "in": "body",
                                            "required": false,
                                            "id": "J7mBITgNkKggGaawgYOGB"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "info",
                                            "desc": "",
                                            "data": "$FooT5996.info",
                                            "in": "body",
                                            "required": false,
                                            "id": "5IP_bF_llKKBPjQatIzNs"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "infocode",
                                            "desc": "",
                                            "data": "$FooT5996.infocode",
                                            "in": "body",
                                            "required": false,
                                            "id": "3hOKP-u5YV2jdl5lRwsBg"
                                        },
                                        {
                                            "type": "array",
                                            "name": "lives",
                                            "desc": "",
                                            "data": [
                                                {
                                                    "type": "object",
                                                    "name": "",
                                                    "desc": "",
                                                    "data": [
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "wind",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].winddirection",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "dspxnpJVBtfPck0dbX1cI"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "windp",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].windpower",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "wl1Zg1BHk1VxPaJRlFR7d"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "repo",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].reporttime",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "kSwLai0kB6NlRGpX4wNn_"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "city",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].city",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "vesUjfe8FxchhG3jYKXJ4"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "adcode",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].adcode",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "e22P7O7dUYX-jPd946671"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "temr",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].temperature",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "e1Nx2d4tsKfrcxt-OSS3q"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "pro",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].province",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "zYL5dNZa4KK1S3FwFdvTQ"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "weat",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].weather",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "GLghrCnebIh4Xhth2kDOu"
                                                        },
                                                        {
                                                            "type": "direct_expr",
                                                            "name": "hum",
                                                            "desc": "",
                                                            "data": "$FooT5996.lives.[0].humidity",
                                                            "in": "body",
                                                            "required": false,
                                                            "id": "m2Zio4j235cKfr6K-jnBQ"
                                                        }
                                                    ],
                                                    "in": "body",
                                                    "required": false,
                                                    "id": "uHQbNv9GQT7rKCmsizVrQ"
                                                }
                                            ],
                                            "in": "body",
                                            "required": false,
                                            "id": "JsTGo6CIxerVNy5tK4RnQ"
                                        }
                                    ],
                                    "in": "body",
                                    "required": false,
                                    "id": "JWjrxsHMxohivTPL0EP0k"
                                },
                                {
                                    "type": "array",
                                    "name": "districts",
                                    "desc": "",
                                    "data": [
                                        {
                                            "type": "direct_expr",
                                            "name": "",
                                            "desc": "",
                                            "data": "$KxJ4NYhA.status",
                                            "in": "body",
                                            "required": false,
                                            "id": "eHdliP7_uj7Cd_N66uWM-"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "",
                                            "desc": "",
                                            "data": "$KxJ4NYhA.info",
                                            "in": "body",
                                            "required": false,
                                            "id": "tUHNLVwFDNl0baZK_fR-Y"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "",
                                            "desc": "",
                                            "data": "$KxJ4NYhA.infocode",
                                            "in": "body",
                                            "required": false,
                                            "id": "aS3mDHCxxPERfsdKZuRcP"
                                        }
                                    ],
                                    "in": "body",
                                    "required": false,
                                    "id": "RlLvmEg9unprYbFFdYxw4"
                                },
                                {
                                    "type": "object",
                                    "name": "weather_new",
                                    "desc": "",
                                    "data": [
                                        {
                                            "type": "direct_expr",
                                            "name": "city",
                                            "desc": "",
                                            "data": "($VH54FFwZ.update_time)+($VH54FFwZ.city)",
                                            "in": "body",
                                            "required": false,
                                            "id": "xPdS8jagICQnmb2deqIHV"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "time",
                                            "desc": "",
                                            "data": "$VH54FFwZ.date",
                                            "in": "body",
                                            "required": false,
                                            "id": "AD7SzNNu6dsvoZin6eDWi"
                                        },
                                        {
                                            "type": "direct_expr",
                                            "name": "weather",
                                            "desc": "",
                                            "data": "$VH54FFwZ.weather",
                                            "in": "body",
                                            "required": false,
                                            "id": "h_fWwETQGOqGraXJe3-4d"
                                        }
                                    ],
                                    "in": "body",
                                    "required": false,
                                    "id": "7jBW6RRfbQLTgS1E42xUF"
                                }
                            ],
                            "in": "body",
                            "required": false,
                            "id": "pZP0oFQ8Gc6WogiystGA3"
                        }
                    ]
                }
            },
            "handles": {
                "left": "end__left"
            }
        },
        {
            "title": "天气新",
            "name": "FooT5996",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/dlrg4/raw/customer/sys_02/test.r",
                "inputs": [
                    {
                        "type": "direct_expr",
                        "name": "city",
                        "required": true,
                        "data": "$start.weather.city ",
                        "in": "query"
                    },
                    {
                        "type": "direct_expr",
                        "name": "key",
                        "required": true,
                        "data": "$start.key ",
                        "in": "query"
                    }
                ],
                "outputs": [
                    {
                        "type": "array",
                        "name": "lives",
                        "data": [
                            {
                                "type": "object",
                                "name": "",
                                "data": [
                                    {
                                        "type": "string",
                                        "name": "winddirection",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].winddirection"
                                    },
                                    {
                                        "type": "string",
                                        "name": "windpower",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].windpower"
                                    },
                                    {
                                        "type": "string",
                                        "name": "reporttime",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].reporttime"
                                    },
                                    {
                                        "type": "string",
                                        "name": "city",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].city"
                                    },
                                    {
                                        "type": "number",
                                        "name": "adcode",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].adcode"
                                    },
                                    {
                                        "type": "number",
                                        "name": "temperature",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].temperature"
                                    },
                                    {
                                        "type": "string",
                                        "name": "province",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].province"
                                    },
                                    {
                                        "type": "string",
                                        "name": "weather",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].weather"
                                    },
                                    {
                                        "type": "number",
                                        "name": "humidity",
                                        "data": null,
                                        "descPath": "天气新.lives.[0].humidity"
                                    }
                                ],
                                "descPath": "天气新.lives.[0]"
                            }
                        ],
                        "descPath": "天气新.lives"
                    },
                    {
                        "type": "string",
                        "name": "info",
                        "data": null,
                        "descPath": "天气新.info"
                    },
                    {
                        "type": "number",
                        "name": "status",
                        "data": null,
                        "descPath": "天气新.status"
                    },
                    {
                        "type": "number",
                        "name": "count",
                        "data": null,
                        "descPath": "天气新.count"
                    },
                    {
                        "type": "number",
                        "name": "infocode",
                        "data": null,
                        "descPath": "天气新.infocode"
                    }
                ]
            },
            "handles": {
                "top": "FooT5996__top",
                "left": "FooT5996__left",
                "right": "FooT5996__right",
                "bottom": "FooT5996__bottom"
            }
        },
        {
            "title": "查询今日天气",
            "name": "VH54FFwZ",
            "type": "request",
            "nextNodes": [
                "end"
            ],
            "detail": {
                "rawPath": "/system/app/dlrg4/raw/customer/sys_03/test04.r",
                "inputs": [
                    {
                        "type": "direct_expr",
                        "name": "city",
                        "data": "$start.weather_new.city ",
                        "in": "query"
                    }
                ],
                "outputs": [
                    {
                        "type": "string",
                        "name": "city",
                        "desc": "城市",
                        "data": null,
                        "descPath": "查询今日天气.城市"
                    },
                    {
                        "type": "string",
                        "name": "update_time",
                        "desc": "时间",
                        "data": null,
                        "descPath": "查询今日天气.时间"
                    },
                    {
                        "type": "string",
                        "name": "date",
                        "desc": "日期",
                        "data": null,
                        "descPath": "查询今日天气.日期"
                    },
                    {
                        "type": "object",
                        "name": "weather",
                        "desc": "详情",
                        "data": [],
                        "descPath": "查询今日天气.详情"
                    }
                ]
            },
            "handles": {
                "top": "VH54FFwZ__top",
                "left": "VH54FFwZ__left",
                "right": "VH54FFwZ__right",
                "bottom": "VH54FFwZ__bottom"
            }
        }
    ],
    "uis": {
        "edges": [
            {
                "id": "estart-KxJ4NYhA-iaF2FfEC",
                "source": "start",
                "target": "KxJ4NYhA",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "start__right",
                "targetHandle": "KxJ4NYhA__left"
            },
            {
                "id": "eKxJ4NYhA-end-LF1wugoF",
                "source": "KxJ4NYhA",
                "target": "end",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "KxJ4NYhA__right",
                "targetHandle": "end__left"
            },
            {
                "id": "estart-FooT5996-xr8DoVLD",
                "source": "start",
                "target": "FooT5996",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "start__right",
                "targetHandle": "FooT5996__left"
            },
            {
                "id": "eFooT5996-end-Sh2oyDQW",
                "source": "FooT5996",
                "target": "end",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "FooT5996__right",
                "targetHandle": "end__left"
            },
            {
                "id": "estart-VH54FFwZ-rFSWY0Zf",
                "source": "start",
                "target": "VH54FFwZ",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "start__right",
                "targetHandle": "VH54FFwZ__left"
            },
            {
                "id": "eVH54FFwZ-end-kMNz82bj",
                "source": "VH54FFwZ",
                "target": "end",
                "type": "smart",
                "style": {
                    "stroke": "#CBD5E1"
                },
                "arrowHeadType": "arrowclosed",
                "sourceHandle": "VH54FFwZ__right",
                "targetHandle": "end__left"
            }
        ],
        "metas": [
            {
                "id": "start",
                "type": "input",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "KxJ4NYhA",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "end",
                "type": "output",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "FooT5996",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            },
            {
                "id": "VH54FFwZ",
                "type": "request",
                "position": {
                    "x": 0,
                    "y": 0
                }
            }
        ]
    }
}
`,
}
