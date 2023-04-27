package mq

import (
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/mq/rocket_mq"
	"go.uber.org/fx"
)

var Module = fx.Options(
	rocket_mq.RocketMqModule,
	rocket_mq.RocketMqProducerModule,
)
