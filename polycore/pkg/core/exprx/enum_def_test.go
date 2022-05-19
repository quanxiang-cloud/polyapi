package exprx

import (
	"reflect"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

func TestEnums(t *testing.T) {
}

func TestEnumSet(t *testing.T) {
	assert := func(got, expect interface{}, msg string) {
		if !reflect.DeepEqual(got, expect) {
			t.Errorf("TestEnumSet fail:%s, expect %v, got %v", msg, expect, got)
		}
	}
	x := newEnumSet(nil)
	x.Reg("foo")
	x.Reg("kind")
	x.Reg("bar")
	x.Reg("great")
	x.Reg("heart")

	enumset.FinishReg()
	assert(x.ShowAll(), "foo | kind | bar | great | heart", "ShowAll")
	assert(x.ShowSorted(), "bar | foo | great | heart | kind", "ShowSorted")
	assert(x.GetAll(), []string{"foo", "kind", "bar", "great", "heart"}, "GetAll")

	assert(x.Verify("not_exists"), false, "Verify(not_exists)")
	assert(x.Verify("foo"), true, "Verify(exists)")
}
