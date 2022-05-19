package adaptor

import (
	"context"
	"encoding/json"
	"net/http"
)

// KMSAuthorizeRespItem KMSAuthorizeRespItem
type KMSAuthorizeRespItem struct {
	In    string `json:"in"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// KMSAuthorizeResp KMSAuthorizeResp
type KMSAuthorizeResp struct {
	Token []*KMSAuthorizeRespItem `json:"token"`
}

// ListKMSCustomerKeyReq req
type ListKMSCustomerKeyReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"limit"`
	Service  string `json:"service"`
	Active   int    `json:"active"`
	Owner    string `json:"-"`
}

// KMSCustomerKey key info except secret
type KMSCustomerKey struct {
	ID          string `json:"id"`    //unique id
	Owner       string `json:"owner"` //owner id
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`     //belong service, eg: system_form
	Host        string `json:"host"`        //service host, eg: api.xxx.com:8080
	AuthType    string `json:"authType"`    //signature/cookie/oauth2...
	AuthContent string `json:"authContent"` //Authorize detail
	KeyID       string `json:"keyID"`       //key id
	KeyContent  string `json:"keyContent"`  //key content
	Active      int    `json:"active"`      //1 active 0 disable
	CreateAt    int64  `json:"createAt"`    //create time
	UpdateAt    int64  `json:"updateAt"`    //update time
}

// ListKMSCustomerKeyResp resp
type ListKMSCustomerKeyResp struct {
	Keys  []*KMSCustomerKey `json:"keys"`
	Total int               `json:"total"`
}

// QueryKMSCustomerKeyResp resp
type QueryKMSCustomerKeyResp struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	OwnerName   string `json:"ownerName"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Service     string `json:"service"`
	Host        string `json:"host"`
	AuthType    string `json:"authType"`
	AuthContent string `json:"authContent"`
	KeyID       string `json:"keyID"`
	KeyContent  string `json:"keyContent"`
	Active      int    `json:"active"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

// UpdateCustomerKeyInBatchReq UpdateCustomerKeyInBatchReq
type UpdateCustomerKeyInBatchReq struct {
	Host        string `json:"host"`
	Service     string `json:"service"`
	AuthType    string `json:"authType"`
	AuthContent string `json:"authContent"`
}

// UpdateCustomerKeyInBatchResp UpdateCustomerKeyInBatchResp
type UpdateCustomerKeyInBatchResp struct {
}

// DeleteCustomerKeyInBatchReq DeleteCustomerKeyInBatchReq
type DeleteCustomerKeyInBatchReq struct {
	Namespace   string   `json:"namespace"`
	ServiceName []string `json:"serviceName"`
}

// DeleteCustomerKeyInBatchResp DeleteCustomerKeyInBatchResp
type DeleteCustomerKeyInBatchResp struct {
}

// DeleteCustomerKeyByPrefixReq DeleteCustomerKeyByPrefixReq
type DeleteCustomerKeyByPrefixReq struct {
	NamespacePath string `json:"namespacePath"`
}

// DeleteCustomerKeyByPrefixResp DeleteCustomerKeyByPrefixResp
type DeleteCustomerKeyByPrefixResp struct {
}

// CheckAuthReq CheckAuthReq
type CheckAuthReq struct {
	AuthType    string `json:"authType"`
	ServicePath string `json:"servicePath"`
	AuthContent string `json:"authContent"`
}

// CheckAuthResp CheckAuthResp
type CheckAuthResp struct {
}

// KMSOper is the interface for kms proxy
type KMSOper interface {
	// Authorize request third party token by key
	Authorize(c context.Context, keyUUID string, body json.RawMessage, header http.Header) (*KMSAuthorizeResp, error)
	ListCustomerAPIKey(c context.Context, req *ListKMSCustomerKeyReq) (*ListKMSCustomerKeyResp, error)
	QueryCustomerAPIKey(c context.Context, keyUUID string) (*QueryKMSCustomerKeyResp, error)
	UpdateCustomerKeyInBatch(c context.Context, req *UpdateCustomerKeyInBatchReq) (*UpdateCustomerKeyInBatchResp, error)
	DeleteCustomerKeyInBatch(c context.Context, req *DeleteCustomerKeyInBatchReq) (*DeleteCustomerKeyInBatchResp, error)
	DeleteCustomerKeyByPrefix(c context.Context, req *DeleteCustomerKeyByPrefixReq) (*DeleteCustomerKeyByPrefixResp, error)
	CheckAuth(c context.Context, req *CheckAuthReq) (*CheckAuthResp, error)
}

// SetKMSOper set the instance of kms oper
func SetKMSOper(f KMSOper) KMSOper {
	i := getInst()
	old := i.kmsOper
	i.kmsOper = f
	return old
}

// GetKMSOper get the instance of kms oper
func GetKMSOper() KMSOper {
	i := getInst()
	return i.kmsOper
}
