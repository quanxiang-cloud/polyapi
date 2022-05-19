package rule

import (
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// const define
const (
	DescMaxLength = 4096 // desc max length
)

// Check rule
var (
	charSetExpr = regexp.MustCompile("^[\u4e00-\u9fa5a-zA-Z0-9,./;'\\[\\]\\\\<>?:\"{}|`~!@#$%^&*()_+-=\\s\\n，。；‘’【】、《》？：“”{}·~！￥…（）]*$")
)

// CheckCharSet check data whether contains unsupport character
func CheckCharSet(data ...string) error {
	for _, v := range data {
		if !charSetExpr.Match([]byte(v)) {
			return errcode.ErrCharacterSet.FmtError(v)
		}
	}
	return nil
}

// CheckDescLength check data length of desc
func CheckDescLength(data string) error {
	return CheckLength(data, DescMaxLength)
}

// CheckLength check data length
func CheckLength(data string, length int) error {
	if len(data) > length {
		return errcode.ErrCharacterTooLong.FmtError(data, length)
	}
	return nil
}
