package value

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	type testCase struct {
		obj      JSONValue
		val      []interface{}
		t        Type
		elemName string
		elem     interface{}
	}
	cases := []*testCase{
		&testCase{
			obj: new(Number),
			val: []interface{}{
				new(Number), new(float64), new(float32),
				new(int), new(int8), new(int16), new(int32), new(int64),
				new(uint), new(uint8), new(uint16), new(uint32), new(uint64),
				Number(10.1), float64(10.2), float32(10.3),
				int(10), int8(10), int16(10), int32(10), int64(10),
				uint(10), uint8(10), uint16(10), uint32(10), uint64(10),
			},
			t: TNumber,
		},
		&testCase{
			obj: new(Integer),
			val: []interface{}{
				new(Integer), new(float64), new(float32),
				new(int), new(int8), new(int16), new(int32), new(int64),
				new(uint), new(uint8), new(uint16), new(uint32), new(uint64),
				Integer(10), float64(10.2), float32(10.3),
				int(10), int8(10), int16(10), int32(10), int64(10),
				uint(10), uint8(10), uint16(10), uint32(10), uint64(10),
			},
			t: TInteger,
		},
		&testCase{
			obj: new(String),
			val: []interface{}{new(string), new(String),
				String("foo"), "foo"},
			t: TString,
		},
		&testCase{
			obj: new(Boolean),
			val: []interface{}{new(bool), new(Boolean), Boolean(true), true},
			t:   TBoolean,
		},
		&testCase{
			obj: new(Object),
			val: []interface{}{&Object{}, Object{}},
			t:   TObject,
		},
		&testCase{
			obj: new(Array),
			val: []interface{}{&Array{}, Array{}},
			t:   TArray,
		},
		&testCase{
			obj:      new(Object),
			val:      []interface{}{&Object{}, Object{}},
			elemName: "x",
			elem:     String("foo"),
			t:        TObject,
		},
		&testCase{
			obj:      new(Array),
			val:      []interface{}{&Array{}, Array{}},
			elemName: "x",
			elem:     "foo",
			t:        TArray,
		},
		&testCase{
			obj: new(Null),
			t:   TNull,
		},
	}
	for i, v := range cases {
		assert := func(got, expect interface{}, msg string) {
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("case %d Value(%s) fail: expect %v, got %v", i+1, msg, expect, got)
				panic("")
			}
		}
		assert(v.t.NewValue(2) != nil, true, v.t.String()+".NewValue()")
		assert(v.obj.Type(), v.t, v.t.String()+".Type()")
		if v.val != nil {
			for _, vv := range v.val {
				assert(v.obj.Set(vv), nil, v.t.String()+".Set()")
			}
			assert(v.obj.Set(Null{}), errNotSupport, v.t.String()+".Set()")
		} else {
			assert(v.obj.Set(nil), errNotSupport, v.t.String()+".Set()")
		}
		if v.elemName != "" {
			assert(v.obj.AddElement(v.elemName, v.elem), nil, v.t.String()+".AddElement()")
		} else {
			if v.t == TObject {
				assert(v.obj.AddElement("dummy", nil), notSupportErr("dummy", nil), v.t.String()+".AddElement()")
			} else {
				assert(v.obj.AddElement("dummy", nil), errNotSupport, v.t.String()+".AddElement()")
			}
		}
	}
}

func TestMockData(t *testing.T) {
	assert := func(got, expect interface{}, msg string) {
		if !reflect.DeepEqual(got, expect) {
			t.Errorf("TestMockData(%s) fail: expect %v, got %v", msg, expect, got)
		}
	}
	assert(MockString() != "", true, "MockString()")
	assert(MockNumber() < 20, true, "MockNumber()")
	assert(RandString("", 0) != "", true, `RandString("", 0)`)
}
