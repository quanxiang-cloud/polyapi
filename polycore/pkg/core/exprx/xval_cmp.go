package exprx

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// CondExpr represents a condition expression
type CondExpr = ValExpr

// CondExprGroup represents a condition expression group
type CondExprGroup = ValExprGroup

//------------------------------------------------------------------------------

// ValExprCmp represents a compare expression, eg: a eq b
type ValExprCmp struct {
	LValue Value  `json:"lvalue"` // Left value
	Cmp    string `json:"cmp"`    // ""|lt|gt|le|ge|eq|ne
	RValue Value  `json:"rvalue"` // Right value
}

// TypeName returns name of the type
func (v ValExprCmp) TypeName() string { return ExprTypeCmp.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValExprCmp) DelayedJSONDecode() error {
	if err := v.LValue.DelayedJSONDecode(); err != nil {
		return err
	}
	if err := v.RValue.DelayedJSONDecode(); err != nil {
		return err
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValExprCmp) ToScript(depth int, e Evaler) (string, error) {
	ls, err := v.LValue.ToScript(depth, e)
	if err != nil {
		return "", err
	}
	rs, err := v.RValue.ToScript(depth, e)
	if err != nil {
		return "", err
	}
	op, err := ConvertOp(v.Cmp)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", ls, op, rs), nil
}

//------------------------------------------------------------------------------

// ValExprSel represents a select expression, eg: cond ? yesVal : noVal
type ValExprSel struct {
	Cond CondExpr `json:"cond"` // check condition
	Yes  Value    `json:"yes"`  // yes value
	No   Value    `json:"no"`   // no value
}

// TypeName returns name of the type
func (v ValExprSel) TypeName() string { return ExprTypeSel.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValExprSel) DelayedJSONDecode() error {
	if err := v.Cond.DelayedJSONDecode(); err != nil {
		return err
	}
	if err := v.Yes.DelayedJSONDecode(); err != nil {
		return err
	}
	if err := v.No.DelayedJSONDecode(); err != nil {
		return err
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValExprSel) ToScript(depth int, e Evaler) (string, error) {
	cs, err := v.Cond.ToScript(depth, e)
	if err != nil {
		return "", err
	}
	ys, err := v.Yes.ToScript(depth, e)
	if err != nil {
		return "", err
	}
	ns, err := v.No.ToScript(depth, e)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%s, %s, %s)", consts.PDSelect, cs, ys, ns), nil
}

//------------------------------------------------------------------------------
