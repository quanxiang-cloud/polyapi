package exprx

// import (
// 	"bytes"
// 	"errors"
// 	"fmt"
// )

// // ValArrayElem represents an array that output as xxx.n
// type ValArrayElem struct {
// 	Name  string  `json:"name"`
// 	Array []Value `json:"array"`
// }

// // TypeName returns name of the type
// func (v ValArrayElem) TypeName() Enum { return ValTypeArrayElem }

// // DelayedJSONDecode delay unmarshal flex json object
// func (v *ValArrayElem) DelayedJSONDecode() error {
// 	if v.Name == "" {
// 		return errors.New("value array_elem missing name")
// 	}
// 	return (*ValArray)(&v.Array).DelayedJSONDecode()
// }

// // ToScript returns the script of this element represent
// func (v ValArrayElem) ToScript(depth int) string {
// 	var buf = bytes.NewBuffer(make([]byte, 0, DefaultBufLen))

// 	lineHead := GenLinehead(depth)
// 	for i := 0; i < len(v.Array); i++ {
// 		p := &v.Array[i]
// 		buf.WriteString(fmt.Sprintf("%s\"%s.%d\": %s,\n", lineHead, v.Name, i+1, p.ToScript(0)))
// 	}
// 	return buf.String()
// }

// // GetName returns Name of the value.
// func (v ValArrayElem) GetName() string {
// 	return v.Name
// }
