package logging

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"sync"
)

var (
	Logger *zap.Logger
	once   sync.Once
)

func InitLogging() {
	once.Do(func() {
		logFile := &lumberjack.Logger{
			Filename:   "internal/common/logs.log", // Log file path
			MaxSize:    100,                        // Max size in MB before rotation
			MaxBackups: 3,                          // Max number of backups
			MaxAge:     14,                         // Max number of days to keep old logs
			Compress:   true,                       // Enable compression (gzipped backups)
		}

		encoderConfig := zapcore.EncoderConfig{
			TimeKey:       "timestamp",
			LevelKey:      "level",
			MessageKey:    "message",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		}

		encoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel)

		core := zapcore.NewTee(fileCore)
		Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	})

	// Log messages
	//logger.Info("This is an info message with a stack trace.")
	//logger.Error("This is an error message with a stack trace.")
}
