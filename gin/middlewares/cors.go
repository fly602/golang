package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ACCESS_CONTROL_ALLOW_ORIGIN      = "Access-Control-Allow-Origin"
	ACCESS_CONTROL_ALLOW_METHODS     = "Access-Control-Allow-Methods"
	ACCESS_CONTROL_ALLOW_HEADERS     = "Access-Control-Allow-Headers"
	ACCESS_CONTROL_EXPOSE_HEADERS    = "Access-Control-Expose-Headers"
	ACCESS_CONTROL_ALLOW_CREDENTIALS = "Access-Control-Allow-Credentials"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header(ACCESS_CONTROL_ALLOW_ORIGIN, "*")
		c.Header(ACCESS_CONTROL_ALLOW_METHODS, "GET,POST,OPTIONS,PUT,DELETE,UPDATE")
		c.Header(ACCESS_CONTROL_ALLOW_HEADERS, "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header(ACCESS_CONTROL_EXPOSE_HEADERS, "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header(ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")

		if method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusNoContent, "Option Request")
		}

	}
}
