package exprx

import (
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/factory"
)

// flexFactory regist some flex json objects
var flexFactory = expr.GetFactory().Clean()

// GetFactory return the factory object
func GetFactory() *factory.FlexObjFactory { return flexFactory }

func init() {
	flexFactory.MustReg(ValNumber("123"))
	flexFactory.MustReg(ValString("foo"))
	flexFactory.MustReg(ValBoolean("true"))
	flexFactory.MustReg(ValObject{
		Value{
			Name: "a",
			Type: "string",
			Data: FlexJSONObject{
				D: NewStringer("foo", ValTypeString),
			},
		},
		Value{
			Name: "b",
			Type: "number",
			Data: FlexJSONObject{
				D: NewStringer("123", ValTypeNumber),
			},
		},
		Value{
			Name: "c",
			Type: "boolean",
			Data: FlexJSONObject{
				D: NewStringer("true", ValTypeBoolean),
			},
		},
	})
	flexFactory.MustReg(ValArray{
		Value{
			Type: "string",
			Data: FlexJSONObject{
				D: ValString("foo"),
			},
		},
		Value{
			Type: "string",
			Data: FlexJSONObject{
				D: ValString("bar"),
			},
		},
	})
	flexFactory.MustReg(ValMergeObj{
		Value{
			Type:  "field",
			Field: "req1.field1",
		},
		Value{
			Type:  "field",
			Field: "req2.field1",
		},
	})
	flexFactory.MustReg(ValFiltObj{
		Source: "req1.arrayField1",
		White: FieldMap{
			"a": "a1",
			"b": "",
		},
		Black: FieldMap{
			"c": "",
		},
		Filter: CondExpr{
			Op: "",
			Value: Value{
				Desc: "field1 /gt 0",
				Type: "exprcmp",
				Data: FlexJSONObject{
					D: ValExprCmp{
						LValue: Value{
							Type:  "number",
							Field: "field1",
						},
						Cmp: "gt",
						RValue: Value{
							Type: "number",
							Data: FlexJSONObject{
								D: ValNumber("0"),
							},
						},
					},
				},
			},
		},
	})
	flexFactory.MustReg(ValTimestamp("2020-12-31T12:34:56Z"))
	flexFactory.MustReg(ValAction("createUser"))
	flexFactory.MustReg(ValArrayString{"foo", "bar"})
	flexFactory.MustReg(
		ValArrayStringElem{
			Name:  "arrayName",
			Array: ValArrayString{"foo", "bar"},
		})

	flexFactory.MustReg(ValExpr{
		Op: "",
		Value: Value{
			Type: "exprgroup",
			Data: FlexJSONObject{
				D: ValExprGroup{
					ValExpr{
						Value: Value{
							Type:  "field",
							Field: "req.field1",
						},
					},
					ValExpr{
						Op: "mul",
						Value: Value{
							Type: "number",
							Data: FlexJSONObject{
								D: ValNumber("2"),
							},
						},
					},
					ValExpr{
						Op: "add",
						Value: Value{
							Type: "number",
							Data: FlexJSONObject{
								D: ValNumber("1"),
							},
						},
					},
				},
			},
		},
	})
	flexFactory.MustReg(ValExprCmp{
		LValue: Value{
			Type:  "field",
			Field: "req1.a",
		},
		Cmp: "le",
		RValue: Value{
			Type: "number",
			Data: FlexJSONObject{
				D: ValNumber("0"),
			},
		},
	})
	flexFactory.MustReg(ValExprSel{
		Cond: CondExpr{
			Op: "",
			Value: Value{
				Desc: "field1 /gt 0",
				Type: "exprcmp",
				Data: FlexJSONObject{
					D: ValExprCmp{
						LValue: Value{
							Type:  "number",
							Field: "req1.field1",
						},
						Cmp: "gt",
						RValue: Value{
							Type: "number",
							Data: FlexJSONObject{
								D: ValNumber("0"),
							},
						},
					},
				},
			},
		},
		Yes: Value{
			Type:  "field",
			Field: "req1.field1"},
		No: Value{
			Type: "number",
			Data: FlexJSONObject{
				D: ValNumber("0"),
			},
		},
	})
	flexFactory.MustReg(ValExprFunc{
		Func: "format",
		Paras: []Value{
			Value{
				Type: "string",
				Data: FlexJSONObject{
					D: ValString("%d days left"),
				},
			},
			Value{
				Type: "number",
				Data: FlexJSONObject{
					D: ValNumber("2"),
				},
			},
		},
	})
	flexFactory.MustReg(ValExprGroup{
		ValExpr{
			Value: Value{
				Type:  "field",
				Field: "req.field1",
			},
		},
		ValExpr{
			Op: "mul",
			Value: Value{
				Type: "number",
				Data: FlexJSONObject{
					D: ValNumber("2"),
				},
			},
		},
		ValExpr{
			Op: "add",
			Value: Value{
				Type: "number",
				Data: FlexJSONObject{
					D: ValNumber("1"),
				},
			},
		},
	})

	flexFactory.MustReg(ValDirectExpr("(req1.x+1)*2"))
}
