package jsonx

import (
	"encoding/json"
	"errors"
)

// RawMessage exports
type RawMessage = json.RawMessage

// FlexJSONObject is an object that can encoding/decoding JSON between flex Go types.
// It implements Marshaler and Unmarshaler and can delay JSON decoding
// from []byte into flex object.
type FlexJSONObject struct {
	D interface{} // flex object for JSON encoding and decoding
}

// MarshalJSON encoding field D as JSON.
func (f FlexJSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.D)
}

// UnmarshalJSON copy data into field D.
func (f *FlexJSONObject) UnmarshalJSON(data []byte) error {
	f.D = append(RawMessage(nil), data...)
	return nil
}

// DelayedUnmarshalJSON unmarshal []byte into instance d.
// It will update field D if unmarshal OK.
func (f *FlexJSONObject) DelayedUnmarshalJSON(d interface{}) error {
	if f.D == nil { //ignore nil inputs
		return nil
	}
	b, ok := f.D.(json.RawMessage)
	if !ok {
		return errors.New("FlexJSONObject: DelayedUnmarshalJSON on non json.RawMessage value")
	}

	//BUG: check miss if len(b)==0 check only
	switch s := string(b); s { // no data
	case "", `""`, `null`:
		f.D = nil
		return nil
	}
	if err := json.Unmarshal(b, d); err != nil {
		return err
	}
	f.D = d
	return nil
}

// Empty check if this object is empty
func (f FlexJSONObject) Empty() bool {
	return f.D == nil
}
