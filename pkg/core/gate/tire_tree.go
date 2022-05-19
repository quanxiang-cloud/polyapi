package gate

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	black byte = 'b'
	white byte = 'w'
	empty byte = 0
	any        = "*"
)

// NewTireTree create a TireTree object
func NewTireTree() *TireTree {
	p := &TireTree{}
	return p.Init()
}

// TireTree impliments tire key for ip config
type TireTree struct {
	root tireTreeNode
}

// Init init the tree
func (t *TireTree) Init() *TireTree {
	t.root.init("", empty)
	return t
}

// BatchInsert insert batch keys
func (t *TireTree) BatchInsert(keys []string, data byte) bool {
	suc := 0
	for _, v := range keys {
		if t.Insert(v, data) {
			suc++
		}
	}
	return suc > 0
}

func (t *TireTree) getData(keys []string, index int, data byte) byte {
	if index == len(keys)-1 {
		return data
	}
	return empty
}

func (t *TireTree) splitKey(key string) []string {
	subs := fastSplit(key, '.', 4) // eg: 127.0.0.1
	// for _, v := range subs[:len(subs)-1] { //allow any at last part only
	// 	if v == any {
	// 		return nil
	// 	}
	// }
	return subs
}

// Insert insert a single key
func (t *TireTree) Insert(key string, data byte) bool {
	subs := t.splitKey(key)
	n := &t.root
	for i, v := range subs {
		d := t.getData(subs, i, data)
		n, _ = n.addChild(v, d)
	}
	return len(subs) > 0
}

// Delete delete a single key
func (t *TireTree) Delete(key string) bool {
	subs := t.splitKey(key)
	nodes := t.root.getNodePath(subs)
	if nodes == nil {
		return false
	}

	d := black
	for i := len(nodes) - 1; i > 0; i-- {
		n, p := nodes[i], nodes[i-1]
		p.removeChild(n, d)
		d = empty
	}
	return true
}

// Match verify if a key matches config keys
func (t *TireTree) Match(key string) (byte, bool) {
	if n := t.root.getLeafNode(t.splitKey(key), true); n != nil {
		return n.data, true
	}
	return 0, false
}

// Show shows the config tree as string
func (t *TireTree) Show() string {
	buf := bytes.NewBuffer(nil)
	type nr struct {
		n     *tireTreeNode
		depth int
	}
	nodes := []nr{nr{&t.root, 1}}
	depth := -1
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		if n.depth > depth {
			if depth > 0 {
				buf.WriteByte('\n')
			}
			depth = n.depth
			buf.WriteString(fmt.Sprintf("%d: ", n.depth))
		}
		for _, v := range n.n.children {
			if v != nil {
				if !v.isEmpty() {
					nodes = append(nodes, nr{v, n.depth + 1})
				}
				buf.WriteString(fmt.Sprintf("%s ", v.show()))
			} else {
				buf.WriteString("# ")
			}
		}
		buf.WriteString(" | ")
	}
	return buf.String()
}

//------------------------------------------------------------------------------

type tireTreeNode struct {
	data     byte
	key      string
	children []*tireTreeNode // first child is for key "*"
}

func newNode(subKey string, data byte) *tireTreeNode {
	p := &tireTreeNode{}
	return p.init(subKey, data)
}

func (t *tireTreeNode) show() string {
	if t.data == empty {
		return t.key
	}
	return fmt.Sprintf("%s(%s)", t.key, string(t.data))
}

func (t *tireTreeNode) getNodePath(subKeys []string) []*tireTreeNode {
	r := make([]*tireTreeNode, 0, len(subKeys)+1)
	n := t
	r = append(r, n)
	for _, v := range subKeys {
		if n, _ = n.findChild(v, false); n == nil {
			return nil
		}
		r = append(r, n)
	}
	return r
}

func (t *tireTreeNode) getLeafNode(subKeys []string, allowAny bool) *tireTreeNode {
	return t._getLeafNode(subKeys, allowAny, 0)
}

func (t *tireTreeNode) _getLeafNode(subKeys []string, allowAny bool, depth int) *tireTreeNode {
	n := t
	index := 0
	anyCount := 0
	for _, v := range subKeys {
		if n, index = n.findChild(v, allowAny); n == nil {
			if allowAny && anyCount > 0 && depth == 0 {
				// eg: [127.*.1,127.0.0] matches 127.0.0
				// it will match fail using *
				// but it matches ok when ignore *
				return t._getLeafNode(subKeys, false, depth+1) // try ignore key *
			}
			return nil
		}
		if index == 0 {
			anyCount++
		}
	}
	return n
}

func (t *tireTreeNode) init(subKey string, data byte) *tireTreeNode {
	t.key = subKey
	t.data = data
	t.children = []*tireTreeNode{nil}
	return t
}

func (t *tireTreeNode) sort() {
	s, _ := t.withoutAny()
	sort.Slice(s, func(i, j int) bool { return s[i].key < s[j].key })
}

func (t *tireTreeNode) withoutAny() ([]*tireTreeNode, int) {
	from := 1
	return t.children[from:], from
}

func (t *tireTreeNode) isEmpty() bool {
	return len(t.children) == 0 || len(t.children) == 1 && t.children[0] == nil
}

func (t *tireTreeNode) isLeaf() bool {
	return t.data != empty
}

func (t *tireTreeNode) hasAny() bool {
	_, index := t.findAny()
	return index >= 0
}

func (t *tireTreeNode) findAny() (*tireTreeNode, int) {
	if t.children[0] != nil {
		return t.children[0], 0
	}
	return nil, -1
}

func (t *tireTreeNode) addAny(data byte) (*tireTreeNode, int) {
	t.children[0] = newNode(any, data)
	return t.findAny()
}

func (t *tireTreeNode) addChild(subKey string, data byte) (*tireTreeNode, int) {
	if n, i := t.findChild(subKey, false); i >= 0 {
		if n.data == empty && data != empty { //update leaf data
			n.data = data
		}
		return n, i
	}
	if subKey == any {
		return t.addAny(data)
	}
	t.children = append(t.children, newNode(subKey, data))
	t.sort()
	return t.findChild(subKey, false)
}

func (t *tireTreeNode) find(s []*tireTreeNode, key string) int {
	low, high := 0, len(s)-1
	for low <= high {
		mid := (low + high) / 2
		switch {
		case key == s[mid].key:
			return mid
		case key < s[mid].key:
			high = mid - 1
		case key > s[mid].key:
			low = mid + 1
		}
	}
	return -1
}

func (t *tireTreeNode) findChild(subKey string, allowAny bool) (*tireTreeNode, int) {
	if isAny := subKey == any; isAny || allowAny {
		if n, index := t.findAny(); index >= 0 || isAny {
			return n, index
		}
	}

	s, from := t.withoutAny()
	if idx := t.find(s, subKey); idx >= 0 {
		return s[idx], idx + from
	}

	return nil, -1
}

func (t *tireTreeNode) _removeChild(index int) {
	if index == 0 {
		t.children[0] = nil
		return
	}
	if index > 0 && index < len(t.children) {
		t.children = append(t.children[:index], t.children[index+1:]...)
	}
}

func (t *tireTreeNode) removeChild(child *tireTreeNode, data byte) bool {
	if data != empty {
		child.data = empty
	}
	if child.isEmpty() {
		_, index := t.findChild(child.key, false)
		t._removeChild(index)
		return true
	}
	return false
}

//------------------------------------------------------------------------------

func fastSplit(s string, sep byte, n int) []string {
	t := s
	r := make([]string, 0, n)
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			r = append(r, t[:i-(len(s)-len(t))])
			t = s[i+1:]
		}
	}
	r = append(r, t)
	return r
}
