package xsvc

import (
	"encoding/base64"
	"encoding/json"
)

// XServiceArgs is the structure defination of X-Poly-Service-Args in header
type XServiceArgs struct {
	Schema      string `json:"schema"`
	Host        string `json:"host"`
	AuthType    string `json:"authType"`
	AuthContent string `json:"authContent"`
	KeyID       string `json:"keyID"`
	KeySecret   string `json:"keySecret"`
}

// Unmarshal serialize X-Poly-Service-Args
func Unmarshal(s string) (*XServiceArgs, error) {
	b, err := decode(s)
	if err != nil {
		return nil, err
	}

	xs := &XServiceArgs{}
	if err := json.Unmarshal(b, xs); err != nil {
		return nil, err
	}
	return xs, nil
}

func decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
