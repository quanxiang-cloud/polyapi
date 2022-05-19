package adaptor

import (
	"database/sql/driver"
	"encoding/json"
)

// APISchemaOper is the interface for api schema operator
type APISchemaOper interface {
}

// Schema json schema
type Schema struct {
	Input  json.RawMessage `json:"input"`
	Output json.RawMessage `json:"output"`
}

// APISchema APISchema
type APISchema struct {
	ID        string `json:"id"`
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	CreateAt  int64  `json:"createAt"`
	UpdateAt  int64  `json:"updateAt"`
}

// Value marshal
func (c Schema) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Scan unmarshal
func (c *Schema) Scan(data interface{}) error {
	if err := json.Unmarshal(data.([]byte), c); err != nil {
		return err
	}
	return nil
}
