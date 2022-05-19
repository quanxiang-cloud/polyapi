package rawapi

import (
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// name=foo => foo.r
func updateSuffix(data *dbRecord) (bool, error) {
	if !strings.ContainsRune(data.Name, '.') {
		oldName := data.Name
		data.Name += ".r"
		apiPath := apipath.Join(data.Namespace, data.Name)
		data.Doc.FmtInOut.SetAccessURL(apiPath)
		renamed.Raw.Add(apipath.Join(data.Namespace, oldName), apiPath)
		return true, nil
	}
	return false, nil
}
