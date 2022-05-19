package schema

import (
	"encoding/json"
	"fmt"

	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

type schemaParser interface {
	parseField(d expr.ValueDefine) error
}

type schema map[string]interface{}

func newAPISchema(name string) schemaParser {
	s := newSchema()
	s[schemaName] = name
	s[schemaType] = docTypeObjectVal
	return s
}

func newSchema() schema {
	return make(map[string]interface{}, 4)
}

func (s schema) put(k string, v interface{}) error {
	if _, ok := s[k]; ok {
		// TODO: throw err, the key exsists
		return fmt.Errorf("the key is duplicate: %s", k)
	}
	s[k] = v
	return nil
}

func (s schema) getSubSchema(k string) schema {
	if _, ok := s[k]; !ok {
		s[k] = newSchema()
	}
	return s[k].(schema)
}

func (s schema) parseField(d expr.ValueDefine) error {
	dataSchema := newSchema()
	dataSchema.parse(d)
	k := d.Name
	if k == "" {
		k = d.In.String()
	}
	return s.getSubSchema(schemaProperties).put(k, dataSchema)
}

func (s schema) parse(d expr.ValueDefine) error {
	fmd, err := unmarshalData(d)
	if err != nil {
		return err
	}

	s[schemaTitle] = fmd[docTitle]
	s[schemaType] = fmd[docType]
	if v, ok := fmd[docIn]; ok {
		s[schemaIn] = v
	}

	data := fmd[docData]
	switch fmd[docType] {
	case docTypeObjectVal:
		s.parseObject(data)
	case docTypeArrayVal:
		s.parseArray(data)
	default:
	}
	return nil
}

func (s schema) parseObject(data interface{}) error {
	pd, err := unmarshalMap(data)
	if err != nil {
		return err
	}

	for _, v := range pd {
		if err := s.parseField(v); err != nil {
			return err
		}
	}
	return nil
}

func (s schema) parseArray(data interface{}) error {
	pd, err := unmarshalMap(data)
	if err != nil {
		return err
	}

	if len(pd) > 0 {
		item := pd[0]
		itemSchema := s.getSubSchema(schemaItems)
		itemSchema.parse(item)
		s[schemaSubType] = item.Type
	}
	return nil
}

func unmarshalMap(data interface{}) ([]expr.ValueDefine, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	pd := make([]expr.ValueDefine, 0)
	err = json.Unmarshal(b, &pd)
	if err != nil {
		return nil, err
	}
	return pd, nil
}

func unmarshalData(d expr.ValueDefine) (schema, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	s := newSchema()
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}
