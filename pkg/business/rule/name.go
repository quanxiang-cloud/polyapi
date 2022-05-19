package rule

import (
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

const (
	// MaxNameLength exports
	MaxNameLength = consts.MaxNameLength
)

var regexNameRule = regexp.MustCompile(`(?sm:^[\w-]+$)`)

// ValidateName verify api name format
func ValidateName(name string, maxLength int, allowEmpty bool) error {
	if allowEmpty && name == "" {
		return nil
	}

	if maxLength > 0 && len(name) > maxLength {
		return errcode.ErrNameTooLong.FmtError(name, maxLength)
	}

	if !regexNameRule.MatchString(name) {
		return errcode.ErrNameInvalid.FmtError(name)
	}
	return nil
}
