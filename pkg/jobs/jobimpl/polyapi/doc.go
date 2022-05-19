package polyapi

import (
	"encoding/json"
	"fmt"

	// BUG: polyapi missing doc
	_ "github.com/quanxiang-cloud/polyapi/polycore/pkg/core/polydoc"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/version"
	"github.com/quanxiang-cloud/polyapi/pkg/business/rule"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/arrange"
)

var (
	curVersion = version.FullVersion()
)

// rebuild api doc
func updateDoc(data *dbRecord) (bool, error) {
	if len(data.Arrange) < 5 || data.Doc == "" || getDocVersion(data.Doc) == curVersion {
		return false, nil
	}

	updated := false
	apiPath := apipath.Join(data.Namespace, data.Name)
	info := arrange.APIInfo{
		Namespace: data.Namespace,
		Name:      data.Name,
		Title:     data.Title,
		Desc:      data.Desc,
		Method:    data.Method,
	}
	switch _, doc, _, err := arrange.BuildJsScript(&info, data.Arrange, ""); {
	case err != nil:
		fmt.Printf("@@%s:%s\n", apiPath, err.Error())
		if data.Script != "" || data.Doc != "" {
			data.Active = rule.ActiveDisable
			data.Doc = ""
			data.Script = ""
			updated = true
		}

	default:
		data.Doc = doc
		updated = true
	}
	return updated, nil
}

type docType struct {
	Version string `json:"version"`
}

func getDocVersion(doc string) string {
	v := ""
	if doc != "" {
		var d docType
		if err := json.Unmarshal([]byte(doc), &d); err == nil {
			v = d.Version
		}
	}
	return v
}
