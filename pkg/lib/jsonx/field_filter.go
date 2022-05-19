package jsonx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// field seprater in white & black
const stringArraySeprater = ","

// FiltJSON filt json object with field white and black config.
// white & black are field nama seprate with comma(,)
func FiltJSON(data string, fieldName string, white, black string) string {
	var d interface{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		fmt.Println(err)
		return data
	}
	f := jsonFilter{
		data:  d,
		white: splitStringArray(white),
		black: splitStringArray(black),
	}
	if err := f.filtData(f.data, fieldName); err != nil {
		return data
	}
	if b, err := json.Marshal(f.data); err == nil {
		return string(b)
	}
	return data
}

// jsonFilter filt JSON object/array
type jsonFilter struct {
	data  interface{}
	white map[string]string
	black map[string]string
}

func (f *jsonFilter) filtData(data interface{}, fieldName string) error {
	switch t := data.(type) {
	case map[string]interface{}:
		f.filtObject(t, fieldName)
	case []interface{}:
		f.filtArray(t, fieldName)
	}
	return nil
}

func (f *jsonFilter) filtObject(data map[string]interface{}, field string) error {
	blackEmpty := len(f.black) == 0
	if field != "" {
		if t, ok := data[field]; ok {
			return f.filtData(t, "")
		}
	}
	for k := range data {
		w, isInWhite := f.isInWhite(k)
		if accept := (!blackEmpty && f.isInBlack(k)) || isInWhite; !accept {
			delete(data, k)
		} else {
			if isInWhite && w != k {
				data[w] = data[k]
				delete(data, k)
			}
		}
	}
	return nil
}

func (f *jsonFilter) filtArray(data []interface{}, field string) error {
	for i, v := range data {
		if d, ok := v.(map[string]interface{}); ok {
			if err := f.filtObject(d, field); err == nil {
				data[i] = d
			}
		}
	}
	return nil
}

func (f *jsonFilter) isInWhite(field string) (string, bool) {
	if v, ok := f.white[field]; ok {
		if v == "" {
			return field, true
		}
		return v, true
	}
	return "", false
}

func (f *jsonFilter) isInBlack(field string) bool {
	if _, ok := f.black[field]; ok {
		return true
	}
	return false
}

// file rename set,  k:kk,k2:kk2,k3,k4
func splitStringArray(s string) map[string]string {
	ss := strings.Split(s, stringArraySeprater)
	out := make(map[string]string, len(ss))
	for i := 0; i < len(ss); i++ {
		fn := strings.TrimSpace(ss[i])
		k, v := fn, ""
		if elems := strings.Split(fn, ":"); len(elems) >= 2 {
			k, v = elems[0], elems[1]
		}
		out[k] = v
	}
	return out
}
