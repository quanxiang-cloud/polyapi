package schema

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"
)

type val struct {
	In   string      `json:"in"`
	Data interface{} `json:"data"`
}

// ParseRequest parse http request
// TODO: remove assertion
func ParseRequest(entity json.RawMessage, header http.Header) (json.RawMessage, error) {
	data := make(map[string]*val)
	json.Unmarshal(entity, &data)

	body := make(map[string]interface{})
	for k, v := range data {
		switch v.In {
		case consts.ParaInHeader:
			if err := parseHeader(header, k, v); err != nil {
				return nil, err
			}
		case consts.ParaInPath:
			if err := parsePath(body, k, v); err != nil {
				return nil, err
			}
		case consts.ParaInBody:
			if err := parseBody(body, k, v); err != nil {
				return nil, err
			}
		case consts.ParaInQuery:
			if err := parseQuery(body, k, v); err != nil {
				return nil, err
			}
		}
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(b), nil
}

func parseHeader(header http.Header, k string, v *val) error {
	if s, ok := v.Data.(string); ok {
		header.Set(k, s)
	} else {
		// error
		return fmt.Errorf("type of %v is not string", v.Data)
	}
	return nil
}

func parsePath(body map[string]interface{}, k string, v *val) error {
	if _, ok := body[polysign.XPolyBodyHideArgs]; !ok {
		body[polysign.XPolyBodyHideArgs] = make(map[string]interface{})
	}

	hideArgs := body[polysign.XPolyBodyHideArgs].(map[string]interface{})
	if _, ok := hideArgs[k]; !ok {
		hideArgs[k] = v.Data
		return nil
	}
	return fmt.Errorf("the key of %s is duplicate", k)
}

func parseBody(body map[string]interface{}, k string, v *val) error {
	bodyArgs, ok := v.Data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("paramter error")
	}

	for k, arg := range bodyArgs {
		if mp, ok := arg.(map[string]interface{}); ok {
			if _, ok := mp["data"].(map[string]interface{}); !ok {
				body[k] = mp["data"]
			} else {
				b, err := json.Marshal(arg)
				if err != nil {
					return err
				}

				subVal := &val{}
				err = json.Unmarshal(b, subVal)
				if err != nil {
					return err
				}
				var arg = make(map[string]interface{})
				parseBody(arg, "", subVal)
				body[k] = arg
			}
		}
	}
	return nil
}

func parseQuery(body map[string]interface{}, k string, v *val) error {
	if _, ok := body[k]; !ok {
		body[k] = v.Data
	}
	return fmt.Errorf("the key of %s is duplicate", k)
}
