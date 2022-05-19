package adaptor

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/apiprovider"

	"github.com/quanxiang-cloud/polyapi/pkg/core/expr"
)

// RawAPIOper is the interface for query raw api
type RawAPIOper interface {
	Query(c context.Context, req *QueryRawAPIReq) (*QueryRawAPIResp, error)
	QueryInBatches(c context.Context, req *QueryRawAPIInBatchesReq) (*QueryRawAPIInBatchesResp, error)
	List(c context.Context, req *RawListReq) (*RawListResp, error)
	ListInService(c context.Context, req *ListInServiceReq) (*ListInServiceResp, error)
	InnerUpdateRawInBatch(ctx context.Context, req *InnerUpdateRawInBatchReq) (*InnerUpdateRawInBatchResp, error)
	InnerDel(c context.Context, req *DelReq) (*DelResp, error)
	ValidInBatches(ctx context.Context, req *RawValidInBatchesReq) (*RawValidInBatchesResp, error)
	ValidByPrefixPath(ctx context.Context, req *RawValidByPrefixPathReq) (*RawValidByPrefixPathResp, error)
	ListByPrefixPath(ctx context.Context, req *ListRawByPrefixPathReq) (*ListRawByPrefixPathResp, error)
	InnerImport(ctx context.Context, req *InnerImportRawReq) (*InnerImportRawResp, error)
	InnerDelByPrefixPath(ctx context.Context, req *InnerDelRawByPrefixPathReq) (*InnerDelRawByPrefixPathResp, error)

	QuerySwagger(ctx context.Context, req *QueryRawSwaggerReq) (*QueryRawSwaggerResp, error)

	QueryDoc(c context.Context, req *apiprovider.QueryDocReq) (*apiprovider.QueryDocResp, error)
}

// QueryRawSwaggerReq QueryRawSwaggerReq
type QueryRawSwaggerReq struct {
	APIPath []string `json:"-"`
}

// QueryRawSwaggerResp QueryRawSwaggerResp
type QueryRawSwaggerResp struct {
	Swagger []byte `json:"swagger"`
}

// InnerDelRawByPrefixPathReq InnerDelRawByPrefixPathReq
type InnerDelRawByPrefixPathReq struct {
	NamespacePath string `json:"-"`
}

// InnerDelRawByPrefixPathResp InnerDelRawByPrefixPathResp
type InnerDelRawByPrefixPathResp struct {
}

// InnerImportRawReq InnerImportRawReq
type InnerImportRawReq struct {
	List []*RawAPIFull `json:"list"`
}

// InnerImportRawResp InnerImportRawResp
type InnerImportRawResp struct {
}

// ListRawByPrefixPathReq ListRawByPrefixPathReq
type ListRawByPrefixPathReq = RawListReq

// ListRawByPrefixPathResp ListRawByPrefixPathResp
type ListRawByPrefixPathResp struct {
	List  []*RawAPIFull `json:"list"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
}

// RawValidByPrefixPathReq RawValidByPrefixPathReq
type RawValidByPrefixPathReq struct {
	NamespacePath string `json:"-"`
	Valid         uint   `json:"valid"`
}

// RawValidByPrefixPathResp RawValidByPrefixPathResp
type RawValidByPrefixPathResp struct {
}

// RawValidInBatchesReq RawValidInBatchesReq
type RawValidInBatchesReq struct {
	APIPath []string
	Valid   uint `json:"valid"`
}

// RawValidInBatchesResp RawValidInBatchesResp
type RawValidInBatchesResp struct{}

// DelReq DelReq
type DelReq struct {
	NamespacePath string   `json:"-" binding:"-"` //
	Names         []string `json:"names"`
	Owner         string   `json:"-"`
}

// DelResp DelResp
type DelResp struct {
}

// InnerUpdateRawInBatchReq InnerUpdateRawInBatchReq
type InnerUpdateRawInBatchReq struct {
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	Host      string `json:"host"`
	Schema    string `json:"schema"`
	AuthType  string `json:"authType"`
}

// InnerUpdateRawInBatchResp InnerUpdateRawInBatchResp
type InnerUpdateRawInBatchResp struct {
}

// ListInServiceReq ListInServiceReq
type ListInServiceReq struct {
	ServicePath string `json:"-"`
	Active      int    `json:"active"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

// ListInServiceResp ListInServiceResp
type ListInServiceResp struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	List  []*RawListNode `json:"list"`
}

// RawListReq RawListReq
type RawListReq struct {
	NamespacePath string `uri:"namespacePath"`
	Active        int    `json:"active"`
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
}

// RawListNode RawListNode
type RawListNode struct {
	ID         string `json:"id"`
	Owner      string `json:"owner"`
	OwnerName  string `json:"ownerName"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	FullPath   string `json:"fullPath"`
	URL        string `json:"url"`
	Version    string `json:"version"`
	Method     string `json:"method"`
	Action     string `json:"action"`
	Active     uint   `json:"active"`
	Valid      uint   `json:"valid"`
	CreateAt   int64  `json:"createAt"`
	UpdateAt   int64  `json:"updateAt"`
	URI        string `json:"uri"`
	AccessPath string `json:"accessPath"`
}

// RawListResp RawListResp
type RawListResp struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	List  []*RawListNode `json:"list"`
}

// QueryRawAPIInBatchesReq QueryRawAPIInBatchesReq
type QueryRawAPIInBatchesReq struct {
	APIPathList []string
}

// QueryRawAPIInBatchesResp QueryRawAPIInBatchesResp
type QueryRawAPIInBatchesResp struct {
	List []*RawListNode `json:"list"`
}

// QueryRawAPIReq QueryRawAPIReq
type QueryRawAPIReq struct {
	APIPath string `json:"-"`
}

// QueryRawAPIResp QueryRawAPIResp
type QueryRawAPIResp struct {
	Content   *RawAPIContent `json:"-"`
	ID        string         `json:"id"`
	URL       string         `json:"url"`
	Method    string         `json:"method"`
	Owner     string         `json:"owner"`
	OwnerName string         `json:"ownerName"`
	Namespace string         `json:"namespace"`
	Name      string         `json:"name"`
	Title     string         `json:"title"`
	Desc      string         `json:"desc"`
	Active    uint           `json:"active"`
	Valid     uint           `json:"valid"`
	Service   string         `json:"service"`
	Schema    string         `json:"schema"`
	Host      string         `json:"host"`
	AuthType  string         `json:"authType"`
	UpdateAt  int64          `json:"updateAt"`
}

// RawAPIContent is the content field detail
type RawAPIContent struct {
	ID          string                 `json:"x-id"`
	Action      string                 `json:"x-action"` // action, extended
	Consts      expr.SwagConstValueSet `json:"x-consts"` // predefined values, extended
	Input       expr.InputNodeDetail   `json:"x-input"`
	Output      expr.OutputNodeDetail  `json:"x-output"`
	BasePath    string                 `json:"basePath"`
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	EncodingIn  string                 `json:"encoding-in"`
	EncodingOut string                 `json:"encoding-out"`
	Summary     string                 `json:"summary"`
	Desc        string                 `json:"desc"`
}

// RawAPIFull is the raw api db scheme
type RawAPIFull struct {
	ID        string
	Owner     string
	OwnerName string
	Namespace string
	Service   string
	Name      string
	Title     string
	Desc      string
	Path      string
	URL       string
	Action    string // action
	Method    string // method GET|POST|...
	Version   string
	Access    uint
	Active    uint
	Valid     uint

	Schema   string
	Host     string
	AuthType string

	CreateAt int64  // create time
	UpdateAt int64  // update time
	DeleteAt *int64 // delete time

	Doc     *APIDoc
	Content *RawAPIContent
}

// GetAction return the action predef value
func (c *RawAPIContent) GetAction() *expr.SwagConstValue {
	return c.Consts.GetAction()
}

// Check verify the inputs
func (c *RawAPIContent) Check() error {
	// if c.URL == "" {
	// 	return errors.New("missing path")
	// }
	// if p.Service == "" {
	// 	return errors.New("missing service")
	// }
	// if c.Schema == "" {
	// 	return errors.New("missing schema")
	// }
	if c.Method == "" {
		return errors.New("missing method")
	}
	if c.EncodingIn == "" {
		return errors.New("missing encoding in")
	}
	if c.EncodingOut == "" {
		return errors.New("missing encoding out")
	}

	expect := (c.Action != "")
	got := (c.GetAction() != nil)
	if expect != got {
		if expect {
			return errors.New("missing action in const value")
		}
		return errors.New("unnecessary action in const value")
	}

	return nil
}

// Value marshal
func (c RawAPIContent) Value() (driver.Value, error) {
	//return json.MarshalIndent(c, "", "  ")
	return json.Marshal(c)
}

// Scan unmarshal
func (c *RawAPIContent) Scan(data interface{}) error {
	if err := json.Unmarshal(data.([]byte), c); err != nil {
		return err
	}
	if err := c.Consts.DelayedJSONDecode(); err != nil {
		return err
	}

	return nil
}

// APIDoc is the api doc schema
type APIDoc struct {
	ID       string           `json:"x-id"`
	Version  string           `json:"version"`
	FmtInOut expr.FmtAPIInOut `json:"x-fmt-inout"` // formated input and output
	Swagger  json.RawMessage  `json:"x-swagger"`

	Defines             json.RawMessage `json:"defines,omitempty"`             // swag format
	SecurityDefinitions json.RawMessage `json:"securityDefinitions,omitempty"` // swag format
	Security            json.RawMessage `json:"security,omitempty"`            // swag format
	Parameters          json.RawMessage `json:"parameters,omitempty"`          // swag format
	Responses           json.RawMessage `json:"responses,omitempty"`           // swag format
}

// Value marshal
func (c APIDoc) Value() (driver.Value, error) {
	//return json.MarshalIndent(c, "", "  ")
	return json.Marshal(c)
}

// Scan unmarshal
func (c *APIDoc) Scan(data interface{}) error {
	d, ok := data.([]byte)
	if !ok {
		return errors.New("unknown scan data type")
	}
	if err := json.Unmarshal(d, c); err != nil {
		return err
	}
	// NOTE: dont auto DelayedJSONDecode here, call it when require
	// if err := c.FmtInOut.DelayedJSONDecode(); err != nil {
	// 	return err
	// }
	return nil
}

// SetRawAPIOper set the instance of RawAPIOper
func SetRawAPIOper(f RawAPIOper) RawAPIOper {
	i := getInst()
	old := i.rawAPIOper
	i.rawAPIOper = f
	return old
}

// GetRawAPIOper get the instance of RawAPIOper
func GetRawAPIOper() RawAPIOper {
	i := getInst()
	return i.rawAPIOper
}
