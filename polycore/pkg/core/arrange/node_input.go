// define input

package arrange

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/exprx"
)

// InputNodeDetail represents the detail of an input node.
type InputNodeDetail exprx.InputNodeDetail

func (d *InputNodeDetail) base() *exprx.InputNodeDetail {
	return (*exprx.InputNodeDetail)(d)
}

// TypeName returns name of the type
func (d InputNodeDetail) TypeName() string { return NodeTypeInput.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (d *InputNodeDetail) DelayedJSONDecode() error {
	return d.base().DelayedJSONDecode()
}

//------------------------------------------------------------------------------

// BuildJsScript generate JS script of this node
func (d *InputNodeDetail) BuildJsScript(n *Node, b *jsBuilder, depth int) error {
	//do some check of input define
	if err := d.check(); err != nil {
		return err
	}

	if err := d.buildObject(n, b, depth); err != nil {
		return err
	}
	if err := d.buildPredefined(n, b, depth); err != nil {
		return err
	}
	return nil
}

func (d *InputNodeDetail) check() error {
	if err := d.checkInput(&d.Inputs, true); err != nil {
		return err
	}
	if err := d.checkValueset(&d.Consts, false); err != nil {
		return err
	}
	return nil
}

func (d *InputNodeDetail) checkValueset(s *exprx.ValueSet, header bool) error {
	for i := 0; i < len(*s); i++ {
		p := &(*s)[i]
		if err := d.checkInputValue(p, header); err != nil {
			return err
		}
	}
	return nil
}

func (d *InputNodeDetail) checkInput(inputs *[]ValueDefine, header bool) error {
	for i := 0; i < len(*inputs); i++ {
		p := &(*inputs)[i]
		if err := d.checkInputValue(&p.InputValue, header); err != nil {
			return err
		}
	}
	return nil
}

func (d *InputNodeDetail) checkValue(p *exprx.Value, header bool) error {
	if err := p.DenyFieldRefer(); err != nil {
		return err
	}
	if !exprx.ValTypeEnum.Verify(p.Type.String()) {
		return fmt.Errorf("input value %s(%s) type unsupported, valid: %v",
			p.Name, p.Type, exprx.ValTypeEnum.GetAll())
	}
	return nil
}

func (d *InputNodeDetail) checkInputValue(p *exprx.InputValue, header bool) error {
	if err := p.DenyFieldRefer(); err != nil {
		return err
	}
	if !exprx.ValTypeEnum.Verify(p.Type.String()) {
		return fmt.Errorf("input value %s(%s) type unsupported, valid: %v",
			p.Name, p.Type, exprx.ValTypeEnum.GetAll())
	}
	return nil
}

// build predefined values
func (d *InputNodeDetail) buildPredefined(n *Node, b *jsBuilder, depth int) error {
	varName := exprx.FullVarName(n.Name)
	b.buf.WriteString(fmt.Sprintf(`  %s.header = %s.header || {}`, varName, varName))
	b.writeLn()
	// b.buf.WriteString(fmt.Sprintf(`  %s.path = %s.path || {}`, varName, varName))
	// b.writeLn()
	for i := 0; i < len(d.Consts); i++ {
		p := (*exprx.XInputValue)(&d.Consts[i])
		in := ""
		switch {
		case p.In.IsHeader():
			in = ".header"
		case p.In.IsPath():
			in = ".path"
		case p.In.IsBody():
			in = ""
		default:
			continue
		}
		s, err := p.ToScript(1, b.evaler)
		if err != nil {
			return err
		}
		b.buf.WriteString(fmt.Sprintf(`  %s%s.%s = %s`, varName, in, p.GetName(false), s))
		b.writeLn()
	}
	return nil
}

// build object code
func (d *InputNodeDetail) buildObject(n *Node, b *jsBuilder, depth int) error {
	b.buf.WriteString(fmt.Sprintf("  %s = %s.body\n",
		exprx.FullVarName(n.Name), consts.PolyAPIInputVarName))
	b.writeLn()

	return nil
}

func (d *InputNodeDetail) checkHasInputData(input ValueSet, name string) bool {
	for i := 0; i < len(input); i++ {
		if p := &input[i]; p.Name == name {
			return !p.Empty()
		}
	}
	return false
}

func (d *InputNodeDetail) checkRequired(input ValueSet) error {
	for i := 0; i < len(d.Inputs); i++ {
		if p := &d.Inputs[i]; p.Required && !d.checkHasInputData(input, p.Name) {
			return errcode.ErrBuildInputRequiredButMissing.FmtError(p.Name)
		}
	}
	return nil
}
