package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log2 "github.com/wt5858/go-ddd-api/infrastructure/pkg/log"
	"go.uber.org/zap"
)

func AddTraceId(l *log2.Logger) gin.HandlerFunc {
	return func(g *gin.Context) {
		traceId := g.GetHeader("traceId")
		if traceId == "" {
			traceId = uuid.New().String()
		}

		ctx, log := l.AddCtx(g.Request.Context(), zap.Any("traceId", traceId))
		g.Request = g.Request.WithContext(ctx)
		log.Info("AddTraceId success")

		g.Next()
	}
}
