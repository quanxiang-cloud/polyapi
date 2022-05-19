package arrange

import (
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
)

func init() {
	enumset.FinishReg()
}

// flexFactory regist some flex json objects
var flexFactory = factory.NewFlexObjFactory("flexObjCreator")

// GetFactory return the factory object
func GetFactory() *factory.FlexObjFactory { return flexFactory }

func init() {
	flexFactory.MustReg(InputNodeDetail{
		Inputs: []ValueDefine{
			ValueDefine{
				InputValue: InputValue{
					Type:     "string",
					Name:     "input1",
					In:       "header",
					Desc:     "input1 from header",
					Required: true,
				},
				Default: "foo",
				Enums:   []string{"foo", "bar"},
			},
			ValueDefine{
				InputValue: InputValue{
					Type:     "number",
					Name:     "input2",
					In:       "body",
					Desc:     "input2 from body",
					Required: true,
				},
				Default: "foo",
			},
		},
		Consts: ValueSet{
			InputValue{
				Type: "number",
				Name: "input3",
				In:   "body",
				Data: FlexJSONObject{
					D: ValNumber("1"),
				},
				Desc: "const input3 from body",
			},
		},
	})
	flexFactory.MustReg(IfNodeDetail{
		Cond: CondExpr{
			Op: "",
			Value: Value{
				Desc: "req1.total_count /gt 0",
				Type: "exprcmp",
				Data: FlexJSONObject{
					D: ValExprCmp{
						LValue: Value{
							Type:  "number",
							Field: "req1.total_count",
						},
						Cmp: "gt",
						RValue: Value{
							Type: "number",
							Data: FlexJSONObject{
								D: "0",
							},
						},
					},
				},
			},
		},
		Yes: "req2",
		No:  "req3",
	})
	flexFactory.MustReg(RequestNodeDetail{
		RawPath: "/system/api/queryUser",
		Inputs: ValueSet{
			InputValue{
				Type: "string",
				Name: "input1",
				In:   "body",
				Data: FlexJSONObject{
					D: ValString("foo"),
				},
			},
			InputValue{
				Type: "string",
				Name: "input2",
				In:   "header",
				Data: FlexJSONObject{
					D: ValString("bar"),
				},
			},
			InputValue{
				Type:  "string",
				Name:  "input3",
				In:    "body",
				Field: "req1.result",
			},
		},
	})
	flexFactory.MustReg(OutputNodeDetail{
		Header: ValueSet{},
		Body: Value{
			Type: "object",
			Data: FlexJSONObject{
				D: ValObject{
					Value{
						Type:  "field",
						Field: "req1",
					},
					Value{
						Type:  "field",
						Field: "req2",
					},
					Value{
						Type:  "field",
						Name:  "req3x",
						Field: "req3",
					},
				},
			},
		},
	})

	flexFactory.MustReg(Arrange{
		Info: nil,
		Nodes: []Node{
			Node{
				Name:      "start",
				Title:     "开始节点",
				Type:      "input",
				Desc:      "the only entrance of the arrange, must name as start",
				NextNodes: []string{"cond1"},
				Detail: FlexJSONObject{
					D: InputNodeDetail{
						Inputs: []ValueDefine{
							ValueDefine{
								InputValue: InputValue{
									Type:     "string",
									Name:     "input1",
									In:       "header",
									Desc:     "input1 from header",
									Required: true,
								},
								Default: "foo",
								Enums:   []string{"foo", "bar"},
							},
							ValueDefine{
								InputValue: InputValue{
									Type:     "number",
									Name:     "input2",
									In:       "body",
									Desc:     "input2 from body",
									Required: true,
								},
								Default: "1",
								//Ranges:   []string{"1", "10"},
							},
						},
						Consts: ValueSet{
							InputValue{
								Type: "number",
								Name: "input3",
								In:   "body",
								Data: FlexJSONObject{
									D: ValNumber("1"),
								},
								Desc: "const input3 from body",
							},
						},
					},
				},
			},
			Node{
				Name:      "cond1",
				Title:     "条件节点1",
				Type:      "if",
				Desc:      "check some condition",
				NextNodes: []string{},
				Detail: FlexJSONObject{
					D: IfNodeDetail{
						Cond: CondExpr{
							Op: "",
							Value: Value{
								Desc: "start.input2 /gt 0",
								Type: "exprcmp",
								Data: FlexJSONObject{
									D: ValExprCmp{
										LValue: Value{
											Type: "direct_expr",
											Data: FlexJSONObject{
												D: ValDirectExpr("start.input2*3 + 1"),
											},
										},
										Cmp: "gt",
										RValue: Value{
											Type: "number",
											Data: FlexJSONObject{
												D: "0",
											},
										},
									},
								},
							},
						},
						Yes: "req1",
						No:  "",
					},
				},
			},
			Node{
				Name:      "req1",
				Title:     "请求节点1",
				Type:      "request",
				Desc:      "do some http request",
				NextNodes: []string{"end"},
				Detail: FlexJSONObject{
					D: RequestNodeDetail{
						RawPath: "/system/raw/testraw",
						Inputs: ValueSet{
							InputValue{
								Type: "string",
								Name: "input1",
								In:   "body",
								Data: FlexJSONObject{
									D: ValString("foo"),
								},
							},
							InputValue{
								Type:  "string",
								Name:  "input2",
								In:    "header",
								Field: "start.input1",
							},
						},
					},
				},
			},
			Node{
				Name:      "end",
				Title:     "结束节点",
				Type:      "output",
				Desc:      "the only end of the arrange, must name as end",
				NextNodes: []string{},
				Detail: FlexJSONObject{
					D: OutputNodeDetail{
						Body: Value{

							Name: "merged",
							Type: "mergeobj",
							Data: FlexJSONObject{
								D: ValMergeObj{
									Value{
										Type: "filter",
										Data: FlexJSONObject{
											D: ValFiltObj{
												Source: "req1",
												White: FieldMap{
													"foo": "fooRename",
													"bar": "",
												},
												Black: FieldMap{},
											},
										},
									},
									Value{
										Type: "filter",
										Data: FlexJSONObject{
											D: ValFiltObj{
												Source: "start",
												White:  FieldMap{},
												Black: FieldMap{
													"input3": "",
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
	})

	//factory.MustReg((*DispDemo)(nil))
}

// node type enum
var (
	// NodeTypeEnum represents node enum set
	NodeTypeEnum    = enumset.New(nil)
	NodeTypeInput   = NodeTypeEnum.MustReg("input")
	NodeTypeIf      = NodeTypeEnum.MustReg("if")
	NodeTypeRequest = NodeTypeEnum.MustReg("request")
	NodeTypeOutput  = NodeTypeEnum.MustReg("output")
	NodeTypeArrange = NodeTypeEnum.MustReg("arrange")
)

var (
	// NodeNameEnum represents node name enum set
	NodeNameEnum = enumset.New(nil)
	// NodeNameStart represents node name start
	NodeNameStart = NodeNameEnum.MustReg("start")
	// NodeNameEnd represents node name end
	NodeNameEnd = NodeNameEnum.MustReg("end")
)
