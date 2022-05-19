package jsonx

import (
	"encoding/json"
	"testing"
)

func TestFlexJSONObject(t *testing.T) {
	type testCase struct {
		obj            *FlexJSONObject
		delayUnmarshal interface{}
		input          string
		expect         string
		expectErr      bool
	}
	cases := []*testCase{
		&testCase{
			obj: &FlexJSONObject{
				D: "foo",
			},
			expect: `"foo"`,
		},
		&testCase{
			obj: &FlexJSONObject{
				D: nil,
			},
			input:  `"foo"`,
			expect: `"foo"`,
		},
		&testCase{
			obj: &FlexJSONObject{
				D: json.RawMessage(`"foo"`),
			},
			delayUnmarshal: new(string),
			expect:         `"foo"`,
		},
		&testCase{
			obj: &FlexJSONObject{
				D: json.RawMessage(``),
			},
			delayUnmarshal: new(string),
			expect:         `null`,
		},
		&testCase{
			obj: &FlexJSONObject{
				D: `foo`,
			},
			delayUnmarshal: new(string),
			expect:         `"foo"`,
			expectErr:      true,
		},
		&testCase{
			obj: &FlexJSONObject{
				D: json.RawMessage(`123`),
			},
			delayUnmarshal: new(string),
			expect:         `123`,
			expectErr:      true,
		},
	}
	for i, v := range cases {
		if v.input != "" {
			if err := v.obj.UnmarshalJSON([]byte(v.input)); err != nil {
				t.Errorf("case %d, UnmarshalJSON() fail: %s", i+1, err.Error())
			}
		}
		if v.delayUnmarshal != nil {
			if err := v.obj.DelayedUnmarshalJSON(v.delayUnmarshal); err != nil {
				if !v.expectErr {
					t.Errorf("case %d, DelayedUnmarshalJSON() fail: %s", i+1, err.Error())
				}
			}
		}
		b, err := v.obj.MarshalJSON()
		if !v.expectErr && err != nil {
			t.Errorf("case %d, MarshalJSON() fail: %s", i+1, err)
		}
		if got := string(b); got != v.expect {
			t.Errorf("case %d, MarshalJSON() fail, expect %q got %q", i+1, v.expect, got)
		}
	}
}
