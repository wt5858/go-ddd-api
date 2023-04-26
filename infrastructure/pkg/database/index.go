package database

import (
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/database/mysql"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/database/redis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	redis.Module,
	mysql.Module,
)
