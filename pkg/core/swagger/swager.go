package swagger

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

// ParseSwagger parse swagger content as raw api list
func ParseSwagger(content []byte, cfg *APIServiceConfig) ([]*adaptor.RawAPIFull, error) {
	var parser SwagParser
	r, err := parser.ParseSwagger(content, cfg)
	return r, err
}

//------------------------------------------------------------------------------

// APIServiceConfig define the config of a service
type APIServiceConfig struct {
	Schema     string
	Host       string
	AuthType   string
	Namespace  string
	Service    string
	StoredURL  string // file url that in fileserver
	SingleName string // operation id for only 1 api swagger
	Version    string
}

func (c *APIServiceConfig) selectAPIName(operationID string, apiCount int) string {
	switch {
	case operationID != "":
		return operationID
	case apiCount <= 1:
		return c.SingleName // return config.SingleName if it contain only 1 api
	}
	return ""
}

func getPathWithAction(path, action string) string {
	if action != "" {
		return fmt.Sprintf("%s?%s", path, action)
	}
	return path
}
