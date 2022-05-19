package timestamp

import (
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
)

// Timestamp get current UTC timestamp by format
func Timestamp(format string) string {
	f := getTimeFormat(format)
	return UTCNow().Format(f)
}

// UTCNow return time.Now as UTC time
func UTCNow() time.Time {
	return time.Now().UTC()
}

func getTimeFormat(format string) string {
	f, _, ok := fmtEnumSet.Content(format)
	if !ok {
		return defaultTimeFormat
	}
	return f
}

// ValidateTimeFormat validate format of timestamp
func ValidateTimeFormat(format string) error {
	if !fmtEnumSet.Verify(format) {
		return errcode.ErrTimestampFormat.FmtError(format, supportedFormat)
	}
	return nil
}
