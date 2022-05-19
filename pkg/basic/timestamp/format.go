package timestamp

import (
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
)

var supportedFormat string

func init() {
	enumset.FinishReg()

	var buf strings.Builder
	buf.WriteString("[")
	for i, v := range fmtEnumSet.GetAll() {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteByte('\'')
		buf.WriteString(v)
		buf.WriteByte('\'')
	}
	buf.WriteString("]")
	supportedFormat = buf.String()
}

const defaultTimeFormat = "2006-01-02T15:04:05Z"

// time format defines
var (
	fmtEnumSet  = enumset.New(nil)
	fmtDefault  = fmtEnumSet.MustRegWithContent("", defaultTimeFormat, "")
	fmtDefault2 = fmtEnumSet.MustRegWithContent("YYYY-MM-DDThh:mm:ssZ", defaultTimeFormat, "")
	fmtISO8601  = fmtEnumSet.MustRegWithContent("YYYY-MM-DDThh:mm:ss+0000", "2006-01-02T15:04:05-0700", "")
)
