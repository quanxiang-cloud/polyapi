package apipath

import (
	"fmt"
	"reflect"
	"testing"
)

const testPrint = false

func TestJoin(t *testing.T) {
	type testCase struct {
		ns     string
		name   string
		expect string
	}
	cases := []*testCase{
		&testCase{ns: "", name: "a", expect: "/a"},
		&testCase{ns: "/a", name: "b", expect: "/a/b"},
		&testCase{ns: "a", name: "b", expect: "/a/b"},
		&testCase{ns: "", name: "", expect: "/"},
		&testCase{ns: "", name: "/", expect: "/"},
		&testCase{ns: "/", name: "", expect: "/"},
		&testCase{ns: "/", name: "/", expect: "/"},
	}
	for i, v := range cases {
		got := Join(v.ns, v.name)
		if testPrint && true {
			fmt.Printf("Join(\"%s\", \"%s\")=\"%s\"\n", v.ns, v.name, got)
		}
		if got != v.expect {
			t.Errorf("case %d ns=%s name=%s expect=%s got=%s",
				i+1, v.ns, v.name, v.expect, got)
		}
	}
}

func TestSplit(t *testing.T) {
	type testCase struct {
		full string
		ns   string
		name string
	}
	cases := []*testCase{
		&testCase{full: "a", ns: "/", name: "a"},
		&testCase{full: "/a", ns: "/", name: "a"},
		&testCase{full: "a/b", ns: "/a", name: "b"},
		&testCase{full: "/a/b", ns: "/a", name: "b"},
		&testCase{full: "/", ns: "/", name: ""},
		&testCase{full: "", ns: "/", name: ""},
	}
	for i, v := range cases {
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d TestSplit(%s) fail: expect %v, got %v", i+1, msg, expect, got)
			}
		}
		ns, name := Split(v.full)
		if testPrint && true {
			fmt.Printf("Split(\"%s\")=\"%s\", \"%s\"\n", v.full, v.ns, v.name)
		}
		assert(Name(v.full), name, "Name()")
		assert(Parent(v.full), ns, "Parent()")
		if ns != v.ns || name != v.name {
			t.Errorf("case %d full=%s ns=%s/%s name=%s/%s, mismatch",
				i+1, v.full, ns, v.ns, name, v.name)
		}
	}
}

func TestFormat(t *testing.T) {
	type testCase struct {
		full   string
		expect string
	}
	cases := []*testCase{
		&testCase{full: "", expect: "/"},
		&testCase{full: "a", expect: "/a"},
		&testCase{full: "/a", expect: "/a"},
		&testCase{full: "a/b", expect: "/a/b"},
		&testCase{full: "/a/b", expect: "/a/b"},
	}
	for i, v := range cases {
		got := Format(v.full)
		if testPrint && true {
			fmt.Printf("Format(\"%s\")=\"%s\", \"%s\"\n", v.full, v.expect, got)
		}
		if got != v.expect {
			t.Errorf("case %d full=%s format=%s/%s, mismatch",
				i+1, v.full, got, v.expect)
		}
	}
}

func TestGetAPIType(t *testing.T) {
	type testCase struct {
		full   string
		expect string
	}
	cases := []*testCase{
		&testCase{full: "", expect: ""},
		&testCase{full: "/a/b/c", expect: ""},
		&testCase{full: "a.r", expect: "r"},
		&testCase{full: "/a.x", expect: "x"},
		&testCase{full: "a/b.c", expect: "c"},
		&testCase{full: "/a.x/b", expect: ""},
	}
	for i, v := range cases {
		got := GetAPIType(v.full)
		if got != v.expect {
			fmt.Printf("APIType(\"%s\")=\"%s\", \"%s\"\n", v.full, v.expect, got)
			t.Errorf("case %d full=%s format=%s/%s, mismatch",
				i+1, v.full, got, v.expect)
		}
	}
}

func TestGetBaseName(t *testing.T) {
	type testCase struct {
		full   string
		expect string
	}
	cases := []*testCase{
		&testCase{full: "", expect: ""},
		&testCase{full: "/a/b/c", expect: "c"},
		&testCase{full: "a.r", expect: "a"},
		&testCase{full: "/a.x", expect: "a"},
		&testCase{full: "a/b.c", expect: "b"},
		&testCase{full: "/a.x/b", expect: "b"},
	}
	for i, v := range cases {
		got := BaseName(v.full)
		if got != v.expect {
			fmt.Printf("BaseName(\"%s\")=\"%s\", \"%s\"\n", v.full, v.expect, got)
			t.Errorf("case %d full=%s format=%s/%s, mismatch",
				i+1, v.full, got, v.expect)
		}
	}
}

func TestGenerateAPIName(t *testing.T) {
	type testCase struct {
		name   string
		ty     string
		expect string
	}
	cases := []*testCase{
		&testCase{name: "", ty: "r", expect: ".r"},
		&testCase{name: "foo", ty: "r", expect: "foo.r"},
		&testCase{name: "foo.r", ty: "r", expect: "foo.r"},
		&testCase{name: "foo.x", ty: "r", expect: "foo.x"},
	}
	for i, v := range cases {
		got := GenerateAPIName(v.name, v.ty)
		if got != v.expect {
			fmt.Printf("Format(\"%s\")=\"%s\", \"%s\"\n", v.name, v.expect, got)
			t.Errorf("case %d name=%s ty=%s format=%s/%s, mismatch",
				i+1, v.name, v.ty, got, v.expect)
		}
	}
}

func TestMakeRequestURL(t *testing.T) {
	type testCase struct {
		schema string
		host   string
		path   string
		expect string
	}
	cases := []*testCase{
		&testCase{schema: "HTTP", host: "api.xx.com", path: "", expect: "http://api.xx.com/"},
		&testCase{schema: "http", host: "api.xx.com", path: "/", expect: "http://api.xx.com/"},
		&testCase{schema: "HTTPS", host: "api.xx.com", path: "/api/v1/foo", expect: "https://api.xx.com/api/v1/foo"},
		&testCase{schema: "https", host: "api.xx.com", path: "api/v1/foo", expect: "https://api.xx.com/api/v1/foo"},
		&testCase{schema: "HTTP", host: "api.xx.com", path: "/iaas/", expect: "http://api.xx.com/iaas/"},
		&testCase{schema: "https", host: "api.xx.com", path: "/iaas/", expect: "https://api.xx.com/iaas/"},
		&testCase{schema: "RPC", host: "api.xx.com", path: "/api/v1/foo.p", expect: "rpc://api.xx.com/api/v1/foo.p"},
		&testCase{schema: "rpc", host: "api.xx.com", path: "api/v1/foo.r", expect: "rpc://api.xx.com/api/v1/foo.r"},
		&testCase{schema: "FILE", host: "api.xx.com", path: "/api/v1/foo.p", expect: "file://api.xx.com/api/v1/foo.p"},
		&testCase{schema: "file", host: "api.xx.com", path: "api/v1/foo.r", expect: "file://api.xx.com/api/v1/foo.r"},
	}
	for i, v := range cases {
		got := MakeRequestURL(v.schema, v.host, v.path)
		if got != v.expect {
			//fmt.Printf("MakeRequestURL(%q, %q, %q)=%q, %q\n", v.schema, v.host, v.path, got, v.expect)
			t.Errorf("case %d %q, %q, %q got=%q ~ %q, mismatch",
				i+1, v.schema, v.host, v.path, got, v.expect)
		}
	}
}
