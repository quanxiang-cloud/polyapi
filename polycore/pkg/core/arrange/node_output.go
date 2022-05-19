// define output

package arrange

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

// OutputNodeDetail represents detail of an output node
type OutputNodeDetail exprx.OutputNodeDetail

func (d *OutputNodeDetail) base() *exprx.OutputNodeDetail {
	return (*exprx.OutputNodeDetail)(d)
}

// TypeName returns name of the type
func (d OutputNodeDetail) TypeName() string { return NodeTypeOutput.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (d *OutputNodeDetail) DelayedJSONDecode() error {
	return d.base().DelayedJSONDecode()
}

// GenerateDoc generate doc from header & body
func (d *OutputNodeDetail) GenerateDoc() error {
	if len(d.Doc) > 0 {
		return nil
	}
	if err := d.genDocFromHeader(); err != nil {
		return err
	}
	if err := d.genDocFromBody(); err != nil {
		return err
	}
	return nil
}

func (d *OutputNodeDetail) genDocFromHeader() error {
	for _, v := range d.Header {
		if !v.In.IsHeader() {
			continue
		}
		vd := ValueDefine{
			InputValue: v,
		}
		d.Doc = append(d.Doc, vd)
	}
	return nil
}

func (d *OutputNodeDetail) genDocFromBody() error {
	body := &d.Body
	switch data := body.Data.D.(type) {
	case *ValObject:
		for _, v := range *data {
			vd := ValueDefine{
				InputValue: InputValue{
					Name:     v.GetName(false),
					Type:     expr.Enum(exprx.GetDocType(v.Type)),
					Title:    v.Title,
					Desc:     v.Desc,
					Appendix: v.Appendix,
					Required: v.Required,
					In:       expr.ParaTypeBody,
					Field:    expr.FieldRef(v.Field),
					Data:     v.Data,
				},
			}
			d.Doc = append(d.Doc, vd)
		}
	}
	return nil
}

//------------------------------------------------------------------------------

// ToScript ToScript
func (d OutputNodeDetail) buildOutput(depth int, b *jsBuilder) (string, error) {
	s, err := d.Body.ToScript(depth, b.evaler)
	if err != nil {
		return "", err
	}
	return b.resolveEmptyNodeFieldRefer(s, b.evaler)
}

// BuildJsScript generate JS script of this node
func (d *OutputNodeDetail) BuildJsScript(n *Node, b *jsBuilder, depth int) error {
	b.writeLn()
	output, err := d.buildOutput(depth, b)
	if err != nil {
		return errcode.ErrBuildNodeFail.FmtError(n.getNameAlias(), err.Error())
	}
	b.buf.WriteString(fmt.Sprintf("  %s = %s\n",
		exprx.FullVarName(n.Name), output))
	b.buf.WriteString(fmt.Sprintf("  return %s(%s)\n",
		consts.PDToJSONP, exprx.FullVarName(n.Name)))
	return nil
}
