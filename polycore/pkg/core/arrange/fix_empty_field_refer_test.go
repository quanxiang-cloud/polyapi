package arrange

import (
	"fmt"
	"testing"
)

func TestFixScriptFields(t *testing.T) {
	s := `
{
	a: d.req.a+d.req2.b*req3.x,
	b:"foo",
	c:d.req2.c+req13.a,
	d:req123.c,
	"e":d.req123.c,
}
`
	var b = jsBuilder{
		visited: map[string]*Node{
			"req": &Node{
				branch: branchRoutes{
					&branchRoute{"c1", false},
					&branchRoute{"c2", true},
				},
			},
			"req2": &Node{
				branch: branchRoutes{
					&branchRoute{"c1", true},
					&branchRoute{"c4", true},
				},
			},
			"req123": &Node{
				branch: branchRoutes{
					&branchRoute{"c3", true},
					&branchRoute{"c4", true},
					&branchRoute{"c1", true},
				},
			},
		},
	}
	fmt.Println(b.resolveEmptyNodeFieldRefer(s, nil))
}
