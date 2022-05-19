package exprx

import (
	"bytes"
	"fmt"
)

// ValExprFunc represents a value of function call
type ValExprFunc struct {
	Func  string  `json:"func"`  // function name
	Paras []Value `json:"paras"` // parameters
}

// TypeName returns name of the type
func (v ValExprFunc) TypeName() string { return ExprTypeFunc.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValExprFunc) DelayedJSONDecode() error {
	for i := 0; i < len(v.Paras); i++ {
		p := &v.Paras[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValExprFunc) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	buf.WriteString(fmt.Sprintf("%s(", v.Func))
	for i := 0; i < len(v.Paras); i++ {
		p := &v.Paras[i]
		if i > 0 {
			buf.WriteString(", ")
		}
		s, err := p.ToScript(0, e)
		if err != nil {
			return "", err
		}
		buf.WriteString(s)
	}
	buf.WriteString(")")
	return buf.String(), nil
}
