package arrange

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

const exprIdent = `(?:[_a-zA-Z]\w*)`

var (
	regexSingleLine = regexp.MustCompile(`(?m-s:^.*$)`)
	regexToReplace  = regexp.MustCompile(`(?m-s:^(?P<HEAD>.*?:\s*)?(?P<EXPR>.*?)(?P<TAIL>\s*,\s*)?$)`)
	exprFields      = fmt.Sprintf(`(?sm:\W%s\.(?P<NODE>%s)(?:\.%s)*)`, consts.PolyAllDataVarName, exprIdent, exprIdent)
	regexFields     = regexp.MustCompile(exprFields)
)

/*
{
	a: d.req.a+d.req2.b:req3.x,
	b:"foo",
	c:d.req2.c,
}

--->

{
	a: sel(if1.yes,d.req.a,undefined),
	b:"foo",
	c:sel(!if2.yes&&if3.yes,d.req2.c+d.req3.d,undefined),
}
*/
func (b *jsBuilder) resolveEmptyNodeFieldRefer(expr string, e Evaler) (string, error) {
	var err error
	replaced := regexSingleLine.ReplaceAllStringFunc(expr, func(src string) string {
		if err != nil {
			return src
		}
		if elems := regexFields.FindAllStringSubmatch(src, -1); len(elems) > 0 {
			parts := regexToReplace.FindAllStringSubmatch(src, -1)[0]
			pre, script, pro := parts[1], parts[2], parts[3]
			//fmt.Printf("parts [%q] [%s] [%s] [%s]\n", parts[0], pre, script, pro)
			var branches branchRoutes
			for _, v := range elems {
				//fmt.Println("subs", v)
				node := v[1]
				n := b.findVisitedNode(node)
				if n == nil {
					err = fmt.Errorf("refer missing node %s", node)
					return src
				}
				branches.MergeFrom(n.branch)
			}
			if err = branches.Dedup(); err != nil {
				return src
			}
			if branches.Size() > 0 {
				return fmt.Sprintf("%s%s(%s,%s,%s)%s",
					pre, consts.PDsel, branches.getCond(), script, consts.ScriptValueUndefined, pro)
			}
		}

		return src
	})
	return replaced, err
}

//------------------------------------------------------------------------------

type branchRoute struct {
	IfNode string
	Yes    bool
}

func (b branchRoute) Cond(index int) string {
	prefix := ""
	if index > 0 {
		prefix = "&&"
	}
	urg := ""
	if !b.Yes {
		urg = "!"
	}
	// &&!ifNode.y
	return fmt.Sprintf("%s%s%s.%s.%s", prefix, urg, consts.PolyAllDataVarName, b.IfNode, consts.PolyIfNodeResultVar)
}

type branchRoutes []*branchRoute

// expr => sel(cond, expr, undefined)
func (b *branchRoutes) ResolveSel(expr string) string {
	if b.Size() == 0 {
		return expr
	}
	return fmt.Sprintf("%s(%s,  %s,  %s)",
		consts.PDsel, b.getCond(), expr, consts.ScriptValueUndefined)
}

func (b branchRoutes) getCond() string {
	buf := bytes.NewBuffer(nil)
	for i, v := range b {
		buf.WriteString(v.Cond(i))
	}
	return buf.String()
}

func (b branchRoutes) Size() int {
	return len(b)
}

func (b *branchRoutes) AddRoute(ifNode string, yes bool) {
	*b = append(*b, &branchRoute{IfNode: ifNode, Yes: yes})
}

func (b *branchRoutes) MergeFrom(other branchRoutes) {
	*b = append(*b, other...)
}

// Dedup remove the duplicate refer of if node branch,
// and return if has error refer both way of one if node
func (b *branchRoutes) Dedup() error {
	m := make(map[string]bool)
	for _, v := range *b {
		if yes, ok := m[v.IfNode]; ok {
			if yes != v.Yes {
				return fmt.Errorf("access both way of node:%s", v.IfNode)
			}
		} else {
			m[v.IfNode] = v.Yes
		}
	}
	*b = make([]*branchRoute, 0, len(m))
	for k, v := range m {
		b.AddRoute(k, v)
	}
	return nil
}

func (b *branchRoutes) Clean() {
	*b = nil
}
