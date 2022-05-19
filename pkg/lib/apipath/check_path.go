package apipath

import (
	"fmt"
	"regexp"
)

// MaxAPIPathLength is max length of api path
const MaxAPIPathLength = 256

// \/[-:~_0-9a-zA-Z\.\*\/\?]*
//var apiPathExp = regexp.MustCompile(`(?s-m:\A\/[-:~_0-9a-zA-Z\.\*\/\?]*\z)`)
var apiPathExp = regexp.MustCompile(`(?s-m:\A\/(?:[:*]?[-\w][-~\.\w]*\/?)*(?:[\?:][-\w][-~\.\w]*)?\z)`)

// ValidateAPIPath check if api path is valid
func ValidateAPIPath(pathPtr *string) error {
	path := formatAPIPath(pathPtr)
	if size := len(path); size > MaxAPIPathLength {
		return fmt.Errorf("too long(%d/%d) %q", size, MaxAPIPathLength, path)
	}
	if !apiPathExp.MatchString(path) {
		return fmt.Errorf("invalid: %q", path)
	}
	return nil
}

var expReplacePathArg = regexp.MustCompile(`\{(?P<ARG>[-\w][-~\.\w]*)\}`)

// FormatAPIPath change /api/{arg} => /api/:arg
func formatAPIPath(pathPtr *string) string {
	*pathPtr = expReplacePathArg.ReplaceAllStringFunc(*pathPtr, func(src string) string {
		arg := expReplacePathArg.FindAllStringSubmatch(src, 1)[0][1]
		return ":" + arg
	})
	return *pathPtr
}
