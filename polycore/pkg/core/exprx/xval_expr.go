package exprx

import (
	"bytes"
)

// ValExprGroup represents expression group, eg: (b+c)
type ValExprGroup []ValExpr

// TypeName returns name of the type
func (v ValExprGroup) TypeName() string { return ExprTypeGroup.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValExprGroup) DelayedJSONDecode() error {
	for i := 0; i < len(*v); i++ {
		p := &(*v)[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValExprGroup) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	buf.WriteString("(")
	for i := 0; i < len(v); i++ {
		p := &v[i]
		if i > 0 {
			op, err := ConvertOp(p.Op)
			if err != nil {
				return "", err
			}
			buf.WriteString(op)
		}
		s, err := p.ToScript(depth, e)
		if err != nil {
			return "", err
		}
		buf.WriteString(s)
	}
	buf.WriteString(")")
	return buf.String(), nil
}

//------------------------------------------------------------------------------

// ValExpr represents an expression, eg: a.x + (b.y + c.z) * 2
type ValExpr struct {
	Op string `json:"op"` // operation, and|or|not add|sub|mul|div
	Value
}

// TypeName returns name of the type
func (v ValExpr) TypeName() string { return ExprTypeExpr.String() }
