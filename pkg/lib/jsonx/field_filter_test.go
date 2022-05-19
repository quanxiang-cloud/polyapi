package jsonx

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	vals := []interface{}{
		map[string]interface{}{
			"a": 1,
			"b": "xyz",
			"c": []string{"foo", "bar"},
			"d": map[string]interface{}{
				"a1": 1,
				"b1": "xyz",
				"c1": []string{"foo", "bar"},
			},
		},
		[]interface{}{
			map[string]interface{}{
				"a": 1,
				"b": "xyz",
				"c": []string{"foo", "bar"},
				"d": map[string]interface{}{
					"a1": 1,
					"b1": "xyz",
					"c1": []string{"foo", "bar"},
				},
			},
			map[string]interface{}{
				"a": 2,
				"b": "xyz2",
				"c": []string{"foo2", "bar2"},
				"d": map[string]interface{}{
					"a1": 3,
					"b1": "xyz3",
					"c1": []string{"foo3", "bar3"},
				},
			},
		},
		"hello",
		[]string{"foo", "bar"},
	}
	for i, v := range vals {
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		s := string(b)
		fmt.Println(i, s)
		{
			f1 := FiltJSON(s, "", "b:bb,d:dd,c1:cc1", "b,c,d")
			fmt.Println(i, f1)
		}
		{
			f1 := FiltJSON(s, "d", "b:bb,d:dd,c1:cc1", "b,c,d")
			fmt.Println(i, f1)
		}
	}
}
