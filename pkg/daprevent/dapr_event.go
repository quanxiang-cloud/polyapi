package daprevent

import (
	"net/http"
)

/*
type DaprEvent struct {
	Topic           string      `json:"topic"`
	Pubsubname      string      `json:"pubsubname"`
	Traceid         string      `json:"traceid"`
	ID              string      `json:"id"`
	Datacontenttype string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
	Type            string      `json:"type"`
	Specversion     string      `json:"specversion"`
	Source          string      `json:"source"`
}
*/

// ReqDataInitAppPath ReqDataInitAppPath
type ReqDataInitAppPath struct {
	APPID     string      `json:"appID" binding:"required,max=64"`
	Owner     string      `json:"-" binding:"-"`
	OwnerName string      `json:"-" binding:"-"`
	Header    http.Header `json:"-" binding:"-"`
}

// EventInitAppPath satisfy DaprEvent
type EventInitAppPath struct {
	Data ReqDataInitAppPath `json:"data"`
}
