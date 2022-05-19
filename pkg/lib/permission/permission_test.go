package permission

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPermit(t *testing.T) {
	ss := PermitBitALL.ToList()
	if !reflect.DeepEqual(ss, []string{
		"read", "execute", "create",
		"update", "delete", "grant",
	}) {
		t.Errorf("%v", ss)
	}

	fmt.Println(ParsePermits([]string{"x", "execute"}))
	fmt.Println(ParsePermits([]string{"read", "update", "execute"}))
}
