// Package validateinit validates variables that need init by other package.
package validateinit

import (
	"errors"
	"fmt"
)

// ValidateInit validate the vairables that is required by other package
func ValidateInit() error {
	return inst.validateAll()
}

// MustRegValidateFunc regist the validate func
func MustRegValidateFunc(name string, fn ValidateFunc) {
	if err := inst.reg(name, fn); err != nil {
		panic(err)
	}
}

// ValidateFunc is the validate func like "func Validate() error"
type ValidateFunc func() error

//------------------------------------------------------------------------------

type validateList []wantInit

var inst validateList

func (vl *validateList) reg(name string, fn ValidateFunc) error {
	if name == "" || fn == nil {
		return errors.New("validateinit missing name or fn")
	}
	*vl = append(*vl, wantInit{name: name, fn: fn})
	return nil
}

func (vl *validateList) validateAll() error {
	for _, v := range *vl {
		if err := v.fn(); err != nil {
			return fmt.Errorf("error validateinit %s:%s", v.name, err.Error())
		}
	}
	return nil
}

type wantInit struct {
	name string
	fn   ValidateFunc
}
