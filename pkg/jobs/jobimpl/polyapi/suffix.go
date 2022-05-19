package polyapi

import (
	"encoding/json"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// name=foo => foo.p
func updateSuffix(data *dbRecord) (bool, error) {
	updated := false
	if !strings.ContainsRune(data.Name, '.') {
		oldPath := apipath.Join(data.Namespace, data.Name)
		data.Name += ".p"
		updated = true

		if len(data.Doc) > 5 {
			var d adaptor.APIDoc
			if err := json.Unmarshal([]byte(data.Doc), &d); err != nil {
				return updated, err
			}

			apiPath := apipath.Join(data.Namespace, data.Name)
			d.FmtInOut.SetAccessURL(apiPath)
			renamed.Poly.Add(oldPath, apiPath)

			b, err := json.Marshal(d)
			if err != nil {
				return updated, err
			}
			data.Doc = string(b)
			updated = true
		}
	}
	return updated, nil
}
