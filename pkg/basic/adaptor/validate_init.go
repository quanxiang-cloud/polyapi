package adaptor

import (
	"fmt"
	"reflect"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/validateinit"
)

func init() {
	// NOTE: check every field of inst is initialized
	validateinit.MustRegValidateFunc("adaptor", func() error {
		v := reflect.ValueOf(inst)
		t := reflect.TypeOf(inst)
		for i := v.NumField() - 1; i >= 0; i-- {
			// NOTE: f.CanInterface() is always false???
			if f := v.Field(i); f.IsNil() {
				return fmt.Errorf("%s uinistialized", t.Field(i).Name)
			}
		}
		return nil
	})
}
