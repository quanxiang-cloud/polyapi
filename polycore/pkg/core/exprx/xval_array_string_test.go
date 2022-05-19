package exprx

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestArrayString(t *testing.T) {
	type tt struct {
		A string
		B ValArrayString
	}
	type testCase struct {
		jsn    string
		expect string
		obj    ValArrayString
	}
	cases := []*testCase{
		&testCase{
			jsn:    `{"A":"x", "B":"foo,bar"}`,
			expect: `{"A":"x","B":"foo,bar"}`,
			obj:    ValArrayString{"foo", "bar"},
		},
		&testCase{
			jsn:    `{"A":"x?,y", "B":"foo?,bar"}`,
			expect: `{"A":"x?,y","B":"foo?,bar"}`,
			obj:    ValArrayString{`foo,bar`},
		},
		&testCase{
			jsn:    `{"A":"x??y", "B":"foo?bar"}`,
			expect: `{"A":"x??y","B":"foo?bar"}`,
			obj:    ValArrayString{`foo?bar`},
		},
		&testCase{
			jsn:    `{"A":"x??y", "B":"foo??bar"}`,
			expect: `{"A":"x??y","B":"foo??bar"}`,
			obj:    ValArrayString{`foo??bar`},
		},
	}
	for i, v := range cases {
		d := tt{}
		if err := json.Unmarshal([]byte(v.jsn), &d); err != nil {
			t.Errorf("case %d, json.Unmarshal error: %s", i+1, err)
		}
		if !reflect.DeepEqual(d.B, v.obj) {
			t.Errorf("case %d, Unmarshaled not equal, expect %#v got %#v", i+1, v.obj, d.B)
		}
		b, err := json.Marshal(d)
		if err != nil {
			t.Errorf("case %d, json.Marshal error: %s", i+1, err)
		}
		if got := string(b); got != v.expect {
			t.Errorf("case %d, expect %q got %q", i+1, v.expect, got)
		}
	}
}

func TestRawArrayString(t *testing.T) {
	type testCase struct {
		txt    string
		script string
		str    string
	}
	cases := []*testCase{
		&testCase{
			txt:    `"foo,bar"`,
			script: `["foo","bar"]`,
			str:    `"foo,bar"`,
		},
		&testCase{
			txt:    `foo,bar`,
			script: `["foo","bar"]`,
			str:    `"foo,bar"`,
		},
		&testCase{
			txt:    `"foo?,bar"`,
			script: `["foo,bar"]`,
			str:    `"foo?,bar"`,
		},
		&testCase{
			txt:    `foo?,bar`,
			script: `["foo,bar"]`,
			str:    `"foo?,bar"`,
		},
		&testCase{
			txt:    `foo??bar`,
			script: `["foo??bar"]`,
			str:    `"foo??bar"`,
		},
	}
	for i, v := range cases {
		var d ValArrayString
		d.SetString(v.txt)
		if got := d.String(); got != v.str {
			t.Errorf("case %d String error: expect %s got %s", i+1, v.str, got)
		}
		if got, err := d.ToScript(0, nil); got != v.script || err != nil {
			t.Errorf("case %d ToScript error: expect %s got %s err=%v", i+1, v.script, got, err)
		}
	}
}
