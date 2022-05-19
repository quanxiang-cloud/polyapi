package schema

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestHttp(t *testing.T) {
	var a = `{
		"Content-Type": {
			"in": "header",
			"data": "application/json"
		},
		"body": {
			"in": "body",
			"data": {
				"name": {
					"in": "",
					"data": "name_value"
				},
				"obj": {
					"in": "",
					"data": {
						"obj_field": {
							"in": "",
							"data": "obj_value"
						}
					}
				}
			}
		},
		"path": {
			"in": "path",
			"data": false
		}
	}`
	// var a = `{
	// 	"body": {
	// 		"in": "body",
	// 		"data": "dsadsadsa"
	// 	}
	// }`

	header := http.Header{}

	body, err := ParseRequest([]byte(a), header)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%#v", body)

	b, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", string(b))

	fmt.Printf("%#v\n", header)
}
