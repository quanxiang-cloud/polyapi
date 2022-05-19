package expr

import (
	"strconv"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/timestamp"
	"github.com/quanxiang-cloud/polyapi/pkg/core/value"
)

// Single values

// ValString represents string value
type ValString string

// TypeName returns name of the type
func (v ValString) TypeName() string { return ValTypeString.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValString) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v ValString) String() string { return string(v) }

// SetString set a string to Value
func (v *ValString) SetString(s string) { *v = ValString(s) }

// GenSample generate a sample JSON value
func (v ValString) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	var d value.String
	if v != "" {
		val = v.String()
	}
	if err := d.Set(val); err != nil {
		d = value.MockString()
	}
	return &d
}

//------------------------------------------------------------------------------

// ValNumber represents number value, deal as string
type ValNumber ValString

// TypeName returns name of the type
func (v ValNumber) TypeName() string { return ValTypeNumber.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValNumber) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v ValNumber) String() string { return string(v) }

// SetString set a string to Value
func (v *ValNumber) SetString(s string) { *v = ValNumber(s) }

// GenSample generate a sample JSON value
func (v ValNumber) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	if v != "" {
		if f, err := strconv.ParseFloat(v.String(), 64); err == nil {
			val = f
		}
	}
	var d value.Number
	if err := d.Set(val); err != nil {
		d = value.MockNumber()
	}
	return &d
}

//------------------------------------------------------------------------------

// ValBoolean represents boolean value
type ValBoolean ValString

// TypeName returns name of the type
func (v ValBoolean) TypeName() string { return ValTypeBoolean.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValBoolean) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v ValBoolean) String() string { return string(v) }

// SetString set a string to Value
func (v *ValBoolean) SetString(s string) { *v = ValBoolean(s) }

// GenSample generate a sample JSON value
func (v ValBoolean) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	var d value.Boolean
	if err := d.Set(val); err != nil {
		d = true
	}
	return &d
}

//------------------------------------------------------------------------------

// ValTimestamp represents timestamp value
type ValTimestamp ValString

// TypeName returns name of the type
func (v ValTimestamp) TypeName() string { return ValTypeTimestamp.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValTimestamp) DelayedJSONDecode() error {
	// NOTE: verify timestamp format
	return timestamp.ValidateTimeFormat(v.String())
}

// String convert v to string
func (v ValTimestamp) String() string { return string(v) }

// SetString set a string to Value
func (v *ValTimestamp) SetString(s string) { *v = ValTimestamp(s) }

// GenSample generate a sample JSON value
func (v ValTimestamp) GenSample(val interface{}, titleFirst bool) value.JSONValue {
	var d value.String
	if err := d.Set(val); err != nil {
		d = "2020-12-31T05:43:21+0800"
	}
	return &d
}

//------------------------------------------------------------------------------

// ValAction represents action string
type ValAction ValString

// TypeName returns name of the type
func (v ValAction) TypeName() string { return ValTypeAction.String() }

// DelayedJSONDecode delay unmarshal flex json object
func (v *ValAction) DelayedJSONDecode() error { return nil }

// String convert v to string
func (v ValAction) String() string { return string(v) }

// SetString set a string to Value
func (v *ValAction) SetString(s string) { *v = ValAction(s) }

//------------------------------------------------------------------------------
