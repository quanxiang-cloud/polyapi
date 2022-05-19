package schema

import (
	"encoding/json"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

// Schema json schema
type Schema = adaptor.Schema

// Gen generate form schema
func Gen(name string, api *expr.FmtAPIInOut) (*Schema, error) {
	return genSchema(name, api)
}

func genSchema(name string, api *expr.FmtAPIInOut) (*Schema, error) {
	ret := &Schema{}

	input := newAPISchema(name)
	for _, v := range api.Input.Inputs {
		input.parseField(v)
	}
	in, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	ret.Input = in

	output := newAPISchema(name)
	for _, v := range api.Output.Doc {
		output.parseField(v)
	}
	out, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	ret.Output = out

	return ret, nil
}
