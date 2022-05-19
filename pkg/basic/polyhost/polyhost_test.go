package polyhost

import (
	"fmt"
	"testing"
)

func TestPolyhost(t *testing.T) {
	type testCase struct {
		val string
		err bool
	}
	testCases := []*testCase{
		&testCase{"http://api.xxx.com", false},
		&testCase{"https://api.xxx.com", false},
		&testCase{"http://api.xxx.com:80", false},
		&testCase{"https://api.xxx.com:443", false},

		&testCase{"http://api.xxx.com:08", true},
		&testCase{"HTTP://api.xxx.com", true},
		&testCase{"HTTPS://api.xxx.com", true},
		&testCase{"rpc://api.xxx.com", true},
		&testCase{"rpc://api.xxx.com:80", true},
		&testCase{"foo", true},
	}
	for i, v := range testCases {
		err := SetSchemaHost(v.val)
		got := err != nil
		if got != v.err {
			fmt.Println(i+1, v.val, err)
			t.Errorf("case %d schemaHost=%q expect %v got %v", i+1, v.val, v.err, got)
		}
	}
}
