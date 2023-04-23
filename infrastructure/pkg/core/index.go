package core

import (
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/core/base"
	"go.uber.org/fx"
)

var Module = fx.Options(
	base.RegisterGin,
)
