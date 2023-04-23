package base

import (
	"github.com/gin-gonic/gin"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/core/middleware"
	"go.uber.org/fx"
)

var RegisterGin = fx.Provide(func() *gin.Engine {
	gin.SetMode(gin.TestMode)
	g := gin.Default()
	g.MaxMultipartMemory = 8 << 20

	// 中间件
	g.Use(middleware.Cors())

	return g
})
