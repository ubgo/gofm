package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Sugar *zap.SugaredLogger
	Plain *zap.Logger
}

func (logger Logger) Version() string {
	return "0.01"
}

func New() Logger {
	path := viper.GetString("logger.file")
	if len(path) == 0 {
		path = "./app.log"
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		LocalTime:  true,
	})
	core := zapcore.NewCore(
		getEncoder(),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	// logger.Info("failed to fetch URL",
	// 	zap.String("url", "http://example.com"),
	// 	zap.Int("attempt", 3),
	// 	zap.Duration("backoff", time.Second),
	// )

	return Logger{
		Sugar: logger.Sugar(),
		Plain: logger,
	}
}

func getEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(cfg)
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}
