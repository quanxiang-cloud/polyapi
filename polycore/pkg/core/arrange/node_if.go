// if logic of arrange config

package arrange

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

// IfNodeDetail represents detail of an if node.
type IfNodeDetail struct {
	Cond exprx.CondExpr `json:"cond"` // check condition
	Yes  string         `json:"yes"`  // next node if condition is yes
	No   string         `json:"no"`   // next node if condition is no
}

// TypeName returns name of the type
func (d IfNodeDetail) TypeName() string { return NodeTypeIf.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (d *IfNodeDetail) DelayedJSONDecode() error {
	return d.Cond.DelayedJSONDecode()
}

//------------------------------------------------------------------------------

func (d *IfNodeDetail) buildCond(depth int, b *jsBuilder) (string, error) {
	s, err := d.Cond.ToScript(depth, b.evaler)
	if err != nil {
		return "", err
	}
	return b.resolveEmptyNodeFieldRefer(s, b.evaler)
}

// BuildJsScript generate JS script of this node
func (d *IfNodeDetail) BuildJsScript(n *Node, b *jsBuilder, depth int) error {
	if err := d.checkNextNodes(n, b); err != nil {
		return err
	}
	lineHead := exprx.GenLinehead(depth)
	b.buf.WriteString(fmt.Sprintf("%s%s = { %s: false, }\n", lineHead, expr.FullVarName(n.Name), consts.PolyIfNodeResultVar))
	cond, err := d.buildCond(depth, b)
	if err != nil {
		return errcode.ErrBuildNodeFail.FmtError(n.getNameAlias(), err.Error())
	}
	b.buf.WriteString(fmt.Sprintf("%sif (%s) {\n", lineHead, cond))
	b.buf.WriteString(fmt.Sprintf("%s  %s.%s = true\n", lineHead, expr.FullVarName(n.Name), consts.PolyIfNodeResultVar))
	if d.Yes != "" {
		if err := d.visitNode(d.Yes, n, true, b, depth+1); err != nil {
			return err
		}
	}
	if d.No != "" {
		b.buf.WriteString(fmt.Sprintf("%s}else{\n", lineHead))
		if err := d.visitNode(d.No, n, false, b, depth+1); err != nil {
			return err
		}
	}
	b.buf.WriteString(fmt.Sprintf("%s}\n", lineHead))

	return nil
}

func (d *IfNodeDetail) visitNode(name string, n *Node, yes bool, b *jsBuilder, depth int) error {
	nn := b.findNode(name)
	if nn == nil {
		return errcode.ErrBuildUnknownNextNode.FmtError(n.getNameAlias(), n.Type, name)
	}
	if err := b.checkVisitNode(nn); err != nil {
		return err
	}
	nn.branch.MergeFrom(n.branch)
	nn.branch.AddRoute(n.Name, yes)
	if err := b.buildNode(nn, depth); err != nil {
		return err
	}
	return nil
}

func (d *IfNodeDetail) checkNextNodes(n *Node, b *jsBuilder) error {
	if d.Yes == "" && d.No == "" {
		return errcode.ErrBuildIfNodeMissingBothYN.FmtError(n.getNameAlias())
	}

	return nil
}
