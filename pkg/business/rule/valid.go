package rule

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/enums"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// Valid value
const (
	Valid   = 1
	Invalid = 0
)

// IsValid check a valid status
func IsValid(valid uint) bool {
	return valid == Valid
}

// CheckValid verify the valid state of object
func CheckValid(valid uint, op adaptor.Operation, obj enums.Enum) error {
	switch {
	case RequireValid(op, obj):
		if !IsValid(valid) {
			return errcode.ErrAPIInvalid.NewError()
		}
	default:
		// do nothing
	}
	return nil
}

// RequireValid check if an operation require valid
func RequireValid(op adaptor.Operation, obj enums.Enum) bool {
	switch op {
	case adaptor.OpRequest:
		return true
	default:
		// do nothing
	}
	return false
}
