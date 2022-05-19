package swagger

import (
	"fmt"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

// GetSampleIndex return sample doc index
func GetSampleIndex(titleFirst bool) int {
	if titleFirst {
		return consts.SampleIndexTitle
	}
	return consts.SampleIndexNormal
}

func selScheme(ss []string) string {
	first := ""
	for _, v := range ss {
		if s, err := consts.ValidSchema(v); err == nil {
			switch {
			case consts.IsDefaultSchema(s):
				return s
			case first == "":
				first = s
			}
		}
	}
	return first
}

func selEncoding(es []string) string {
	first := ""
	for _, v := range es {
		if enc, err := consts.FromMIME(strings.ToLower(v)); err == nil {
			switch {
			case consts.IsDefaultEncoding(enc): // prefer
				return enc
			case enc == consts.EncodingYAML: // second prefer
				return enc
			case first == "":
				first = enc
			}
		}
	}
	if first != "" {
		return first
	}
	return consts.DefaultEncoding
}

// selPriority select string by priority
func selPriority(a ...string) string {
	for _, v := range a {
		if v != "" {
			return v
		}
	}
	return ""
}

// SwagConstValue export
type SwagConstValue = expr.SwagConstValue

func findVal(vals []SwagConstValue, exists *SwagConstValue, merge bool) (SwagConstValue, bool) {
	for i := 0; i < len(vals); i++ {
		if p := &vals[i]; p.Name == exists.Name {
			if merge {
				merged := SwagConstValue{
					Data: exists.Data,
					Name: exists.Name,
					In:   p.In,
					Type: p.Type,
					Desc: p.Desc,
				}
				merged.DelayedJSONDecode() // NOTE: delay decode for merged api
				return merged, true
			}
			return SwagConstValue{}, false // return false if don't merge
		}
	}

	return *exists, true
}

func verifyVal(p SwagConstValue) error {
	if !p.Type.IsPredefineable() {
		return fmt.Errorf("invalid type [%s] for predef const '%s'", p.Type, p.Name)
	}

	if p.Name == "" {
		return fmt.Errorf("missing name for predef const value '%s'", p.Name)
	}

	if p.Data.Empty() && !p.Type.IsNullable() {
		return fmt.Errorf("missing data for predef const value '%s'", p.Name)
	}

	switch p.In {
	case expr.ParaTypeBody, expr.ParaTypeHeader, expr.ParaTypePath, expr.ParaTypeQuery:
		//expr.ParaTypeFormData,
		break
	default:
		return fmt.Errorf("invalid in [%s] for predef const [%s]", p.In, p.Name)
	}
	return nil
}

func mergePredefValue(left, right []SwagConstValue) ([]SwagConstValue, error) {
	totalSize := len(left) + len(right)
	if limit := 80; totalSize >= limit { // size assert for liner search
		return nil, fmt.Errorf("sizeof consts overflow %d", limit)
	}
	var out = make([]SwagConstValue, 0, totalSize)

	// merge api & global
	for i := 0; i < len(left); i++ {
		p := &left[i]
		if found, ok := findVal(right, p, true); ok {
			if err := verifyVal(found); err != nil {
				return nil, err
			}
			out = append(out, found)
		}
	}

	// add global-only
	for i := 0; i < len(right); i++ {
		p := &right[i]
		if found, ok := findVal(left, p, false); ok {
			if err := verifyVal(found); err != nil {
				return nil, err
			}
			out = append(out, found)
		}
	}

	return out, nil
}
