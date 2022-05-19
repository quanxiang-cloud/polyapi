package exprx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
)

func init() {
	testMode = true
}

// verify interfaces

var _ = []DelayedJSONDecoder{
	(*XInputValue)(nil),
	(*Value)(nil),
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	//(*ValArrayElem)(nil),
	(*ValArrayString)(nil),
	(*ValArrayStringElem)(nil),
	(*ValMergeObj)(nil),
	(*ValFiltObj)(nil),
	// (*ValPath)(nil),
	// (*ValHeader)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),

	(*ValExpr)(nil),
	(*ValExprCmp)(nil),
	(*ValExprSel)(nil),
	(*ValExprFunc)(nil),
	(*ValExprGroup)(nil),
	(*ValDirectExpr)(nil),

	(*ValueSet)(nil),

	(*CondExpr)(nil),
	(*CondExprGroup)(nil),
}

var _ = []ScriptElem{
	(*XInputValue)(nil),
	(*Value)(nil),
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	//(*ValArrayElem)(nil),
	(*ValArrayStringElem)(nil),
	(*ValArrayString)(nil),
	(*ValArrayString)(nil),
	(*ValMergeObj)(nil),
	(*ValFiltObj)(nil),
	//(*ValPath)(nil),
	//(*ValHeader)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),

	(*ValExpr)(nil),
	(*ValExprCmp)(nil),
	(*ValExprSel)(nil),
	(*ValExprFunc)(nil),
	(*ValExprGroup)(nil),
	(*ValDirectExpr)(nil),

	(*ValueSet)(nil),

	(*CondExpr)(nil),
	(*CondExprGroup)(nil),
}

var _ = []NamedScriptElem{
	(*XInputValue)(nil),
	(*Value)(nil),
	//	(*ValArrayElem)(nil),
	(*ValArrayStringElem)(nil),
	(*FieldRef)(nil),
	(*ValExpr)(nil),
	(*ValFiltObj)(nil),
}

var _ = []NamedType{
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	//	(*ValArrayElem)(nil),
	(*ValArrayStringElem)(nil),
	(*ValArrayString)(nil),
	(*ValMergeObj)(nil),
	(*ValFiltObj)(nil),
	// (*ValPath)(nil),
	// (*ValHeader)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),

	(*ValExpr)(nil),
	(*ValExprCmp)(nil),
	(*ValExprSel)(nil),
	(*ValExprFunc)(nil),
	(*ValExprGroup)(nil),
	(*ValDirectExpr)(nil),
}

var _ = []GenSampler{
	(*ValTimestamp)(nil),
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	(*ValArrayStringElem)(nil),
	(*ValArrayString)(nil),
}

var _ = []Stringer{
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	// (*ValPath)(nil),
	// (*ValHeader)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),
}

var _ = []json.Marshaler{
	(*FlexJSONObject)(nil),
	(*ValArrayString)(nil),
}

var _ = []json.Unmarshaler{
	(*FlexJSONObject)(nil),
	(*ValArrayString)(nil),
}

var _ = []StringerWithError{
	(*ValDirectExpr)(nil),
}

func TestValuesScript(t *testing.T) {
	type testCase struct {
		id   int
		obj  interface{}
		expr string

		// check
		empty     bool
		expect    string
		String    string
		name      string
		nameTitle string
	}
	f := GetFactory()
	newFieldRef := func() *FieldRef {
		var field FieldRef
		return &field
	}

	cases := []*testCase{
		&testCase{
			id:     1,
			obj:    f.MustCreate(ExprTypeDirectExpr.String()),
			expr:   `a+b*c`,
			expect: `a+b*c`,
		},
		&testCase{
			obj:    f.MustCreate(ExprTypeDirectExpr.String()),
			expr:   `a+b*2.5`,
			expect: `a+b*2.5`,
		},
		&testCase{
			obj:    f.MustCreate(ExprTypeDirectExpr.String()),
			expr:   `a+b*c.x`,
			expect: `a+b*d.c.x`,
		},
		&testCase{
			obj:    f.MustCreate(ExprTypeDirectExpr.String()),
			expr:   `a+(b*c.x)`,
			expect: `a+(b*d.c.x)`,
		},
		&testCase{
			obj:    newFieldRef(),
			expr:   `req1.x`,
			expect: `d.req1.x`,
			name:   `x`,
		},
		&testCase{
			obj:    newFieldRef(),
			expr:   ``,
			expect: XValTypeUndefined.String(),
			name:   ``,
		},
		&testCase{
			obj:    newFieldRef(),
			expr:   `$`, //$
			expect: `d`, //d
			name:   `d`, //d
		},
		&testCase{
			obj:    f.MustCreate(ValTypeString.String()),
			expr:   `foo`,
			expect: `"foo"`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeBoolean.String()),
			expr:   `true`,
			expect: `true`,
		},
		&testCase{
			id:     10,
			obj:    f.MustCreate(ValTypeNumber.String()),
			expr:   `12.34`,
			expect: `12.34`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeTimestamp.String()),
			expr:   `YYYY-MM-DDThh:mm:ssZ`,
			expect: `timestamp("YYYY-MM-DDThh:mm:ssZ")`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeTimestamp.String()),
			expr:   `YYYY-MM-DDThh:mm:ss+0000`,
			expect: `timestamp("YYYY-MM-DDThh:mm:ss+0000")`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeTimestamp.String()),
			expr:   ``,
			expect: `timestamp("")`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeAction.String()),
			expr:   `foo`,
			expect: `"foo"`,
		},
		&testCase{
			obj: &ValExprFunc{
				Func: "fn",
				Paras: []Value{
					Value{
						Type: "string",
						Data: FlexJSONObject{
							D: json.RawMessage(`"foo"`),
						},
					},
					Value{
						Type: "number",
						Data: FlexJSONObject{
							D: json.RawMessage(`"12.34"`),
						},
					},
				},
			},
			expect: `fn("foo", 12.34)`,
		},
		&testCase{
			obj: &ValExprCmp{
				Cmp: CmpGE.String(),
				LValue: Value{
					Type:  "field",
					Field: "req1.x",
				},
				RValue: Value{
					Type: "number",
					Data: FlexJSONObject{
						D: json.RawMessage(`"12.34"`),
					},
				},
			},
			expect: `d.req1.x>=12.34`,
		},
		&testCase{
			obj: &ValExprSel{
				Cond: ValExpr{
					Value: Value{
						Type: ExprTypeDirectExpr,
						Data: FlexJSONObject{
							D: json.RawMessage(`"req1.flag==1"`),
						},
					},
				},
				Yes: Value{
					Type:  "field",
					Field: "req1.x",
				},
				No: Value{
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"Bob"`),
					},
				},
			},
			expect: `pdSelect(d.req1.flag==1, d.req1.x, "Bob")`,
		},
		&testCase{
			obj: &ValExprGroup{
				ValExpr{
					Op: "and",
					Value: Value{
						Type: ExprTypeDirectExpr,
						Data: FlexJSONObject{
							D: json.RawMessage(`"req1.flag==1"`),
						},
					},
				},
				ValExpr{
					Op: "and",
					Value: Value{
						Type: ExprTypeDirectExpr,
						Data: FlexJSONObject{
							D: json.RawMessage(`"req2.flag==1"`),
						},
					},
				},
				ValExpr{
					Op: "add",
					Value: Value{
						Type: ExprTypeDirectExpr,
						Data: FlexJSONObject{
							D: json.RawMessage(`"req3.flag"`),
						},
					},
				},
			},
			expect: `(d.req1.flag==1&&d.req2.flag==1+d.req3.flag)`,
		},
		&testCase{
			obj:    f.MustCreate(ValTypeArrayString.String()),
			expr:   `"foo,bar"`,
			expect: `["foo","bar"]`,
		},
		&testCase{
			id: 20,
			obj: &ValArrayStringElem{
				Name:  "ary",
				Array: []string{"foo", "bar"},
			},
			expect: "  \"ary.1\": \"foo\",\n  \"ary.2\": \"bar\",\n",
			name:   `ary`,
		},
		&testCase{
			obj: &ValExpr{
				Op: "&&",
				Value: Value{
					Type: ExprTypeDirectExpr,
					Data: FlexJSONObject{
						D: json.RawMessage(`"req1.flag==1"`),
					},
				},
			},
			expect: `d.req1.flag==1`,
			name:   `_`,
		},
		&testCase{
			obj: &ValObject{
				Value{
					Name: "x",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"foo"`),
					},
				},
				Value{
					//Name:  "y",
					Type:  "boolean",
					Field: "req1.y",
				},
				Value{
					//Name:  "z",
					Type: "array_string_elem",
					Data: FlexJSONObject{
						D: json.RawMessage(`{"name":"z","array":"foo,bar"}`),
					},
				},
			},
			expect: "{\n    \"x\": \"foo\",\n    \"y\": d.req1.y,\n    \"z.1\": \"foo\",\n    \"z.2\": \"bar\",\n  }",
		},
		&testCase{
			obj: &ValArray{
				Value{
					Name: "x",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"foo"`),
					},
				},
				Value{
					Name: "y",
					Type: "boolean",
					Data: FlexJSONObject{
						D: json.RawMessage(`"true"`),
					},
				},
			},
			expect: "[\n    \"foo\",\n    true,\n  ]",
		},
		&testCase{
			obj:    &Value{},
			expect: `undefined`,
			name:   "_",
			empty:  true,
		},
		&testCase{
			obj:    &XInputValue{},
			expect: `undefined`,
			name:   "_",
		},
		&testCase{
			obj: &Value{
				Type:  "string",
				Name:  `n`,
				Title: `title`,
				Desc:  `desc`,
				Data: FlexJSONObject{
					D: json.RawMessage(`"foo"`),
				},
			},
			name:      `n`,
			nameTitle: `title`,
			expect:    `"foo"`,
			String:    `foo`,
		},
		&testCase{
			obj: &Value{
				Type:  "string",
				Name:  `n`,
				Title: ``,
				Desc:  `desc`,
				Data: FlexJSONObject{
					D: json.RawMessage(`"foo"`),
				},
			},
			name:      `n`,
			nameTitle: `n`,
			String:    `foo`,
			expect:    `"foo"`,
		},
		&testCase{
			obj: &XInputValue{
				Type:  "string",
				Field: "req1.x",
				Data: FlexJSONObject{
					D: json.RawMessage(`"foo"`),
				},
				In: "header",
			},
			name:   "_",
			String: "req1.x",
			expect: `d.req1.x`,
		},
		&testCase{
			obj: &XInputValue{
				Type:  "string",
				Name:  `n`,
				Title: `title`,
				Desc:  `desc`,
				Data: FlexJSONObject{
					D: json.RawMessage(`"foo"`),
				},
			},
			name:      `n`,
			nameTitle: `title`,
			expect:    `"foo"`,
			String:    `foo`,
		},
		&testCase{
			id: 30,
			obj: &XInputValue{
				Type:  "string",
				Name:  `n`,
				Title: ``,
				Desc:  `desc`,
				Data: FlexJSONObject{
					D: json.RawMessage(`"foo"`),
				},
			},
			name:      `n`,
			nameTitle: `n`,
			expect:    `"foo"`,
			String:    `foo`,
		},
		&testCase{
			obj: &ValFiltObj{
				Source: "req1.x",
				White:  FieldMap{"x": "xx", "y": "yy"},
			},
			name:   `x`,
			expect: "pdFiltObject(\n    d.req1.x,\n    {white:{x:\"xx\", y:\"yy\", }, black:{}, },\n    function (d) { return undefined }\n  )",
		},
		&testCase{
			obj: &ValMergeObj{
				Value{
					Field: "req1.x",
				},
				Value{
					Field: "req1.y",
				},
			},
			expect: "pdMergeObjs(\n      d.req1.x,\n      d.req1.y\n    )",
		},
		&testCase{
			obj: &ValueSet{
				InputValue{
					Name: "x",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"foo"`),
					},
					In: "body",
				},
				InputValue{
					Type: "array_string_elem",
					Data: FlexJSONObject{
						D: json.RawMessage(`{"name":"z","array":"foo,bar"}`),
					},
					In: "body",
				},
				InputValue{
					Name: "y",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"bar"`),
					},
					In: "header",
				},
			},
			expect: "{\n    \"x\": \"foo\",\n    \"z.1\": \"foo\",\n    \"z.2\": \"bar\",\n  }",
		},
		&testCase{
			obj: &FmtAPIInOut{
				Input: InputNodeDetail{
					Inputs: []ValueDefine{
						ValueDefine{
							InputValue: InputValue{
								Name:     "b",
								Type:     "number",
								In:       "body",
								Required: true,
								Data: FlexJSONObject{
									D: json.RawMessage(`""`),
								},
							},
						},
					},
					Consts: ValueSet{
						InputValue{
							Name: "x",
							Type: "string",
							Data: FlexJSONObject{
								D: json.RawMessage(`"bar"`),
							},
							In: "header",
						},
					},
				},
				Output: OutputNodeDetail{
					Doc: []ValueDefine{
						ValueDefine{
							InputValue: InputValue{
								Name: "x",
								Type: "string",
								Data: FlexJSONObject{
									D: json.RawMessage(`""`),
								},
							},
						},
					},
				},
			},
		},
	}
	verifyValueX := func(i int, p *testCase, t *testing.T) {
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d verifyValueX(%s) fail: expect %v, got %v", i+1, msg, expect, got)
			}
		}
		switch x := p.obj.(type) {
		case *Value:
			assert(x.Empty(), p.empty, "Value.Empty")
			assert(x.GetAsString(), p.String, "Value.GetAsString")
			assert(x.DenyFieldRefer(), nil, "Value.DenyFieldRefer")
		case *InputValue:
			assert(x.GetAsString(), p.String, "InputValue.GetAsString")
			assert(x.DenyFieldRefer(), x.DenyFieldRefer(), "InputValue.DenyFieldRefer")
		case *ValueSet:
			x.AddKV("", "", ValTypeString, ParaTypeHeader)
			x.AddKV("x", "foo", ValTypeString, ParaTypePath)
			x.AddExKV("y", "bar", ValTypeString, ParaTypePath)
			assert(x.GetAction() == nil, true, "ValueSet.GetAction")
			x.AddExKV("y", "bar", ValTypeAction, ParaTypeHide)
			x.AddExKV("t", "zz", ValTypeTimestamp, ParaTypeBody)
			assert(x.GetAction().GetAsString(), "bar", "ValueSet.GetAction")

			args := &RequestArgs{
				Action: "action",
				URL:    "http://sample.com/api/v1/:x/:y/request",
				Body: []byte(fmt.Sprintf(`{"x":"foo","%s":{"x":"bar","y":"yy","security_key":"skey"}}`,
					polysign.XPolyBodyHideArgs)),
				Header: http.Header{},
			}
			err := x.PrepareRequest(args)
			assert(err, nil, "FmtAPIInOut.PrepareRequest")
			assert(args.URL, "http://sample.com/api/v1/bar/yy/request", "FmtAPIInOut.PrepareRequest")

			f, arg, err := x.ResolvePathArgs("/api/v1/:x/:y/request", "node")
			assert(err, nil, "FmtAPIInOut.ResolvePathArgs")
			assert(f, "/api/v1/%v/%v/request", "FmtAPIInOut.ResolvePathArgs")
			assert(arg, []string{`"foo"`, `"bar"`}, "FmtAPIInOut.ResolvePathArgs")

			url, err := x.ReplacePathArgs("/api/v1/:x/:y/request", "node")
			assert(err, nil, "FmtAPIInOut.ReplacePathArgs")
			assert(url, "/api/v1/foo/bar/request", "FmtAPIInOut.ReplacePathArgs")

		case *FmtAPIInOut:
			x.SetAccessURL("/system/raw")
			assert(x.URL, "/api/v1/polyapi/request/system/raw", "FmtAPIInOut")
			x.SetAccessURL("/system/poly")
			assert(x.URL, "/api/v1/polyapi/request/system/poly", "FmtAPIInOut")
		}
	}
	for i, v := range cases {
		o := v.obj
		getSub := func(main, sub string) string {
			if sub != "" {
				return sub
			}
			return main
		}
		if x, ok := o.(Stringer); ok && v.expr != "" {
			x.SetString(v.expr)
			str := getSub(v.expr, v.String) // the same as expr
			if s := x.String(); s != str {
				t.Errorf("case %d String() fail: expect %s, got %s", i+1, str, s)
			}
		}

		if x, ok := o.(DelayedJSONDecoder); ok {
			if err := x.DelayedJSONDecode(); err != nil {
				t.Errorf("case %d DelayedJSONDecode error: %s", i+1, err.Error())
			}
		}

		if x, ok := o.(StringerWithError); ok && v.expr != "" {
			if err := x.SetStringWithError(v.expr); err != nil {
				t.Errorf("case %d SetStringWithError error: %s", i+1, err.Error())
			}
		}

		if x, ok := o.(NamedScriptElem); ok {
			if name := x.GetName(false); name != v.name {
				t.Errorf("case %d GetName(false) fail: expect %s, got %s", i+1, v.name, name)
			}
			tn := getSub(v.name, v.nameTitle)
			if title := x.GetName(true); title != tn {
				t.Errorf("case %d GetName(true) fail: expect %s, got %s", i+1, tn, title)
			}
		}

		if x, ok := o.(GenSampler); ok {
			if d := x.GenSample(nil, false); d == nil {
				t.Errorf("case %d GenSample() fail: got %v", i+1, d)
			}
		}
		if x, ok := o.(CreateSampleDataor); ok {
			if d := x.CreateSampleData(nil, false); d == nil {
				// t.Errorf("case %d GenSampleData() fail: got %v", i+1, d)
			}
		}

		if x, ok := o.(ScriptElem); ok {
			got, err := x.ToScript(1, nil)
			if got != v.expect || err != nil {
				t.Errorf("case %d ToScript error: expect %q got %q err=%v", i+1, v.expect, got, err)
				fmt.Println(got)
			}
		}

		verifyValueX(i, v, t)
	}
}

func TestFieldRefer(t *testing.T) {
	type testCase struct {
		name   string
		expect string
	}
	testCases := []*testCase{
		&testCase{"123", "123"},
		&testCase{"123.45", "123.45"},
		&testCase{"$123", "$123"},
		&testCase{"$123.45", "$123.45"},
		&testCase{"a", "a"},
		&testCase{"$.a", "d.a"},
		&testCase{"a.b", "d.a.b"},
		&testCase{"$.a.b", "d.a.b"},
		&testCase{"a.b", "d.a.b"},
		&testCase{"$$a", "d.a"},
		&testCase{"$$.a", "d.a"},
		&testCase{"$$.a.b", "d.a.b"},
		&testCase{"$a", "d.a"},
		&testCase{"a+b*c.x", "a+b*d.c.x"},
		&testCase{"a+(b*c.x)", "a+(b*d.c.x)"},

		&testCase{"a.$.b", "a.$.b"},
		&testCase{"a+(b*c.0.x)", "a+(b*c.0.x)"},
		&testCase{"a+(b*c.[0 ].x)", "a+(b*d.c[0].x)"},
		&testCase{"a+(b*c.[0].x)", "a+(b*d.c[0].x)"},
		&testCase{"a+(b*c.e[0].x)", "a+(b*d.c.e[0].x)"},
		&testCase{"a+(b*c.e[ 0 ].x)", "a+(b*d.c.e[ 0 ].x)"},
		&testCase{"a+(b*c.[0].[11].x)", "a+(b*d.c[0][11].x)"},
		&testCase{"a+(b*c.e[0][12].x)", "a+(b*d.c.e[0][12].x)"},
		&testCase{"a+(b*c.e[ 0 ][ 13 ].x)", "a+(b*d.c.e[ 0 ][ 13 ].x)"},
		&testCase{"($a.b.[0 ].c+$a.b.[ 1].c)*$b+c", "(d.a.b[0].c+d.a.b[1].c)*d.b+c"},
		&testCase{"($$$a.b.[0 ].c+$$a.b.[ 1].c)*$b+c.d", "(d.a.b[0].c+d.a.b[1].c)*d.b+d.c.d"},
	}

	for i, v := range testCases {
		got, err := fixFullFieldRefer(v.name)
		if got != v.expect || err != nil {
			fmt.Println(i, v.name, got, err)
			t.Errorf("case %d expect %q got %q err=%v", i+1, v.expect, got, err)
		}
	}
}
