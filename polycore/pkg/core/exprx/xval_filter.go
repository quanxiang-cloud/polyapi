package exprx

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// special on map range in test mode
var testMode = false

// FieldMap represents a object filed mapping
type FieldMap map[string]string

// ToScript returns the script of this element represent
func (m FieldMap) ToScript(depth int) string {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	lineHead := GenLinehead(depth)
	buf.WriteString(fmt.Sprintf("%s{", ""))

	type kv struct {
		k string
		v string
	}
	kvs := make([]kv, 0, len(m))
	for k, v := range m {
		kvs = append(kvs, kv{k, v})
	}
	if testMode { // sort result in test mode
		sort.Slice(kvs, func(i, j int) bool { return kvs[i].k < kvs[j].k })
	}
	for i := 0; i < len(kvs); i++ {
		p := &kvs[i]
		buf.WriteString(fmt.Sprintf(`%s:"%s", `, p.k, p.v))
	}

	buf.WriteString(fmt.Sprintf("%s}", lineHead))
	return buf.String()
}

// ValFiltObj represents data filter expression
type ValFiltObj struct {
	Source FieldRef `json:"source"` // field of object or array type
	White  FieldMap `json:"white"`  // white list of field name mapping, oldName->newName
	Black  FieldMap `json:"black"`  // black list of field name mapping
	Filter CondExpr `json:"filter"` // data fielter for an array

	config   string // config code of js object
	filtFunc string // function for filt d
}

// TypeName returns name of the type
func (v ValFiltObj) TypeName() string { return XValTypeFilter.String() }

func (v ValFiltObj) getConfig(depth int) string {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	lineHead := GenLinehead(depth)
	buf.WriteString(fmt.Sprintf("%s{", ""))
	buf.WriteString(fmt.Sprintf(`%s%s:%s, `, lineHead, "white", v.White.ToScript(depth)))
	buf.WriteString(fmt.Sprintf(`%s%s:%s, `, lineHead, "black", v.Black.ToScript(depth)))
	buf.WriteString(fmt.Sprintf("%s}", lineHead))
	return buf.String()
}

func (v ValFiltObj) getFiltFunc(depth int) string {
	return fmt.Sprintf(`function (%s) { return %s }`,
		polyAllDataVarName, first(v.Filter.ToScript(depth, nil)))
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValFiltObj) DelayedJSONDecode() error {
	if err := v.Filter.DelayedJSONDecode(); err != nil {
		return err
	}
	v.config = v.getConfig(0)
	v.filtFunc = v.getFiltFunc(0)
	return nil
}

// GetName returns Name of the elem
func (v ValFiltObj) GetName(titleFirst bool) string {
	return v.Source.GetName(false)
}

// ToScript returns the script of this element represent
func (v ValFiltObj) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen*2))
	lineHead := GenLinehead(depth)
	buf.WriteString(fmt.Sprintf("%s(\n", consts.PDFiltObject))
	buf.WriteString(fmt.Sprintf("%s  %s,\n", lineHead, FullVarName(v.Source.String())))
	buf.WriteString(fmt.Sprintf("%s  %s,\n", lineHead, v.config))
	buf.WriteString(fmt.Sprintf("%s  %s\n", lineHead, v.filtFunc))
	buf.WriteString(fmt.Sprintf("%s)", lineHead))
	return buf.String(), nil
}

//------------------------------------------------------------------------------

// ValMergeObj merge multi objects as one. eg: {a,b}+{c,d} => {a,b,c,d}
type ValMergeObj []Value

// TypeName returns name of the type
func (v ValMergeObj) TypeName() string { return XValTypeMergeObj.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValMergeObj) DelayedJSONDecode() error {
	for i := 0; i < len(*v); i++ {
		p := &(*v)[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// ToScript returns the script of this element represent
func (v ValMergeObj) ToScript(depth int, e Evaler) (string, error) {
	var buf = bytes.NewBuffer(make([]byte, 0, defaultBufLen))
	lineHead := GenLinehead(depth + 1)
	buf.WriteString(fmt.Sprintf("%s(", consts.PDMergeObjs))
	for i := 0; i < len(v); i++ {
		p := &v[i]
		script, err := p.ToScript(depth+2, e)
		if err != nil {
			return "", err
		}
		buf.WriteString(fmt.Sprintf("\n%s  %s", lineHead, script))
		if i < len(v)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(fmt.Sprintf("\n%s)", lineHead))
	return buf.String(), nil
}

//------------------------------------------------------------------------------
