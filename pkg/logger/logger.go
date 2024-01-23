package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type settings struct {
	isdebug bool
}

type Option func(*settings)

func ProdLogger() Option {
	return func(s *settings) {
		s.isdebug = false
	}
}

func DevLogger() Option {
	return func(s *settings) {
		s.isdebug = true
	}
}

func Create(options ...Option) (*zap.Logger, error) {
	stdout := zapcore.AddSync(os.Stdout)

	settings := settings{}
	for _, op := range options {
		op(&settings)
	}

	cfg := zap.NewDevelopmentEncoderConfig()
	lvl := zap.DebugLevel
	stdoutenc := zapcore.NewConsoleEncoder(cfg)
	if settings.isdebug {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionEncoderConfig()
		lvl = zap.InfoLevel
		cfg.TimeKey = "timestamp"
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder
		stdoutenc = zapcore.NewJSONEncoder(cfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(stdoutenc, stdout, zap.NewAtomicLevelAt(lvl)),
	}
	return zap.New(zapcore.NewTee(cores...)), nil
}
