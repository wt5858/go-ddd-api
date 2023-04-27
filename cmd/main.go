package main

import (
	router "github.com/wt5858/go-ddd-api/adpter/http/routers"
	"github.com/wt5858/go-ddd-api/infrastructure/conf"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/core"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/database"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/log"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/mq"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/search"
	"go.uber.org/fx"
)

var Module = fx.Options(
	conf.Module,
	log.Module,
	database.Module,
	search.Module,
	mq.Module,
	core.Module,   // http
	router.Module, // router handle
)

func main() {
	err := fx.ValidateApp(Module)
	if err != nil {
		panic(err)
	}

	fx.New(Module).Run()

	//defer system.close
}
