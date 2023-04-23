package router

import (
	"context"
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/wt5858/go-ddd-api/infrastructure/conf"
	"go.uber.org/fx"
)

var CfgGin = fx.Invoke(func(lifecycle fx.Lifecycle, g *gin.Engine, cfg *conf.Config) {
	httpServer := http.Server{
		Addr:    cfg.ListenAddr,
		Handler: g,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if cfg.Debug {
					// http://127.0.0.1:18000/debug/pprof/
					ginpprof.Wrap(g)
				}

				// http: //127.0.0.1:18000/swagger/index.html
				//if cfg.SwaggerConfig.SwitchConfig == true {
				//	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
				//}

				err := httpServer.ListenAndServe()
				if err != nil {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})

})
