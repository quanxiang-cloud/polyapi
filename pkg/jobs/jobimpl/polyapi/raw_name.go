package polyapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/jobs/jobimpl/renamed"

	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"
)

// update data.arrange.RawPath
func updateRawPath(data *dbRecord) (bool, error) {
	updated := false
	var a arrange.Arrange
	if err := json.Unmarshal([]byte(data.Arrange), &a); err != nil {
		return false, err
	}
	updateMap := make(map[string]string)
	for i := 0; i < len(a.Nodes); i++ {
		p := &a.Nodes[i]
		if p.Type == arrange.NodeTypeRequest {
			if err := p.DelayedJSONDecode(); err != nil {
				return false, err
			}
			if d, ok := p.Detail.D.(*arrange.RequestNodeDetail); ok {
				if d.RawPath != "" {
					if n, ok := renamed.Raw.Query(d.RawPath); ok {
						updateMap[d.RawPath] = n
						d.RawPath = n
						continue
					}
					if !strings.ContainsRune(d.RawPath, '.') {
						n := d.RawPath + ".r"
						updateMap[d.RawPath] = n
						d.RawPath = n
						continue
					}
				}
			}
		}
	}

	if len(updateMap) > 0 {
		//NOTE:don't marshal a to arrange, because it will lost UI data.
		for k, v := range updateMap {
			// BUG: a=>a.r will make a.r=>a.r.r , while "a"=>"a.r" is safe
			// apiPath is always with "APIpath" format
			data.Arrange = strings.Replace(data.Arrange, fmt.Sprintf("%q", k), fmt.Sprintf("%q", v), -1)
		}
		updated = true
	}

	return updated, nil
}
