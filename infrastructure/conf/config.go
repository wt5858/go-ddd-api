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

type LogConfig struct {
	LogFileDir string
	AppName    string
	Platform   string
	MaxSize    int //文件多大开始切分
	MaxBackups int //保留文件个数
	MaxAge     int //文件保留最大实际
	Level      string
}

type Config struct {
	App
	MySQLConf
	RedisConf
	LimiterConf
	SwaggerConfig
	LogConfig
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
		LogConfig{
			LogFileDir: GetKeyByConf("log.logFileDir", "str").(string),
			AppName:    GetKeyByConf("log.appName", "str").(string),
			Platform:   GetKeyByConf("log.platform", "str").(string),
			MaxSize:    GetKeyByConf("log.maxSize", "int").(int),
			MaxBackups: GetKeyByConf("log.maxBackups", "int").(int),
			MaxAge:     GetKeyByConf("log.maxAge", "int").(int),
			Level:      GetKeyByConf("log.level", "str").(string),
		},
	}
})
