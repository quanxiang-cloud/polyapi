package rule

import "github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"

// Array data maximum length
const (
	ArrayMax = 50
)

// CheckArrayLength check array data length
func CheckArrayLength(length int) error {
	if length > ArrayMax {
		return errcode.ErrExceedingMaximumLimit.FmtError(ArrayMax)
	}
	return nil
}
