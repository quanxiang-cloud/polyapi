package docview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"
)

func joinURL(host, url string) string {
	return host + url
}

// GetAPIDocView return the doc of raw API
func GetAPIDocView(doc *adaptor.APIDoc, schemaHost string, docType string, titleFirst bool) ([]byte, error) {
	f := &doc.FmtInOut
	var d *apiDocView
	so := &f.SampleOut[getSampleIndex(titleFirst)]
	switch dt := enumset.Enum(docType); dt {
	case DocTypeRaw:
		d = &apiDocView{
			URL:    joinURL(schemaHost, f.URL),
			Method: f.Method,
			Input:  f.Input,
			Output: f.Output,
		}

	case DocTypeSwag:
		return []byte(doc.Swagger), nil

	case DocTypeCurl, DocTypeJavascript, DocTypePython:
		input, err := genScriptDoc(dt, doc, schemaHost, titleFirst)
		if err != nil {
			return nil, err
		}
		d = &apiDocView{
			URL:    joinURL(schemaHost, f.URL),
			Method: f.Method,
			Input:  input,
			Output: prettyBody(so.Resp, ""),
			Desc:   "",
		}

	default:
		return nil, errcode.ErrInvalidDocType.FmtError(docType, DocTypeEnum.GetAll())
	}

	if d != nil {
		b, err := json.Marshal(d)
		if err != nil {
			return nil, err
		}
		//println(string(b))
		return b, err
	}
	return nil, nil
}

func genScriptDoc(docType enumset.Enum, doc *adaptor.APIDoc, schemaHost string, titleFirst bool) (string, error) {
	switch docType {
	case DocTypeCurl:
		return genCurlDoc(doc, schemaHost, docType.String(), titleFirst)
	case DocTypeJavascript:
		return genJavascriptDoc(doc, schemaHost, docType.String(), titleFirst)
	case DocTypePython:
		return genPythonDoc(doc, schemaHost, docType.String(), titleFirst)
	default:
		return "", errcode.ErrInvalidDocType.FmtError(docType, DocTypeEnum.GetAll())
	}
}

func genCurlDoc(doc *adaptor.APIDoc, host string, docType string, titleFirst bool) (string, error) {
	f := &doc.FmtInOut
	s := &f.SampleIn[getSampleIndex(titleFirst)]
	buf := bytes.NewBuffer(nil)
	buf.WriteString("curl ")
	url := joinURL(host, f.URL)
	data := ""
	if f.Method != expr.MethodGet.String() {
		data = prettyBody(s.Body, "")
		buf.WriteString(fmt.Sprintf(`-X %s `, f.Method))
	} else {
		url = fmt.Sprintf("%s?%s", url, httputil.BodyToQuery(string(s.Body)))
	}
	buf.WriteString(url)
	buf.WriteString(" \\\n")
	for k, v := range s.Header {
		for _, vv := range v {
			buf.WriteString(fmt.Sprintf(`  -H "%s: %s" \`, k, vv))
			buf.WriteByte('\n')
		}
	}
	if data != "" {
		buf.WriteString(fmt.Sprintf("  --data '%s' ", data))
	}

	//println(buf.String())

	return buf.String(), nil
}

// url = "http://127.0.0.1:8080/";
// xhr = new XMLHttpRequest();
// xhr.open("post", url, true);
// var data;
// xhr.setRequestHeader("Content-Type", "application/json");
// data = JSON.stringify({"key": "value"});
// xhr.send(data);
func genJavascriptDoc(doc *adaptor.APIDoc, host string, docType string, titleFirst bool) (string, error) {
	f := &doc.FmtInOut
	s := &f.SampleIn[getSampleIndex(titleFirst)]
	url := joinURL(host, f.URL)
	data := "null"
	if f.Method == expr.MethodGet.String() {
		url = fmt.Sprintf("%s?%s", url, httputil.BodyToQuery(string(s.Body)))
	} else {
		data = prettyBody(s.Body, "")
	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf(`url = "%s";`, url))
	buf.WriteString(fmt.Sprintf(`
xhr = new XMLHttpRequest();
xhr.open("%s", url, true);
`, strings.ToLower(f.Method)))
	for k, v := range s.Header {
		for _, vv := range v {
			buf.WriteString(fmt.Sprintf(`xhr.setRequestHeader("%s", "%s");`, k, vv))
			buf.WriteByte('\n')
		}
	}
	buf.WriteString(fmt.Sprintf(`data = %s;`, data))
	buf.WriteByte('\n')
	buf.WriteString("xhr.send(data);")

	//println(buf.String())

	return buf.String(), nil
}

func genPythonDoc(doc *adaptor.APIDoc, host string, docType string, titleFirst bool) (string, error) {
	f := &doc.FmtInOut
	s := &f.SampleIn[getSampleIndex(titleFirst)]
	url := joinURL(host, f.URL)
	data := "None"
	if f.Method == expr.MethodGet.String() {
		url = fmt.Sprintf("%s?%s", url, httputil.BodyToQuery(string(s.Body)))
	} else {
		data = prettyBody(s.Body, "")
	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(`import json
import requests
`)
	buf.WriteString("headers = {\n")
	for k, v := range s.Header {
		for _, vv := range v {
			buf.WriteString(fmt.Sprintf(`  '%s': '%s',`, k, vv))
			buf.WriteString("\n")
		}
	}
	buf.WriteString("}\n")

	buf.WriteString(fmt.Sprintf(`data = json.dumps(%s)`, data))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf(`r = requests.%s("%s", data=data, headers=headers)`, strings.ToLower(f.Method), url))
	buf.WriteString("\n")
	buf.WriteString("print(r.text)\n")

	//println(buf.String())

	return buf.String(), nil
}

func prettyBody(body json.RawMessage, endLine string) string {
	var d interface{}
	if err := json.Unmarshal(body, &d); err != nil {
		return string(body)
	}

	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return string(body)
	}

	s := string(b)
	if endLine != "" {
		s = strings.ReplaceAll(s, "\n", endLine)
	}
	return s
}
