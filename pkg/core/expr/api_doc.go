package expr

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
)

// ValueDefine represents the value of an input.
type ValueDefine struct {
	InputValue          //value
	Key        bool     `json:"key,omitempty"`     // key field, for name check of referenced
	Default    string   `json:"default,omitempty"` // default value
	Mock       string   `json:"mock,omitempty"`    // mock value
	Enums      []string `json:"enums,omitempty"`   // valid value enum of this input
	Ranges     []string `json:"ranges,omitempty"`  // valid value ranges [min,max)[min,max)... of this input
}

// InputNodeDetail represents the detail of an input node.
type InputNodeDetail struct {
	Inputs []ValueDefine `json:"inputs,omitempty"` // input from header, path, body or uri(GET)
	Consts ValueSet      `json:"consts,omitempty"` // const values provide by arrange
}

// DelayedJSONDecode delay unmarshal flex json object
func (d *InputNodeDetail) DelayedJSONDecode() error {
	for i := 0; i < len(d.Inputs); i++ {
		p := &d.Inputs[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}

	for i := 0; i < len(d.Consts); i++ {
		p := &d.Consts[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

// OutputNodeDetail represents detail of an output node
type OutputNodeDetail struct {
	Header ValueSet      `json:"header,omitempty"` // output from header
	Body   Value         `json:"body,omitempty"`   // output from body
	Doc    []ValueDefine `json:"doc,omitempty"`    // output from body, for doc only
}

// DelayedJSONDecode delay unmarshal flex json object
func (d *OutputNodeDetail) DelayedJSONDecode() error {
	if err := d.Header.DelayedJSONDecode(); err != nil {
		return nil
	}
	if err := d.Body.DelayedJSONDecode(); err != nil {
		return nil
	}
	for i := 0; i < len(d.Doc); i++ {
		p := &d.Doc[i]
		if err := p.DelayedJSONDecode(); err != nil {
			return err
		}
	}
	return nil
}

//------------------------------------------------------------------------------

// APISampleInput is the sample input of an API
type APISampleInput struct {
	Header http.Header     `json:"header,omitempty"`
	Body   json.RawMessage `json:"body,omitempty"`
}

// APISampleOutput is the sample output of an API
type APISampleOutput struct {
	Header http.Header     `json:"header,omitempty"`
	Resp   json.RawMessage `json:"resp,omitempty"`
}

// FmtAPIInOut is the formated API input and output
type FmtAPIInOut struct {
	Method    string             `json:"method"`
	URL       string             `json:"url"`
	Input     InputNodeDetail    `json:"input"`
	Output    OutputNodeDetail   `json:"output"`
	SampleIn  [2]APISampleInput  `json:"sampleInput"`  // [0]normal [1]tilteFirst
	SampleOut [2]APISampleOutput `json:"sampleOutput"` // [0]normal [1]tilteFirst
}

// SetAccessURL update the api access path
func (d *FmtAPIInOut) SetAccessURL(apiPath string) {
	if len(apiPath) > 0 && apiPath[0] == '/' {
		apiPath = apiPath[1:]
	}
	d.URL = fmt.Sprintf(consts.APIRequestPath, apiPath)
}

// DelayedJSONDecode delay unmarshal flex json object
func (d *FmtAPIInOut) DelayedJSONDecode() error {
	if err := d.Input.DelayedJSONDecode(); err != nil {
		return err
	}
	if err := d.Output.DelayedJSONDecode(); err != nil {
		return err
	}

	return nil
}
