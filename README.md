[![Go Reference](https://pkg.go.dev/badge/github.com/stupidwebauthn/swa_sdk_go.svg)](https://pkg.go.dev/github.com/stupidwebauthn/swa_sdk_go)

# Stupid Webauthn Sdk <small>GoLang</small>

## Gin example

```golang
import (
	"github.com/gin-gonic/gin"
	"github.com/stupidwebauthn/swa_sdk_go"
)

r.Use(func(c *gin.Context) {
	header := c.Writer.Header()
	res, status, err := swa.AuthMiddleware(c.Request, &header)
	if err != nil {
		c.AbortWithError(status, err)
		return
	}
	c.Set("swa", res)
	c.Next()
})
r.GET("/api/data", func(c *gin.Context) {
	rawAuth, _ := c.Get("swa")
	auth := rawAuth.(*swa_sdk_go.AuthResponse)
	...
})
```