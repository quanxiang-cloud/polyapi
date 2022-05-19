package polyapi

import (
	"encoding/json"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

func TestVersion(t *testing.T) {
	var doc = adaptor.APIDoc{
		Version: "123",
	}
	b, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var d docType
	if err := json.Unmarshal(b, &d); err != nil {
		panic(err)
	}
	if d.Version != doc.Version {
		t.Errorf("doc.version mismatch, expect %q got %q", doc.Version, d.Version)
	}
}
