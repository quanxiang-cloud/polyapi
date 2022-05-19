package polysign

import (
	"time"
)

// PolySignatureTimeout is the timeout for signature in polyapi
const PolySignatureTimeout = time.Second * 10

// header define
const (
	XHeaderPolyAccessToken = "Access-Token" // NOTE: Token authorize mode
)

// PolySignatureInfo is the data structure for signature generator
type PolySignatureInfo struct {
	AccessKeyID string `json:"X-Polysign-Access-Key-Id,omitempty"` // header
	Timestamp   string `json:"X-Polysign-Timestamp,omitempty"`     // header
	SignMethod  string `json:"X-Polysign-Method,omitempty"`        // header
	SignVersion string `json:"X-Polysign-Version,omitempty"`       // header

	Signature string `json:"x_polyapi_signature,omitempty"` // from body

	// NOTE: body XPolyRaiseUpFieldName defined, signature will raise up children for this field
	// Body:{Child:foo} will generate as Child=foo in query
	Body map[string]interface{} `json:"$$*_raise_*$$"`
}
