package arrange

import (
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

const testShowEnums = false

func TestEnum(t *testing.T) {
	if testShowEnums {
		fmt.Println("- nodeTypes   :", NodeTypeEnum.ShowAll())
		fmt.Println("- nodeNames   :", NodeNameEnum.ShowAll())
		fmt.Println("- valueTypes  :", exprx.ValTypeEnum.ShowAll())
		fmt.Println("- xvalueTypes :", exprx.XValTypeEnum.ShowAll())
		fmt.Println("- exprTypes   :", exprx.ExprTypeEnum.ShowAll())
		fmt.Println("- valueOpTypes:", exprx.OpEnum.ShowAll())
		fmt.Println("- condTypes   :", exprx.CondEnum.ShowAll())
		fmt.Println("- cmpTypes    :", exprx.CmpEnum.ShowAll())
		fmt.Println("- encodingFmt :", exprx.EncodingEnum.ShowAll())
		fmt.Println("- schemaTypes :", exprx.SchemaEnum.ShowAll())
		fmt.Println("- paraTypes   :", exprx.ParaTypeEnum.ShowAll())

		fmt.Println("- methodTypes :", exprx.MethodEnum.ShowAll())
	}
}

func TestNodeName(t *testing.T) {
	type testCase struct {
		name string
		err  bool
	}
	testCases := []*testCase{
		&testCase{"abc", false},
		&testCase{"abC", false},
		&testCase{"abC_4", false},
		&testCase{"ab_", false},
		&testCase{"_abc", false},

		&testCase{"-abc", true},
		&testCase{"1_abc", true},
		&testCase{"foo.bar", true},
		&testCase{"ab-c", true},
		&testCase{"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890", true},
	}
	for i, v := range testCases {
		err := ValidateNodeName(v.name)
		if got := err != nil; got != v.err {
			t.Errorf("case %d name=%q expect=%v got %v err=%v", i+1, v.name, v.err, got, err)
		}
	}
}
