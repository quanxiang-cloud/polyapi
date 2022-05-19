package namespace

import (
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/business/app"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

var subTitle = map[string]string{
	"":                      "应用",
	"poly":                  "API编排",
	"raw":                   "原生API",
	"raw/customer":          "代理第三方API",
	"raw/inner":             "平台API",
	"raw/inner/form":        "表单模型API",
	"raw/inner/form/form":   "表单API",
	"raw/inner/form/custom": "模型API",
}

// update namespace title
func updateTitle(data *dbRecord) (bool, error) {
	updated := false
	if data.Title == "" {
		apiPath := apipath.Join(data.Parent, data.Namespace)
		_, sub, err := app.SplitAsAppPath(apiPath)
		if err == nil {
			if title, ok := subTitle[sub]; ok {
				//fmt.Printf("@%s %s\n", apiPath, title)
				data.Title = title
				updated = true
			} else {
				//fmt.Println("@@", apiPath, sub)
			}
		} else {
			fmt.Printf("*%s err=%v\n", apiPath, err)
		}
	}

	return updated, nil
}
