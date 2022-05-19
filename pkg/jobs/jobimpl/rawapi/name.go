package rawapi

import (
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/hash"

	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// uuid api name 'raw_APzutFc4wq-ZBrdBbLOfnfIGus7fQShjhmx15QaUiKwp.r' => short name
func updateName(data *dbRecord) (bool, error) {
	if strings.HasPrefix(data.Name, "raw_") {
		id, err := hash.ShortIDWithError(0)
		if err != nil {
			return false, err
		}
		oldPath := apipath.Join(data.Namespace, data.Name)
		data.Name = id + ".r"
		apiPath := apipath.Join(data.Namespace, data.Name)
		data.Doc.FmtInOut.SetAccessURL(apiPath)
		renamed.Raw.Add(oldPath, apiPath)
		fmt.Printf("@rename@ %s => %s\n", oldPath, apiPath)
		return true, nil
	}
	return false, nil
}
