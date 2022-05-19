package gate

import (
	"net/http"
	"strings"

	"github.com/quanxiang-cloud/polyapi/pkg/basic/polysign"

	"github.com/gin-gonic/gin"
)

var allowHeaders = getAllowHeaders()

// Cors handle CORS request
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func getAllowHeaders() string {
	var buf strings.Builder
	buf.WriteString("Content-Type,Access-Token,X-CSRF-Token,Authorization,Token")
	buf.WriteByte(',')
	buf.WriteString(polysign.XHeaderPolySignKeyID)
	buf.WriteByte(',')
	buf.WriteString(polysign.XHeaderPolySignMethod)
	buf.WriteByte(',')
	buf.WriteString(polysign.XHeaderPolySignTimestamp)
	buf.WriteByte(',')
	buf.WriteString(polysign.XHeaderPolySignVersion)
	return buf.String()
}
