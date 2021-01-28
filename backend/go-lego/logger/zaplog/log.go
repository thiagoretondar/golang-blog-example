package zaplog

import (
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCustomZap(config logger.Config) (Logger, error) {
	var logConfig zap.Config

	if config.Production {
		logConfig = zap.NewProductionConfig()
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.EncoderConfig.TimeKey = "@timestamp"
		logConfig.DisableStacktrace = true
		logConfig.DisableCaller = false
	} else {
		// TODO don't call swat if other than Production
		logConfig = zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.DisableStacktrace = true
		logConfig.DisableCaller = true
		logConfig.EncoderConfig.CallerKey = "caller"
	}

	// Parse and verify given configuration log level
	if config.LogLevel != "" {
		errLogLevel := logConfig.Level.UnmarshalText([]byte(config.LogLevel))
		if errLogLevel != nil {
			panic(errLogLevel)
		}
	}

	// create zaplogger with configurations
	zapLog, err := logConfig.Build(
		// supplying this option prevents zaplog from always reporting the wrapper code as the caller
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, err
	}
	defer zapLog.Sync()

	// Override zaplog default zaplogger
	zap.ReplaceGlobals(zapLog)

	return &zaplogger{logger: zapLog}, nil
}
