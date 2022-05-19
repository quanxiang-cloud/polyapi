package rule

import (
	"fmt"
	"testing"
)

func TestNameRule(t *testing.T) {
	type testCase struct {
		name   string
		expect bool
	}
	const tooLong = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890___"
	cases := []*testCase{
		&testCase{name: "", expect: true},
		&testCase{name: "0", expect: true},
		&testCase{name: "A", expect: true},
		&testCase{name: "a", expect: true},
		&testCase{name: "0_abCD", expect: true},
		&testCase{name: "_0123abCD", expect: true},
		&testCase{name: "-0123abCD", expect: true},
		&testCase{name: "_", expect: true},
		&testCase{name: "-", expect: true},
		&testCase{name: "β", expect: false},
		&testCase{name: "不允许中文", expect: false},
		&testCase{name: "_a/b", expect: false},
		&testCase{name: tooLong, expect: false},
	}
	for i, v := range cases {
		err := ValidateName(v.name, MaxNameLength, true)
		got := err == nil
		if got != v.expect {
			fmt.Println(i, v.name, err)
			t.Errorf("case %d name=%s err=%v, unexpected", i+1, v.name, err)
		}
	}
}
