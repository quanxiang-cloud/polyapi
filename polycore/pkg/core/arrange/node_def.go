// functional arrange protocol
// start with one input node
// end with one output node

package arrange

import (
	"encoding/json"
	"errors"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/defines/consts"
	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/enumset"
	"github.com/quanxiang-cloud/polyapi/polycore/pkg/core/protocol"
)

// client defined data field name
const (
	//dispTypeName = "disp"

	polyRequestTempVar = consts.PolyRequestTempVar // temp var to call request
)

// Evaler export
type Evaler = protocol.Evaler

// // client define, this is just a sample...
// type DispDemo struct {
// 	X int `json:"x"`
// 	Y int `json:"y"`
// }

// // TypeName returns name of the type
// func (v *DispDemo) TypeName() expr.Enum { return dispTypeName }

//------------------------------------------------------------------------------

// Node represents an arrange node
type Node struct {
	Name      string              `json:"name"`      // node name, unique, English only
	Title     string              `json:"title"`     // node name alias, unique, allow Chinese
	Desc      string              `json:"desc"`      // description of this node
	Type      enumset.Enum        `json:"type"`      // input|request|output|if
	NextNodes []string            `json:"nextNodes"` // name of parent nodes
	Detail    expr.FlexJSONObject `json:"detail"`    // detail of this node
	//Disp      expr.FlexJsonObject `json:"disp"`      // display of this node, for client only

	branch branchRoutes // record branch routes of "if" nodes
}

func (v *Node) getNameAlias() string {
	if v.Title != "" {
		return v.Title
	}
	return v.Name
}

// DelayedJSONDecode delay unmarshal flex json object
func (v *Node) DelayedJSONDecode() error {
	if err := flexFactory.DelayedUnmarshalFlexJSONObject(v.Type.String(), &v.Detail); err != nil {
		return err
	}
	// if err := factory.DelayedUnmarshalFlexJsonObject(dispTypeName, &v.Disp); err != nil {
	// 	return err
	// }
	d := v.Detail.D.(expr.DelayedJSONDecoder)
	if err := d.DelayedJSONDecode(); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------

// Arrange represent an arrange config.
type Arrange struct {
	// for poly doc
	Info  *APIInfo `json:"info,omitempty"`
	Nodes []Node   `json:"nodes"` // nodes of this arrange
}

// TypeName returns name of the type
func (a Arrange) TypeName() string { return NodeTypeArrange.String() }

// Builder impliments script build function for nodes
type Builder interface {
	BuildJsScript(n *Node, b *jsBuilder, depth int) error
}

// InitArrange return an initialized arrange text
func InitArrange(template string, info *Arrange) (string, error) {
	if template != "" {
		return copyArrange(template, info)
	}

	var arrange = Arrange{
		Info: info.Info,
	}
	b, err := json.Marshal(arrange)
	if err != nil {
		return "", err
	}
	return unsafeByteString(b), nil
}

func copyArrange(template string, info *Arrange) (string, error) {
	if template == "" {
		return "", errors.New("missing arrange template when copy")
	}
	var arrange map[string]interface{}
	if err := json.Unmarshal(unsafeStringBytes(template), &arrange); err != nil {
		return "", err
	}

	arrange["info"] = info.Info
	b, err := json.Marshal(arrange)
	if err != nil {
		return "", err
	}
	return unsafeByteString(b), nil
}
