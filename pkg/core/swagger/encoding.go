package swagger

import (
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding"
)

// ConvertToJSON converts given encoding text to JSON.
func ConvertToJSON(encode string, source string, pretty bool) (string, error) {
	return encoding.ConvertEncoding(encode, source, expr.EncodingJSON.String(), pretty)
}
