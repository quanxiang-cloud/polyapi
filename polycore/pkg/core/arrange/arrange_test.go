package arrange

import (
	"testing"

	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

// verify interfaces

var _ = []exprx.DelayedJSONDecoder{

	(*Node)(nil),
	(*InputNodeDetail)(nil),
	(*IfNodeDetail)(nil),
	(*RequestNodeDetail)(nil),
	(*OutputNodeDetail)(nil),

	//(*DispDemo)(nil),

}

var _ = []exprx.ScriptElem{
	// (*InputNodeDetail)(nil),
	// (*IfNodeDetail)(nil),
	// (*RequestNodeDetail)(nil),
	// (*OutputNodeDetail)(nil),

}

var _ = []exprx.NamedScriptElem{}

var _ = []exprx.NamedType{
	(*InputNodeDetail)(nil),
	(*IfNodeDetail)(nil),
	(*RequestNodeDetail)(nil),
	(*OutputNodeDetail)(nil),

	//(*DispDemo)(nil),
}

var _ = []exprx.Stringer{}

func TestExpr(t *testing.T) {
	type testCase struct {
		txt, expect string
	}
	cases := []*testCase{
		&testCase{
			`
"123"
`,
			`\n"123"\n`,
		},
		&testCase{
			`foo\\n`,
			`foo\\n`,
		},
		&testCase{
			"\rfoo\r\nbar\n",
			"\\rfoo\\r\\nbar\\n",
		},
	}
	for i, v := range cases {
		got := toJsString(v.txt)
		if got != v.expect {
			t.Errorf("case %d expect %q got %q", i+1, v.expect, got)
		}
	}
}
