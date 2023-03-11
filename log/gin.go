package log

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.FullPath()
		if url == "" {
			url = "[unmatched]"
		}
		t := time.Now()
		c.Next()
		elapsed := time.Since(t)

		status := c.Writer.Status()
		if status >= 400 {
			Default.Errorf("[%s] %s (%d) from %s took %dms", c.Request.Method, url, c.Writer.Status(), c.ClientIP(), elapsed.Milliseconds())
		} else {
			Default.Infof("[%s] %s (%d) from %s took %dms", c.Request.Method, url, c.Writer.Status(), c.ClientIP(), elapsed.Milliseconds())
		}
	}
}
