package polyhost

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/validateinit"
)

var (
	polySchemaHost string // https://polyapi.quanxiang.dev:443
	polySchema     string // https
	polyHost       string // polyapi.quanxiang.dev:443
)

func init() {
	validateinit.MustRegValidateFunc("polyhost", func() error {
		if polySchemaHost == "" {
			return errors.New("polySchemaHost uninitialized")
		}
		return nil
	})
}

//------------------------------------------------------------------------------

// https://polyapi.quanxiang.dev:443
var regexScmemaHost = regexp.MustCompile(`(?sm:^(?P<SCHEMA>https?)://(?P<HOST>[-\w\.]+(:[1-9]\d{0,4})?)$)`)

// SetSchemaHost initialize schema host of polyapi
func SetSchemaHost(schemaHost string) error {
	elems := regexScmemaHost.FindAllStringSubmatch(schemaHost, 1)
	if len(elems) == 0 {
		return fmt.Errorf("invalid schemaHost config %q", schemaHost)
	}
	polySchemaHost = schemaHost
	polySchema, polyHost = elems[0][1], elems[0][2]
	return nil
}

// GetSchemaHost return polyapi schemahost
func GetSchemaHost() string {
	return polySchemaHost
}

// GetSchema return polyapi schema
func GetSchema() string {
	return polySchema
}

// GetHost return polyapi host
func GetHost() string {
	return polyHost
}
