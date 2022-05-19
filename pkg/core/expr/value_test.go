package expr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"
)

// verify interfaces

var _ = []DelayedJSONDecoder{
	(*InputValue)(nil),
	(*Value)(nil),
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	//(*ValArrayElem)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),

	(*ValueSet)(nil),
}

var _ = []NamedType{
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),
}

var _ = []GenSampler{
	(*ValTimestamp)(nil),
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValObject)(nil),
	(*ValArray)(nil),
}

var _ = []Stringer{
	(*ValNumber)(nil),
	(*ValString)(nil),
	(*ValBoolean)(nil),
	(*ValTimestamp)(nil),
	(*ValAction)(nil),
}

var _ = []json.Marshaler{
	(*FlexJSONObject)(nil),
}

var _ = []json.Unmarshaler{
	(*FlexJSONObject)(nil),
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
			obj:    newFieldRef(),
			expr:   `req1.x`,
			expect: `d.req1.x`,
			name:   `x`,
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
			obj:    &InputValue{},
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
			obj: &InputValue{
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
			obj: &InputValue{
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
			obj: &InputValue{
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
			obj: &ValueSet{
				InputValue{
					Name: "x",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"foo.bar"`),
					},
					In: "body",
				},
				InputValue{
					Name: "t",
					Type: "timestamp",
					Data: FlexJSONObject{
						D: json.RawMessage(`""`),
					},
					In: "body",
				},
				InputValue{
					Name: "y",
					Type: "string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"bar"`),
					},
					In: "body",
				},
				InputValue{
					Name: "z",
					Type: "array_string",
					Data: FlexJSONObject{
						D: json.RawMessage(`"foo,bar"`),
					},
					In: "query",
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
			fmt.Println(string(args.Body))
			fmt.Println(httputil.BodyToQuery(string(args.Body)))

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

		if x, ok := o.(NamedValue); ok {
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

		verifyValueX(i, v, t)
	}
}
