package application

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GoGoTools struct {
	Log *zap.SugaredLogger
}

func New() *GoGoTools {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.TimeKey = zapcore.OmitKey
	encCfg.LevelKey = zapcore.OmitKey
	encCfg.NameKey = zapcore.OmitKey
	encCfg.CallerKey = zapcore.OmitKey
	encCfg.FunctionKey = zapcore.OmitKey
	encCfg.StacktraceKey = zapcore.OmitKey

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := zap.Must(cfg.Build())
	return &GoGoTools{Log: logger.Sugar()}
}

func (ggt *GoGoTools) Verbose() {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.TimeKey = zapcore.OmitKey
	encCfg.LevelKey = zapcore.OmitKey
	encCfg.NameKey = zapcore.OmitKey
	encCfg.CallerKey = zapcore.OmitKey
	encCfg.FunctionKey = zapcore.OmitKey
	encCfg.StacktraceKey = zapcore.OmitKey

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger := zap.Must(cfg.Build())
	ggt.Log = logger.Sugar()
}
