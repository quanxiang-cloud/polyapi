package polysign

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestSignature(t *testing.T) {
	fmt.Println(url.PathEscape(XHeaderPolySignVersion))
	fmt.Println(url.PathEscape(XHeaderPolySignMethod))
	fmt.Println(url.PathEscape(XHeaderPolySignKeyID))
	fmt.Println(url.PathEscape(XHeaderPolySignTimestamp))
	fmt.Println(url.PathEscape(XHeaderPolyAccessToken))

	fmt.Println(url.PathEscape(XBodyPolySignSignature))
	fmt.Println(url.PathEscape(XPolyBodyHideArgs))
}

func TestRaiseUpFiledTag(t *testing.T) {
	body := map[string]interface{}{"body": fmt.Sprintf(`"%s"`, `this is raise up field(body)`)}
	var testCase = PolySignatureInfo{
		Signature:   "this is XBodyPolySignSignature",
		AccessKeyID: "this is XHeaderPolySignKeyID",
		Timestamp:   "this is XHeaderPolySignTimestamp",
		SignMethod:  "this is XHeaderPolySignMethod",
		SignVersion: "this is XHeaderPolySignVersion",
		Body:        body,
	}
	b, err := json.MarshalIndent(testCase, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	var got map[string]interface{}
	if err := json.Unmarshal(b, &got); err != nil {
		panic(err)
	}
	verify := func(name string, expect interface{}) {
		s, ok := got[name]
		if !ok {
			t.Errorf("missing filed %s", name)
		}
		if !reflect.DeepEqual(expect, s) {
			t.Errorf("field %s mismatch. expect=%q got=%q", name, expect, s)
		}
	}
	verify(XHeaderPolySignVersion, testCase.SignVersion)
	verify(XHeaderPolySignMethod, testCase.SignMethod)
	verify(XHeaderPolySignKeyID, testCase.AccessKeyID)
	verify(XHeaderPolySignTimestamp, testCase.Timestamp)
	verify(XBodyPolySignSignature, testCase.Signature)
	verify(XPolyRaiseUpFieldName, body)
}
