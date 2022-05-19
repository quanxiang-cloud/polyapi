package arrange

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/validateinit"
)

// PolyDocGennerator is a function for generate swagger for the poly api
// this function is provide from outside package
var PolyDocGennerator func(poly *Arrange, in *InputNodeDetail,
	out *OutputNodeDetail) (*adaptor.APIDoc, error)

func init() {
	validateinit.MustRegValidateFunc("arrange", func() error {
		if PolyDocGennerator == nil {
			return fmt.Errorf("PolyDocGennerator uninitialized")
		}
		return nil
	})
}
