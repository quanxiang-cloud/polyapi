package xml

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/beevik/etree"
)

const valueType = "t"
const valueCount = "c"
const elemType = "e"

// Encoder implements a xml encoder
type Encoder struct{}

// Encode marshal a Go data into []byte, it use buf firstly.
func (e *Encoder) Encode(v interface{}, buf []byte, pretty bool) ([]byte, error) {
	d := etree.NewDocument()

	d.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	e.value(&d.Element, "root", reflect.ValueOf(v))

	if pretty {
		d.WriteSettings.UseCRLF = true
		d.Indent(2)
	}

	b := bytes.NewBuffer(buf)
	_, err := d.WriteTo(b)
	return b.Bytes(), err
}

// value encode a Go value recursively
func (e *Encoder) value(parent *etree.Element, name string, v reflect.Value) error {
	switch k := v.Kind(); k {
	case reflect.Int, reflect.Uint, reflect.Float32, reflect.Float64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String, reflect.Bool:
		elem := parent.CreateElement(name)
		elem.CreateAttr(valueType, typeName(k))
		elem.SetText(fmt.Sprintf("%v", v.Interface()))
		return nil

	case reflect.Array, reflect.Slice:
		l := v.Len()
		elem := parent.CreateElement(name)
		elem.CreateAttr(valueType, typeName(k))
		elem.CreateAttr(elemType, elemName(v.Type().Elem().Kind()))
		c := 0
		for i := 0; i < l; i++ {
			if err := e.value(elem, "e", v.Index(i)); err != nil {
				return err
			}
			c++
		}
		elem.CreateAttr(valueCount, fmt.Sprintf("%d", c))
		return nil

	case reflect.Struct:
		t := v.Type()
		elem := parent.CreateElement(name)
		elem.CreateAttr(valueType, typeName(k))
		c := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			f := t.Field(i)
			fv := v.Field(i)
			if err := e.value(elem, f.Name, fv); err != nil {
				return err
			}
			c++
		}
		elem.CreateAttr(valueCount, fmt.Sprintf("%d", c))
		return nil

	case reflect.Map:
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()
		if kt.Kind() != reflect.String {
			return fmt.Errorf("unsupported map key %s", kt.String())
		}
		keys := v.MapKeys()
		elem := parent.CreateElement(name)
		elem.CreateAttr(valueType, typeName(k))
		elem.CreateAttr(elemType, elemName(vt.Kind()))

		c := 0
		for _, key := range keys {
			if err := e.value(elem, key.String(), v.MapIndex(key)); err != nil {
				return err
			}
			c++
		}
		elem.CreateAttr(valueCount, fmt.Sprintf("%d", c))
		return nil

	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			return e.value(parent, name, v.Elem())
		}

		elem := parent.CreateElement(name)
		elem.SetText("null")
		return nil

	default:
		// do nothing
	}
	return fmt.Errorf("unsupported type %s", v.Type().String())
}
