package mysql

import (
	"database/sql"

	"github.com/wt5858/go-ddd-api/infrastructure/conf"
	"github.com/wt5858/go-ddd-api/infrastructure/pkg/log"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db      *DataBase
	connect *gorm.DB
)

type DataBase struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

type GormLogger struct {
	log *log.Logger
}

func (gl *GormLogger) Print(v ...interface{}) {
	gl.log.ZapLogger.Info("[sql-info]",
		zap.String("module", "gorm"),
		zap.String("type", "sql"),
		zap.Any("sql", v[3]),
		zap.Any("values", v[4]),
		zap.Any("duration", v[2]),
		zap.Any("src", v[1]),
		zap.Any("rows_returned", v[5]),
	)
}

func getConnect(driver, dsn string, logger *log.Logger) *gorm.DB {
	sqlDB, _ := sql.Open(driver, dsn)
	connect, _ = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:                 logger,
		SkipDefaultTransaction: false,
	})
	connect.Debug()
	db, _ := connect.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	return connect
}

var Module = fx.Provide(func(cfg *conf.Config, logger *log.Logger) *DataBase {
	db = &DataBase{
		Master: getConnect(
			cfg.MySQLConf.Driver,
			cfg.MySQLConf.MasterDsn,
			logger,
		),
		Slave: getConnect(
			cfg.MySQLConf.Driver,
			cfg.MySQLConf.SlaveDsn,
			logger,
		),
	}
	return db
})
