package swagger

import (
	"encoding/json"
)

// SwagResponses represents response in swagger
type SwagResponses map[string]*SwagResponseObject

// SwagResponseObject represents response object in swagger
type SwagResponseObject struct {
	Desc   string             `json:"description"`
	Schama SwagValue          `json:"schema"`
	Header SwagResponseHeader `json:"headers"`
}

// SwagResponseHeader represents header in response object
type SwagResponseHeader map[string]*SwagResponseHeaderItem

// SwagResponseHeaderItem represents item in response header
type SwagResponseHeaderItem struct {
	Type string `json:"type"`
	Desc string `json:"description"`
}

//------------------------------------------------------------------------------

// SwagInputValue SwagInputValue
type SwagInputValue struct {
	Ref      string `json:"$ref,omitempty"`
	Name     string `json:"name"`                  // Common
	In       string `json:"in"`                    // Common
	Desc     string `json:"description,omitempty"` // Common
	Required bool   `json:"required,omitempty"`    // Common

	Schema *SwagSchema `json:"schema,omitempty"` // body only

	// others
	Type            string `json:"type,omitempty"`
	Format          string `json:"format,omitempty"`
	AllowEmptyValue bool   `json:"allowEmptyValue,omitempty"`
	//Default         interface{} `json:"default,omitempty"`
}

// SwagSchema SwagSchema
type SwagSchema struct {
	AllOf      []json.RawMessage      `json:"allOf,omitempty"` // todo merge shcema
	Ref        string                 `json:"$ref,omitempty"`
	Type       interface{}            `json:"type"` // Type maybe "array" or ["array", "null"]
	Format     string                 `json:"format,omitempty"`
	Title      string                 `json:"title,omitempty"`
	Desc       string                 `json:"description,omitempty"`
	Required   []string               `json:"required,omitempty"`   // object only
	Properties map[string]*SwagSchema `json:"properties,omitempty"` // object only
	Items      *SwagSchema            `json:"items,omitempty"`      // array only
	//Default    interface{}            `json:"default,omitempty"`
}

// SwagResponsesSchema represents the response objects
type SwagResponsesSchema struct {
	Success *SwagResponseObjectSchema `json:"200"`
}

// SwagResponseObjectSchema represents the response object
type SwagResponseObjectSchema struct {
	Desc   string             `json:"description"`
	Schama *SwagSchema        `json:"schema,omitempty"`
	Header SwagResponseHeader `json:"headers"`
}
