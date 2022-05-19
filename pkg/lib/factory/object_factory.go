package factory

import (
	"fmt"
	"reflect"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/jsonx"
)

// NamedType define a type with type name
type NamedType interface {
	TypeName() string
}

// FlexObjFactory impliments a factory that can create multi products by type name
type FlexObjFactory struct {
	mp     map[string]reflect.Type
	sample map[string]interface{}
	name   string
}

// NewFlexObjFactory create a new factory
func NewFlexObjFactory(name string) *FlexObjFactory {
	return &FlexObjFactory{
		mp:     make(map[string]reflect.Type),
		sample: make(map[string]interface{}),
		name:   name,
	}
}

// MustReg register the creator by name, it panic if name is duplicate
func (f *FlexObjFactory) MustReg(v NamedType) {
	err := f.Reg(v)
	if err != nil {
		panic(err)
	}
}

// Clean clean the factroy
func (f *FlexObjFactory) Clean() *FlexObjFactory {
	f.mp = make(map[string]reflect.Type)
	f.sample = make(map[string]interface{})
	return f
}

// Reg register the creator by name
func (f *FlexObjFactory) Reg(v NamedType) error {
	if _, ok := f.mp[v.TypeName()]; ok {
		return fmt.Errorf("duplicate reg of %s,%#v", v.TypeName(), v)
	}
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	f.mp[v.TypeName()] = t
	f.sample[v.TypeName()] = v
	return nil
}

// MustCreate make product by name, it panic if errors
func (f *FlexObjFactory) MustCreate(name string) interface{} {
	d, err := f.Create(name)
	if err != nil {
		panic(err)
	}
	return d
}

// Create make product by name
func (f *FlexObjFactory) Create(name string) (interface{}, error) {
	t, ok := f.mp[name]
	if !ok {
		return nil, fmt.Errorf("product [%s] cannot create from factory %s", name, f.name)
	}
	return reflect.New(t).Interface(), nil
}

// DelayedUnmarshalFlexJSONObject create the real object and ummarshal it for FlexJSONObject
func (f *FlexObjFactory) DelayedUnmarshalFlexJSONObject(kind string, obj *jsonx.FlexJSONObject) error {
	if kind == "" { // ignore empty input
		obj.D = nil
		return nil
	}

	p, err := f.Create(kind)
	if err != nil {
		return err
	}

	return obj.DelayedUnmarshalJSON(p)
}

// CreateSample make sample product by name
func (f *FlexObjFactory) CreateSample(name string) (interface{}, error) {
	t, ok := f.sample[name]
	if !ok {
		return nil, fmt.Errorf("product [%s] cannot create sample from factory %s", name, f.name)
	}

	return t, nil
}
