package rule

import (
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// MaxHostLen is max length of host
const MaxHostLen = 128

var hostExp = regexp.MustCompile(`(?s-m:\A[-\w\.]+(:[1-9]\d{0,4})?\z)`)

// ValidateHost check host format
func ValidateHost(host string) error {
	if len(host) > MaxHostLen {
		return errcode.ErrTooLong.FmtError("host", MaxHostLen)
	}
	if !hostExp.MatchString(host) {
		return errcode.ErrInvalidHost.FmtError(host)
	}
	return nil
}
