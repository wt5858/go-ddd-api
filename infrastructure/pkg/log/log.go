package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/wt5858/go-ddd-api/infrastructure/conf"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

var (
	l              *Logger
	sp             = string(filepath.Separator)
	outWrite       zapcore.WriteSyncer       // IO输出
	debugConsoleWS = zapcore.Lock(os.Stdout) // 控制台标准输出
)

var Module = fx.Provide(NewLogger)

type Logger struct {
	ZapLogger *zap.Logger
	sync.RWMutex
	Opts                      *conf.Config
	zapConfig                 zap.Config
	inited                    bool
	SlowThreshold             time.Duration
	LogLevel                  gormLogger.LogLevel
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func init() {
	l = &Logger{
		Opts: &conf.Config{},
	}
}

func NewLogger(cf *conf.Config) *Logger {
	l.Lock()
	defer l.Unlock()
	if l.inited {
		l.ZapLogger.Info("[initLogger] logger Inited")
		return nil
	}
	l.Opts = cf
	l.loadCfg()
	l.init()
	l.ZapLogger.Info("[initLogger] zap plugin initializing completed")
	l.inited = true
	return l
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	l.ZapLogger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer l.ZapLogger.Sync()
}

func GetLogger() *Logger {
	if l == nil {
		fmt.Println("Please initialize the log service first")
		return nil
	}
	return l
}

func (l *Logger) GetLevel() (level zapcore.Level) {
	switch strings.ToLower(l.Opts.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel //默认为调试模式
	}
}

func (l *Logger) setSyncers() {
	outWrite = zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + ".log",
		MaxSize:    l.Opts.MaxSize,
		MaxBackups: l.Opts.MaxBackups,
		MaxAge:     l.Opts.MaxAge,
		Compress:   true,
		LocalTime:  true,
	})
	return
}

func (l *Logger) loadCfg() {
	if l.GetLevel() == zapcore.DebugLevel {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	l.zapConfig.OutputPaths = []string{"stdout"}
	l.zapConfig.OutputPaths = []string{"stderr"}
	// 默认输出到程序运行目录的logs子目录
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}
	if l.Opts.MaxSize == 0 {
		l.Opts.MaxSize = 100
	}
	if l.Opts.MaxBackups == 0 {
		l.Opts.MaxBackups = 60
	}
	if l.Opts.MaxAge == 0 {
		l.Opts.MaxAge = 30
	}
	l.SlowThreshold = time.Second
	l.LogLevel = gormLogger.Info
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l.GetLevel()
	})
	var cores []zapcore.Core
	if l.GetLevel() == zapcore.DebugLevel {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, debugConsoleWS, priority),
		}...)
	} else {
		if l.Opts.Platform == "k8s" {
			cores = append(cores, []zapcore.Core{
				zapcore.NewCore(fileEncoder, debugConsoleWS, priority),
			}...)
		} else {
			cores = append(cores, []zapcore.Core{
				zapcore.NewCore(fileEncoder, outWrite, priority),
			}...)
		}
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

func (l *Logger) GetCtx(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.Opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.ZapLogger
}

func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(l.Opts.CtxKey).(*zap.Logger)
	if ok {
		return log
	}
	return l.ZapLogger
}

func (l *Logger) AddCtx(ctx context.Context, field ...zap.Field) (context.Context, *zap.Logger) {
	log := l.ZapLogger.With(field...)
	ctx = context.WithValue(ctx, l.Opts.CtxKey, log)
	return ctx, log
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}
