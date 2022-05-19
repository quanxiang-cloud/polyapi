package arrange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/errcode"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/apipath"
)

// APIInfo is info of api
type APIInfo struct {
	Schema    string `json:"schema,omitempty"`
	Host      string `json:"host,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"` // name of this arrange
	Title     string `json:"title,omitempty"`
	Desc      string `json:"desc,omitempty"`    // description of this arrange
	Version   string `json:"version,omitempty"` // version of this arrange
	Method    string `json:"method,omitempty"`
	Encoding  string `json:"encoding,omitempty"` // json only current
}

// APIPath get api full path
func (p *APIInfo) APIPath() string {
	return apipath.Join(p.Namespace, p.Name)
}

// BuildJsScript generate js code from arrange config.
// Which generate a temp one-time function and execute it immediately.
func BuildJsScript(info *APIInfo, arrangeJSON string, owner string) (string, string, []string, error) {
	a := &Arrange{}
	if err := json.Unmarshal(unsafeStringBytes(arrangeJSON), a); err != nil {
		return "", "", nil, err
	}
	a.Info = info
	if a.Info == nil {
		a.Info = &APIInfo{}
	}

	b := newBuilder(owner)
	if b.evaler != nil {
		defer b.evaler.Free() //NOTE: ensure free b.evaler to avoid leak
	}

	script, err := b.build(a)
	if err != nil {
		return "", "", nil, err
	}

	// auto generate swagger for this poly api
	d, err := PolyDocGennerator(
		a,
		b.startNode.Detail.D.(*InputNodeDetail),
		b.endNode.Detail.D.(*OutputNodeDetail),
	)
	if err != nil {
		return "", "", nil, err
	}

	docBytes, err := json.Marshal(d)
	if err != nil {
		return "", "", nil, err
	}

	doc := unsafeByteString(docBytes)

	return script, doc, b.getRawAPIList(), err
}

func newBuilder(owner string) *jsBuilder {
	b := make([]byte, 0, 1024)
	var e Evaler
	if oper := adaptor.GetEvalerOper(); oper != nil {
		e = oper.CreateEvaler()
	}
	return &jsBuilder{
		owner:        owner,
		buf:          bytes.NewBuffer(b),
		q:            newNodeDeque(-1),
		visited:      make(map[string]*Node),
		visitedAlias: make(map[string]*Node),
		rawList:      make(map[string]struct{}),
		config:       nil,
		startNode:    nil,
		endNode:      nil,
		evaler:       e,
	}
}

// javascript builder
type jsBuilder struct {
	owner        string
	buf          *bytes.Buffer
	q            *nodeDeque
	config       *Arrange
	visited      map[string]*Node
	visitedAlias map[string]*Node
	startNode    *Node
	endNode      *Node
	rawList      map[string]struct{}
	evaler       Evaler
}

func (b *jsBuilder) putRawAPI(apiPath string) {
	b.rawList[apiPath] = struct{}{}
}

func (b *jsBuilder) getRawAPIList() []string {
	ret := make([]string, 0, len(b.rawList))
	for k := range b.rawList {
		ret = append(ret, k)
	}
	return ret
}

// write new line
func (b *jsBuilder) writeLn() {
	b.buf.WriteByte('\n')
}

func (b *jsBuilder) build(arrange *Arrange) (string, error) {
	buf := b.buf
	comment := fmt.Sprintf("// polyTmpScript_%s_%s\n",
		arrange.Info.APIPath(), time.Now().Format("2006-01-02T15:04:05MST"))
	buf.WriteString(comment)                    // script name comment
	buf.WriteString("var _tmp = function(){\n") // tmp function begin
	b.config = arrange

	const d = consts.PolyAllDataVarName
	const input = consts.PolyAPIInputVarName
	buf.WriteString(fmt.Sprintf("  var %s = { \"%s\": %s, } // qyAllLocalData\n\n",
		d, input, input))

	if err := b.buildNodes(); err != nil {
		return "", err
	}
	buf.WriteString("}; _tmp();\n") // tmp function end

	return buf.String(), nil
}

func (b *jsBuilder) buildNodes() error {
	if b.config == nil || len(b.config.Nodes) < 2 {
		return errcode.ErrBuildNodeCnt.NewError()
	}

	start := &b.config.Nodes[0]
	start.branch.Clean()
	b.q.PushBack(start) // push the first node
	for !b.q.Empty() {
		n, _ := b.q.PopFront()

		if err := b.checkVisitNode(n); err != nil {
			return err
		}

		if n != b.endNode {
			if err := b.buildNode(n, 1); err != nil {
				return err
			}
		}
	}

	if err := b.checkAfterVisitAllNodes(); err != nil {
		return err
	}
	// build output code
	if err := b.buildNode(b.endNode, 1); err != nil {
		return err
	}

	return nil
}

// check when visit all nodes
func (b *jsBuilder) checkAfterVisitAllNodes() error {
	// check if mission end node
	if b.endNode == nil {
		return errcode.ErrBuildMissOutput.FmtError(NodeNameEnd, NodeTypeOutput)
	}
	for i := 0; i < len(b.config.Nodes); i++ {
		n := &b.config.Nodes[i]
		if _, ok := b.visited[n.Name]; !ok {
			return errcode.ErrBuildIsolateNode.FmtError(n.getNameAlias(), n.Type)
		}
	}
	return nil
}

// check rules of node
func (b *jsBuilder) checkVisitNode(n *Node) error {
	if n.Name == "" {
		return errcode.ErrBuildNoNameNode.FmtError(n.Type)
	}

	// record first node
	if b.startNode == nil {
		if n.Type != NodeTypeInput || !NodeNameStart.Equal(n.Name) {
			return errcode.ErrBuildMissStart.FmtError(NodeNameStart, NodeTypeInput)
		}
		b.startNode = n
	}

	// have been visited
	if old, ok := b.visited[n.Name]; ok && n.Type != NodeTypeOutput {
		if n == old { // duplicate visit the same node
			return errcode.ErrBuildDuplicateNode.FmtError(old.getNameAlias(), old.Type)
		}

		// duplicate name of two node
		return errcode.ErrBuildDuplicateNodeName.FmtError(n.getNameAlias(), n.Type, old.getNameAlias(), old.Type)
	}

	aliasName := n.getNameAlias()
	if old, ok := b.visitedAlias[aliasName]; ok && n.Type != NodeTypeOutput {
		if n == old { // duplicate visit the same node
			return errcode.ErrBuildDuplicateNode.FmtError(old.getNameAlias(), old.Type)
		}

		// duplicate name of two node
		return errcode.ErrBuildDuplicateNodeAlias.FmtError(aliasName, n.Type, old.getNameAlias(), old.Type)
	}

	// check node config
	switch n.Type {
	case NodeTypeInput:
		if n != b.startNode {
			return errcode.ErrBuildMultiInput.FmtError(b.startNode.getNameAlias(), n.getNameAlias())
		}
	case NodeTypeOutput:
		if b.endNode == nil {
			if !NodeNameEnd.Equal(n.Name) {
				return errcode.ErrBuildOutputName.FmtError(NodeNameEnd, n.getNameAlias())
			}
			if len(n.NextNodes) > 0 {
				return errcode.ErrBuildOutputWithNext.FmtError(NodeNameEnd, n.NextNodes)
			}
			b.endNode = n
		}
		if n != b.endNode {
			return errcode.ErrBuildMultiOutput.FmtError(b.endNode.getNameAlias(), n.getNameAlias())
		}
	case NodeTypeIf:
		if len(n.NextNodes) > 0 {
			n.NextNodes = nil // NOTE: ignore next node of if node
			// NOTE: don't check next node for "if" node
			//return errcode.ErrBuildIfWithNext.FmtError(n.getNameAlias(), n.NextNodes)
		}
	case NodeTypeRequest:
		if len(n.NextNodes) == 0 {
			return errcode.ErrBuildRequestWithoutNext.FmtError(n.getNameAlias())
		}
		// if !EncodingEnum.Verify(n.Encoding) {
		// 	return fmt.Errorf(`"invalid encoding [%s] for node [%s(%s)]: %v`,
		// 		n.Encoding, n.Name, n.Type, EncodingEnum.GetAll())
		// }
	default:
		return errcode.ErrBuildUnknownNodeType.FmtError(n.Type, n.getNameAlias())
	}

	b.visited[n.Name] = n
	b.visitedAlias[aliasName] = n
	// push next nodes
	for _, v := range n.NextNodes {
		nn := b.findNode(v)
		if nn == nil {
			return errcode.ErrBuildUnknownNextNode.FmtError(n.getNameAlias(), n.Type, v)
		}
		nn.branch.MergeFrom(n.branch)
		b.q.PushBack(nn)
	}
	return nil
}

// build the node by config
// n is verified node
func (b *jsBuilder) buildNode(n *Node, fromIfDepth int) (err error) {
	defer func() {
		if err != nil {
			err = errcode.ErrBuildNodeFail.FmtError(n.getNameAlias(), err.Error())
		}
	}()

	if err := ValidateNodeName(n.Name); err != nil { // check node name rule
		return err
	}

	// delay decoding JSON for node detail
	if err := n.DelayedJSONDecode(); err != nil {
		return err
	}

	d, ok := n.Detail.D.(Builder)
	if !ok {
		return fmt.Errorf(`node %s(%s) detail doesn't implement Builder: %#v`,
			n.Name, n.Type, n.Detail.D)
	}
	return d.BuildJsScript(n, b, fromIfDepth)
}

// find node by name
func (b *jsBuilder) findNode(name string) *Node {
	if n := b.findVisitedNode(name); n != nil {
		return n
	}

	for i := 0; i < len(b.config.Nodes); i++ {
		n := &b.config.Nodes[i]
		if n.Name == name {
			return n
		}
	}

	return nil
}

// find visited node by name
func (b *jsBuilder) findVisitedNode(name string) *Node {
	if n, ok := b.visited[name]; ok {
		return n
	}
	return nil
}
