// Package enumset is enumset define provider
//
// NOTE: must call FinishReg in init()
/*
func init() {
	enumset.FinishReg()
}
*/
package enumset

import (
	"bytes"
	"fmt"
	"sort"
)

// Enum represents a enum
type Enum string

// String converts an enum to string
func (e Enum) String() string {
	return string(e)
}

// Equal judege if enum equals to string
func (e Enum) Equal(s string) bool {
	return string(e) == s
}

//------------------------------------------------------------------------------

// New create a new EnumSet
func New(exists *EnumSet) *EnumSet {
	ret := &EnumSet{}
	if exists != nil {
		ret.list = append(ret.list, exists.list...)
		if exists.contentMap != nil {
			ret.contentMap = make(map[string]*enumExContent)
			for k, v := range exists.contentMap {
				ret.contentMap[k] = v
			}
		}
	}
	allEnumSet = append(allEnumSet, ret)
	return ret
}

type enumExContent struct {
	op    string // eg "="
	title string // eg "等于"
}

// EnumSet represents a set of enum
type EnumSet struct {
	list       []string                  // enum set
	sorted     []string                  // ordered list
	contentMap map[string]*enumExContent // "eq" -> {"=","等于"}
}

// Reg regist a new enum to the set
func (es *EnumSet) Reg(val string) (Enum, error) {
	e := Enum(val)
	for _, v := range es.list {
		if e.Equal(v) {
			return "", fmt.Errorf("dupicate enum value %s", val)
		}
	}

	es.list = append(es.list, e.String())
	return e, nil
}

// MustReg regist a new enum to the set.
// It panic if val duplicate.
func (es *EnumSet) MustReg(val string) Enum {
	e, err := es.Reg(val)
	if err != nil {
		panic(err)
	}
	return e
}

// RegWithContent regist a new enum with content to the set. eg ("eq", "=")
func (es *EnumSet) RegWithContent(val string, op, title string) (Enum, error) {
	e := Enum(val)
	for _, v := range es.list {
		if e.Equal(v) {
			return "", fmt.Errorf("dupicate enum value %s", val)
		}
	}

	es.list = append(es.list, e.String())
	if es.contentMap == nil {
		es.contentMap = make(map[string]*enumExContent)
	}
	es.contentMap[val] = &enumExContent{
		op:    op,
		title: title,
	}
	return e, nil
}

// MustRegWithContent regist a new enum with content to the set. eg ("eq", "=").
// It panic if val duplicate.
func (es *EnumSet) MustRegWithContent(val string, op, title string) Enum {
	e, err := es.RegWithContent(val, op, title)
	if err != nil {
		panic(err)
	}
	return e
}

// ShowAll show enum list of the set
func (es *EnumSet) ShowAll() string {
	var b bytes.Buffer
	for i := 0; i < len(es.list); i++ {
		v := es.list[i]
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(v)
	}
	return b.String()
}

// ShowSorted show enum list of the set with sort
func (es *EnumSet) ShowSorted() string {
	var b bytes.Buffer
	for i := 0; i < len(es.sorted); i++ {
		v := es.sorted[i]
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(v)
	}
	return b.String()
}

// GetAll return enum list of the set
func (es *EnumSet) GetAll() []string {
	return es.list
}

// Sort sort the enums
func (es *EnumSet) Sort() {
	if len(es.sorted) > 0 {
		return
	}
	es.sorted = append([]string{}, es.list...)
	sort.Strings(es.sorted)
}

// Verify check if a enum is valid, binary search
func (es *EnumSet) Verify(e string) bool {
	s := es.sorted
	low, high := 0, len(s)-1
	for low <= high {
		mid := (low + high) / 2
		switch {
		case e == s[mid]:
			return true
		case e > s[mid]:
			low = mid + 1
		case e < s[mid]:
			high = mid - 1
		}
	}

	return false
}

// Content get content of the enum
func (es *EnumSet) Content(e string) (string, string, bool) {
	if es.contentMap != nil {
		if p, ok := es.contentMap[e]; ok {
			return p.op, p.title, true
		}
	}
	return "", "", false
}

// FindContent find content from multi enumset
func FindContent(e string, sets ...*EnumSet) (string, string, bool) {
	for _, v := range sets {
		if op, title, ok := v.Content(e); ok {
			return op, title, ok
		}
	}
	return "", "", false
}

//------------------------------------------------------------------------------

var allEnumSet []*EnumSet

// FinishReg sort the enums.
// NOTE: must call this at init() in customer package.
func FinishReg() {
	for _, v := range allEnumSet {
		v.Sort()
	}
}
