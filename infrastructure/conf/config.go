package conf

import (
	"time"

	"go.uber.org/fx"
)

type App struct {
	ListenAddr string
	Debug      bool
	Stat       bool
	StaticDir  string
}

type MySQLConf struct {
	Driver    string
	MasterDsn string
	SlaveDsn  string
}

type RedisConf struct {
	Protocol    string
	Host        string
	Port        string
	Auth        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Db          string
}

type LimiterConf struct {
	Speed      float64
	Concurrent int
	Timeout    time.Duration
}

type SwaggerConfig struct {
	SwitchConfig bool
}

type Config struct {
	App
	MySQLConf
	RedisConf
	LimiterConf
	SwaggerConfig
}

func init() {
	if err := InitConfig(""); err != nil {
		panic(err)
	}
}

var Module = fx.Provide(func() *Config {
	return &Config{
		App{
			ListenAddr: GetKeyByConf("app.host", "str").(string) + GetKeyByConf("app.port", "str").(string),
			Debug:      GetKeyByConf("app.debug", "bool").(bool),
			Stat:       GetKeyByConf("app.stat", "bool").(bool),
			StaticDir:  GetKeyByConf("app.static", "str").(string),
		},
		MySQLConf{
			Driver:    GetKeyByConf("mysql.driver", "str").(string),
			MasterDsn: GetKeyByConf("mysql.master.dsn", "str").(string),
			SlaveDsn:  GetKeyByConf("mysql.slave.dsn", "str").(string),
		},
		RedisConf{
			Protocol:    GetKeyByConf("redis.protocol", "str").(string),
			Host:        GetKeyByConf("redis.host", "str").(string),
			Port:        GetKeyByConf("redis.port", "str").(string),
			Auth:        GetKeyByConf("redis.auth", "str").(string),
			MaxIdle:     GetKeyByConf("redis.maxIdle", "int").(int),
			MaxActive:   GetKeyByConf("redis.maxActive", "int").(int),
			IdleTimeout: GetKeyByConf("redis.idleTimeout", "duration").(time.Duration),
			Db:          GetKeyByConf("redis.db", "str").(string),
		},
		LimiterConf{
			Speed:      GetKeyByConf("limiter.speed", "float64").(float64),
			Concurrent: GetKeyByConf("limiter.concurrent", "int").(int),
			Timeout:    GetKeyByConf("limiter.timeout", "duration").(time.Duration),
		},
		SwaggerConfig{
			SwitchConfig: true,
		},
	}
})
