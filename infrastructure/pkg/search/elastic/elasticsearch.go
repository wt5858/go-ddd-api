package elastic

import (
	"go.uber.org/zap"

	"github.com/wt5858/go-ddd-api/infrastructure/pkg/log"

	"github.com/wt5858/go-ddd-api/infrastructure/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
)

var EsModule = fx.Provide(func(cfg *conf.Config, log *log.Logger) *elasticsearch.Client {
	config := elasticsearch.Config{}
	// 可以修改为集群的方式
	config.Addresses = []string{cfg.EsConfig.Host + ":" + cfg.EsConfig.Port}
	config.Username = cfg.EsConfig.Username
	config.Password = cfg.EsConfig.Password
	esClient, err := elasticsearch.NewClient(config)
	if err != nil {
		log.ZapLogger.Error("[elastic-conn-error]", zap.Any("error", err.Error()))
		panic(err)
	}
	return esClient
})
