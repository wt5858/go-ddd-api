package handle

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var CfgHealthHandler = fx.Invoke(func(g *gin.Engine) {
	g.GET("health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"time":   time.Now().Format("2006-01-02 15:04:05"),
			"status": "ok",
		})
	})
})
