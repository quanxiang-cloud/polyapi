package apipath

import (
	"fmt"
	"net/url"
	"testing"
)

func TestEscape(t *testing.T) {
	type testCase struct {
		full   string
		expect bool
	}
	cases := []*testCase{
		&testCase{full: "/K3Cloud/Kingdee.BOS.WebApi.ServicesStub.AuthService.ValidateUser.common.kdsvc", expect: true},
		&testCase{full: "/distributor.action", expect: true},
		&testCase{full: "/iaas/", expect: true},
		&testCase{full: "/", expect: true},
		&testCase{full: "/?create", expect: true},
		&testCase{full: "/iaas/?create", expect: true},
		&testCase{full: "/iaas?create", expect: true},
		&testCase{full: "/192.135.177.8", expect: true},
		&testCase{full: "/foo~bar/", expect: true},
		&testCase{full: "/foo.bar/", expect: true},
		&testCase{full: "/foo......bar/", expect: true},
		&testCase{full: "/foo..__--.bar/", expect: true},
		&testCase{full: "/_foo-bar/", expect: true},
		&testCase{full: "/-foo~bar/", expect: true},
		&testCase{full: "/--", expect: true},
		&testCase{full: "/__", expect: true},
		&testCase{full: "/foo/:bar/x", expect: true},
		&testCase{full: "/:foo/{bar}/*x", expect: true},
		&testCase{full: "/foo/:bar/*x?y", expect: true},
		&testCase{full: "/v1/some/resource/name:customVerb", expect: true}, // https://cloud.google.com/apis/design/custom_methods
		&testCase{full: "/foo-bar/", expect: true},
		&testCase{full: "/foo_bar/", expect: true},
		// /*
		&testCase{full: "/iaas??create", expect: false},
		&testCase{full: "/iaas?create?delete", expect: false},
		&testCase{full: "iaas", expect: false},
		&testCase{full: "iaas/?create", expect: false},
		&testCase{full: "/foo\nbar", expect: false},
		&testCase{full: "/foo\rbar", expect: false},
		&testCase{full: "/foo$bar", expect: false},
		&testCase{full: "/foo%bar", expect: false},
		&testCase{full: "/foo#bar", expect: false},
		&testCase{full: "/foo@bar", expect: false},
		&testCase{full: "/foo!bar", expect: false},
		&testCase{full: "/foo\bar", expect: false},
		&testCase{full: "/.foo.bar/", expect: false},
		&testCase{full: "/~foo.bar/", expect: false},
		&testCase{full: "/~~", expect: false},
		&testCase{full: "/..", expect: false},
		&testCase{full: "", expect: false},
		&testCase{full: "/foo/::bar/*x", expect: false},
		&testCase{full: "/foo/:bar/**x", expect: false},
		&testCase{full: `
			toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo long
			toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo long
			toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo long
			toooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo long`,
			expect: false},
		// */
	}
	for i, v := range cases {
		err := ValidateAPIPath(&v.full)
		got := (err == nil)
		if got != v.expect {
			fmt.Printf("ValidateAPIPath(\"%s\")=\"%v\", \"%v\" err=%v\n", v.full, v.expect, got, err)
			t.Errorf("case %d full=%s expect=%v/%v, mismatch",
				i+1, v.full, got, v.expect)
		}
	}

	es := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_+=-~:$@&"
	fmt.Println(es)
	fmt.Println(url.PathEscape(es))
	es2 := "()[]{}<>;\"'<>/?!#%^*|"
	fmt.Println(es2)
	fmt.Println(url.PathEscape(es2))
}
