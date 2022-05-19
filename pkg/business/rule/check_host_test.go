package rule

import (
	"testing"
)

func TestHost(t *testing.T) {
	type testCase struct {
		host   string
		expect bool
	}
	hosts := []*testCase{
		&testCase{"localhost", true},
		&testCase{"123.com", true},
		&testCase{"_123.com", true},
		&testCase{"-123.com", true},
		&testCase{"xx135.com", true},
		&testCase{"192.168.1.1:123", true},
		&testCase{"foo.bar.com:123", true},
		&testCase{"foo:bar.com:123", false},
		&testCase{"foo.bar.com:123a", false},
		&testCase{"foo.bar.com:-1", false},
		&testCase{"foo.bar.com:01234", false},
		&testCase{"foo.bar.com\n:01234", false},
		&testCase{"foo.bar.com:123456", false},
		&testCase{"/foo.bar.com:123", false},
		&testCase{"http://foo.bar.com:123", false},
	}
	for i, v := range hosts {
		if got := ValidateHost(v.host) == nil; got != v.expect {
			t.Errorf("case %d, checkHost(%s) fail, expect %t, got %t", i+1, v.host, v.expect, !v.expect)
		}
	}
}
