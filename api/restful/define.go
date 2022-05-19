package restful

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"

	"github.com/quanxiang-cloud/polyapi/pkg/core/gate"
	"github.com/quanxiang-cloud/polyapi/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
)

// APIType enum
const (
	APIInner = gate.APIInner
	APIRead  = gate.APIRead
	APIWrite = gate.APIWrite
)

//------------------------------------------------------------------------------

type pingResp struct {
	Msg       string `json:"msg"`
	Timestamp string `json:"timestamp"`
}

// PingPong return a pong response for ping
func PingPong(c *gin.Context) {
	r := &pingResp{
		Msg:       "pong",
		Timestamp: time.Now().Format("2006-01-02T15:04:05MST"),
	}
	resp.Format(r, nil).Context(c)
}

//------------------------------------------------------------------------------

// getPathArg parse a path arg
func getPathArg(c *gin.Context, name string, out *string, err *error) bool {
	*out = c.Param(name)
	if *out == "" {
		*err = fmt.Errorf("missing path arg %s", name)
		return false
	}
	return true
}

var getRequestArgs = httputil.GetRequestArgs // get params by method

var bindBody = httputil.BindBody // func (c *gin.Context, d interface{}) error

func paraErr(err string) error2.Error {
	return error2.NewErrorWithString(error2.ErrParams, err)
}

func innerErr(err string) error2.Error {
	return error2.NewErrorWithString(error2.Internal, err)
}

// unsafeByteString convert []byte to string without copy
// the origin []byte **MUST NOT** accessed after that
func unsafeByteString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// unsafeStringBytes return GoString's buffer slice
// ** NEVER modify returned []byte **
func unsafeStringBytes(s string) []byte {
	var bh reflect.SliceHeader
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return *(*[]byte)(unsafe.Pointer(&bh))
}
