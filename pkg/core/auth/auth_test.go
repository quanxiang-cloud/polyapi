package auth

import (
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/adaptor"
)

func TestAppendBody(t *testing.T) {
	add := adaptor.KMSAuthorizeResp{
		Token: []*adaptor.KMSAuthorizeRespItem{
			&adaptor.KMSAuthorizeRespItem{
				Name:  "x",
				Value: "xx",
			},
			&adaptor.KMSAuthorizeRespItem{
				Name:  "a.b.c",
				Value: "a.b.c++",
			},
		},
	}
	b, err := appendToBody(`{}`, add.Token)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
