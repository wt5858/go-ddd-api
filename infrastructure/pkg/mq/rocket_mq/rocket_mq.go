package rocket_mq

import (
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/log"
	"go.uber.org/zap"

	"github.com/wt5858/go-ddd-api/infrastructure/conf"

	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/fx"
)

var RocketMqModule = fx.Provide(func(cfg *conf.Config, log *log.Logger) rocketmq.PushConsumer {
	r, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(cfg.RocketMqConfig.Group),
		consumer.WithConsumerOrder(true), // 增加顺序消费
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.RocketMqConfig.Host + ":" + cfg.RocketMqConfig.Port})),
	)

	if err != nil {
		log.ZapLogger.Error("[rocket_mq-consumer-error]", zap.Any("error", err.Error()))
		panic(err)
	}
	return r
})

var RocketMqProducerModule = fx.Provide(func(cfg *conf.Config, log *log.Logger) rocketmq.Producer {
	r, err := rocketmq.NewProducer(
		producer.WithGroupName(cfg.RocketMqConfig.Group),
		producer.WithRetry(2),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{cfg.RocketMqConfig.Host + ":" + cfg.RocketMqConfig.Port})),
	)
	if err != nil {
		log.ZapLogger.Error("[rocket_mq-producer-error]", zap.Any("error", err.Error()))
		panic(err)
	}

	err = r.Start()
	if err != nil {
		log.ZapLogger.Error("[rocket_mq-producer-error] start producer error", zap.Any("error", err.Error()))
		panic(err)
	}
	return r
})
