package rawapi

import (
	"fmt"

	//"github.com/quanxiang-cloud/polyapi/pkg/lib/encoding"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/version"
	"github.com/quanxiang-cloud/polyapi/pkg/core/swagger"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

var (
	swagCfg    swagger.APIServiceConfig
	curVersion = version.FullVersion()
)

// rebuild api doc
func updateDoc(data *dbRecord) (update bool, err error) {
	if data.Doc.Version == curVersion {
		return false, nil
	}

	switch list, err := swagger.ParseSwagger(data.Doc.Swagger, &swagCfg); {
	case err != nil:
		timestamp := time2.NowUnix()
		data.DeleteAt = &timestamp
		return true, fmt.Errorf("@[DEL]@swag:%s", err.Error())
	case len(list) > 0:
		data.Doc = list[0].Doc
		apiPath := apipath.Join(data.Namespace, data.Name)
		data.Doc.FmtInOut.SetAccessURL(apiPath)

		/*
			if true {
				j, _ := encoding.ToJSON(data.Doc, false)
				fmt.Println("ToJSON", j)
				obj, err := encoding.FromJSON(j)
				if err != nil {
					panic(err)
				}
				if true {
					msg := "YAML"
					s, err := encoding.ToYAML(obj, false)
					if err != nil {
						panic(err)
					}
					fmt.Println("To", msg, s)
					o, err := encoding.FromYAML(s)
					if err != nil {
						panic(err)
					}
					fmt.Printf("from %s %#v\n", msg, o)
					ss, err := encoding.ToJSON(o, false)
					if err != nil {
						fmt.Println(err)
						//panic(err)
					}
					fmt.Println(msg, ss)
				}
				if true {
					msg := "XML"
					s, err := encoding.ToXML(obj, false)
					if err != nil {
						panic(err)
					}
					fmt.Println(msg, s)
					o, err := encoding.FromXML(s)
					if err != nil {
						panic(err)
					}
					fmt.Printf("from %s %#v\n", msg, o)
					ss, err := encoding.ToJSON(o, false)
					if err != nil {
						fmt.Println(err)
						//panic(err)
					}
					fmt.Println(msg, ss)
				}
			}
		*/

		return true, nil
	}
	return false, nil
}
