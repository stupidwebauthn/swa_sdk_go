# Stupid Webauthn Sdk <small>GoLang</small>

## Gin example

```golang
import (
	"github.com/gin-gonic/gin"
	"github.com/stupidwebauthn/swa_sdk_go"
)

r.Use(func(c *gin.Context) {
	res, status, err := swa.Middleware(c.Request)
	if err != nil {
		if status == http.StatusUnauthorized {
			swa.RemoveAuthCookie(c.Writer)
		}
		c.Status(status)
		c.Error(err)
		c.Abort()
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