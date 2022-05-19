package apiprovider

import (
	"encoding/json"
	"net/http"
)

// QueryDocReq is request of api doc
type QueryDocReq struct {
	APIType    string `json:"-" binding:"-"`
	APIPath    string `json:"-" binding:"-"`
	DocType    string `json:"docType" binding:"max=64"`
	TitleFirst bool   `json:"titleFirst"`
}

// QueryDocResp is response of api doc
type QueryDocResp struct {
	Title   string          `json:"title"`
	APIPath string          `json:"apiPath"`
	DocType string          `json:"docType"`
	Doc     json.RawMessage `json:"doc"`
}

// RequestReq is api request arg
type RequestReq struct {
	Owner          string          `json:"-"`
	APIService     string          `json:"-" binding:"-"` // NOTE: optional, use specified service
	APIServiceArgs string          `json:"-" binding:"-"` // 	   optional, use specified service args
	APIPath        string          `json:"-" binding:"-"`
	APIType        string          `json:"-" binding:"-"`
	Method         string          `json:"-" binding:"-"`
	Body           json.RawMessage `json:"-"` // input json string
	Header         http.Header     `json:"-"` // input header
}

// RequestResp is api request response
type RequestResp struct {
	APIPath    string          `json:"apiPath"`  //
	Response   json.RawMessage `json:"response"` // response json string
	Header     http.Header     `json:"header"`   // response header(cookies)
	StatusCode int             `jons:"statusCode"`
	Status     string          `json:"status"`
}
