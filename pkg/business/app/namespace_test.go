package app

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

func check(s string, msg string) {
	if s == "" {
		panic(fmt.Errorf("missing %s", msg))
	}
}

func TestAPPNamespace(t *testing.T) {
	check(relPathRaw3Party, "relPathRaw3Party")
	check(relPathPoly, "relPathPoly")
	check(appRoot, "appRoot")

	fmt.Println(SplitAsAppPath("/system/app"))
	fmt.Println(SplitAsAppPath("/system/app/xxx"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/poly"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/raw"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/raw/customer"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/raw/inner"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/raw/inner/form"))
	fmt.Println(SplitAsAppPath("/system/app/xxx/raw/inner/form/form"))
	fmt.Println(relPathRaw3Party)
	fmt.Println(relPathPoly)

	fmt.Println(Path("XXX", "raw"))
	fmt.Println(Path("XXX", "raw.3party"))

	fmt.Println(PathEnum.GetAll())

	fmt.Printf("%#v\n", strings.Split("x/y/$", "$"))
	fmt.Printf("%#v\n", strings.Split("x/y/$", "y"))

	for k, v := range predefinedPath {
		fmt.Printf("%-20s %-10s %s\n", k, v.Name, v.Path)
	}
	fmt.Println("----------------------------------------")
	fmt.Println(RootOfAPPs)
	fmt.Println("----------------------------------------")

	app := "appX"
	for _, v := range createAppPathOrder {
		fmt.Println(Path(app, v.String()))
	}

}

func TestAppPath(t *testing.T) {
	type testCase struct {
		owner     string
		op        adaptor.Operation
		path      string
		expectErr bool
	}
	cases := []*testCase{
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app/bar.r", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app//foo.r/bar", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app//foo.r/bar.r", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app/bar", true},
		&testCase{"foo", adaptor.OpAddRawAPI, "/system/app/bar", true},
		&testCase{"foo", adaptor.OpAddSub, "/system/app/bar", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/x/app/bar", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/x/bar/poly", true},
		&testCase{"system", adaptor.OpAddPolyAPI, "x/app/bar/poly", false},
		&testCase{"", adaptor.OpAddPolyAPI, "x/app/bar/poly", false},
		&testCase{"foo", adaptor.OpAddSub, "system/app/bar/form", true},
		&testCase{"foo", adaptor.OpAddSub, "/system/app/bar/form", true},
		&testCase{"foo", adaptor.OpAddPolyAPI, "system/app/bar/poly", false},

		&testCase{"foo", adaptor.OpAddRawAPI, "system/app/bar/customer", true},
		&testCase{"foo", adaptor.OpAddRawAPI, "system/app/bar/raw/customer", false},
		&testCase{"foo", adaptor.OpAddService, "system/app/bar/customer", true},
		&testCase{"foo", adaptor.OpAddService, "system/app/bar/raw/customer", false},
		&testCase{"foo", adaptor.OpAddSub, "system/app/bar/customer", true},
		&testCase{"foo", adaptor.OpAddSub, "system/app/bar/raw/customer", false},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app/bar/poly", false},
		&testCase{"foo", adaptor.OpAddRawAPI, "/system/app/bar/customer", true},
		&testCase{"foo", adaptor.OpAddRawAPI, "/system/app/bar/raw/customer", false},
		&testCase{"foo", adaptor.OpAddService, "/system/app/bar/customer", true},

		&testCase{"foo", adaptor.OpAddService, "/system/app/bar/raw/customer", false},
		&testCase{"foo", adaptor.OpAddPolyAPI, "system/app/bar/poly/x", false},
		&testCase{"foo", adaptor.OpAddRawAPI, "system/app/bar/customer/x", true},
		&testCase{"foo", adaptor.OpAddRawAPI, "system/app/bar/raw/customer/x", false},
		&testCase{"foo", adaptor.OpAddService, "system/app/bar/customer/x", true},
		&testCase{"foo", adaptor.OpAddService, "system/app/bar/raw/customer/x", false},
		&testCase{"foo", adaptor.OpAddPolyAPI, "/system/app/bar/poly/x", false},
		&testCase{"foo", adaptor.OpAddRawAPI, "/system/app/bar/customer/x", true},
		&testCase{"foo", adaptor.OpAddRawAPI, "/system/app/bar/raw/customer/x", false},
		&testCase{"foo", adaptor.OpAddService, "/system/app/bar/customer/x", true},
		&testCase{"foo", adaptor.OpAddService, "/system/app/bar/raw/customer/x", false},
	}
	for i, v := range cases {
		err := ValidateNamespace(v.owner, v.op, v.path)
		if got := err != nil; got != v.expectErr {
			t.Errorf("case %d, ValidateNamespace(%s,%v,%s) fail, expect %t, got %t,%v", i+1, v.owner, v.op, v.path, v.expectErr, got, err)
		}
	}
}

func TestAppServicePath(t *testing.T) {
	type testCase struct {
		owner     string
		path      string
		expectErr bool
	}
	cases := []*testCase{
		&testCase{"foo", "/system/app/bar", true},
		&testCase{"foo", "/system/app/bar", true},
		&testCase{"foo", "/system/app/bar", true},
		&testCase{"foo", "/x/app/bar", true},
		&testCase{"foo", "/system/x/bar/poly", true},
		&testCase{"system", "x/app/bar/poly", false},
		&testCase{"", "x/app/bar/poly", false},
		&testCase{"foo", "system/app/bar/customer/x", true},
		&testCase{"foo", "system/app/bar/customer/x", true},
		&testCase{"foo", "system/app/bar/poly/x/x", true},

		&testCase{"foo", "system/app/bar/customer/x/x", true},
		&testCase{"foo", "system/app/bar/customer/x/x", true},
		&testCase{"foo", "/system/app/bar/poly/x/x", true},
		&testCase{"foo", "/system/app/bar/customer/x/x", true},
		&testCase{"foo", "/system/app/bar/customer/x/x", true},

		&testCase{"foo", "system/app/bar/raw/customer/x/x", false},
		&testCase{"foo", "system/app/bar/raw/customer/x/x", false},
		&testCase{"foo", "/system/app/bar/raw/customer/x/x", false},
		&testCase{"foo", "/system/app/bar/raw/customer/x/x", false},
	}
	for i, v := range cases {
		err := ValidateServicePath(v.owner, adaptor.OpCreate, v.path)
		if got := err != nil; got != v.expectErr {
			t.Errorf("case %d, ValidateNamespace(%s,%s) fail, expect %t, got %t,%v", i+1, v.owner, v.path, v.expectErr, got, err)
		}
	}
}

func prettyShow(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(err, v)
	}
	fmt.Println(string(b))
}

func TestCreateAppPathInfo(t *testing.T) {
	creates := GetCreateAppPaths("x")
	expect := []NamespaceDefine{
		NamespaceDefine{
			Name:   "x",
			Title:  "应用",
			Parent: "/system/app",
			Path:   "/system/app/x",
		},
		NamespaceDefine{
			Name:   "poly",
			Title:  "API编排",
			Parent: "/system/app/x",
			Path:   "/system/app/x/poly",
		},
		NamespaceDefine{
			Name:   "raw",
			Title:  "原生API",
			Parent: "/system/app/x",
			Path:   "/system/app/x/raw",
		},
		NamespaceDefine{
			Name:   "faas",
			Title:  "函数服务",
			Parent: "/system/app/x/raw",
			Path:   "/system/app/x/raw/faas",
		},
		NamespaceDefine{
			Name:   "customer",
			Title:  "代理第三方API",
			Parent: "/system/app/x/raw",
			Path:   "/system/app/x/raw/customer",
		},
		NamespaceDefine{
			Name:   "default",
			Title:  "默认分组",
			Parent: "/system/app/x/raw/customer",
			Path:   "/system/app/x/raw/customer/default",
		},
		NamespaceDefine{
			Name:   "inner",
			Title:  "平台API",
			Parent: "/system/app/x/raw",
			Path:   "/system/app/x/raw/inner",
		},
		NamespaceDefine{
			Name:   "form",
			Title:  "表单模型API",
			Parent: "/system/app/x/raw/inner",
			Path:   "/system/app/x/raw/inner/form",
		},
	}
	if !reflect.DeepEqual(creates, expect) {
		//fmt.Printf("%#v\n", creates)
		prettyShow(creates)
		t.Errorf("GetCreateAppPaths mismatch")
	}
}

func TestMakeRequestPath(t *testing.T) {
	fmt.Println(MakeRequestPath("/"))
	fmt.Println(MakeRequestPath(""))
	fmt.Println(MakeRequestPath("/system/app"))
	fmt.Println(MakeRequestPath("system/app"))
}
