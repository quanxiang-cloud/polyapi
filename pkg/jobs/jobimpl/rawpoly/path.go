package rawpoly

import (
	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"
)

// update namespace title
func updateAPIPath(data *dbRecord) (bool, error) {
	updated := false
	if n, ok := renamed.Raw.Query(data.RawAPI); ok {
		data.RawAPI = n
		updated = true
	}
	if n, ok := renamed.Poly.Query(data.PolyAPI); ok {
		data.PolyAPI = n
		updated = true
	}

	return updated, nil
}
