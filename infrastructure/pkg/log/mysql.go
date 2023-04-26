package log

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func (l *Logger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Info {
		return
	}
	l.ZapLogger.Sugar().Infof(str, args...)
}

func (l *Logger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Warn {
		return
	}
	l.ZapLogger.Sugar().Warnf(str, args...)
}

func (l *Logger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormLogger.Error {
		return
	}
	l.ZapLogger.Sugar().Errorf(str, args...)
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &Logger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		l.ZapLogger.Error(
			"[error-sql]",
			zap.Error(err),
			zap.Any("elapsed", elapsed.String()),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		l.ZapLogger.Info(
			"[slow-sql]",
			zap.Any("elapsed", elapsed.String()),
			zap.Any("rows", rows),
			zap.Any("sql", sql),
		)
	case l.LogLevel >= gormLogger.Info:
		sql, rows := fc()
		l.ZapLogger.Info(
			"[sql-info]",
			zap.Any("elapsed", elapsed.String()),
			zap.Any("rows", rows),
			zap.Any("sql", sql),
		)
	}
}
