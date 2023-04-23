package router

import (
	"github.com/wt5858/go-ddd-api/adpter/http/routers/handle"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/core/router"
	"go.uber.org/fx"
)

var Module = fx.Options(
	handle.CfgHealthHandler,
	router.CfgGin,
)
