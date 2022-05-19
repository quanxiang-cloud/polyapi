package rule

import (
	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/enums"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// active value
const (
	ActiveEnable  = 1
	ActiveDisable = 0
	ActiveAny     = -1
)

// IsActive check an active state
func IsActive(active uint) bool {
	return active == ActiveEnable
}

// ValidateActive verify the active state of object
func ValidateActive(active uint, op adaptor.Operation, obj enums.Enum) error {
	switch {
	case RequireActive(op, obj):
		if !IsActive(active) {
			return errcode.ErrActiveDisabled.NewError()
		}
	case DenyActive(op, obj):
		if IsActive(active) {
			return errcode.ErrActiveEnabled.NewError()
		}
	default:
		// do nothing
	}
	return nil
}

// RequireActive check if an operation require avtive
func RequireActive(op adaptor.Operation, obj enums.Enum) bool {
	switch op {
	case adaptor.OpRequest, adaptor.OpAddRawAPI, adaptor.OpAddPolyAPI,
		adaptor.OpAddService, adaptor.OpAddSub: // require active
		return true
	default:
		// do nothing
	}
	return false
}

// DenyActive check if an operation require none avtive
func DenyActive(op adaptor.Operation, obj enums.Enum) bool {
	switch op {
	case adaptor.OpDelete, adaptor.OpBuild: // require inactive
		switch obj {
		case enums.ObjectNamespace, enums.ObjectService:
			// do nothing(dont deny active)
		default:
			return true
		}
	default:
		// do nothing
	}
	return false
}

// ValidActive check if an operation is ok via current active state
func ValidActive(active uint, op adaptor.Operation, obj enums.Enum) bool {
	switch {
	case RequireActive(op, obj):
		return IsActive(active)
	case DenyActive(op, obj):
		return !IsActive(active)
	default:
		// do nothing
	}
	return true
}

// IgnoreCache check if this operation ignore cache
func IgnoreCache(op adaptor.Operation) bool {
	switch op {
	case adaptor.OpCreate:
		return true
	}
	return false
}
