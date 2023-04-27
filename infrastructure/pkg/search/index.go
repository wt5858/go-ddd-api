package search

import (
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/search/elastic"
	"go.uber.org/fx"
)

var Module = fx.Options(
	elastic.EsModule,
)
