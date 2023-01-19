package application

import (
	"go.uber.org/zap"
)

type GoGoTools struct {
	Log *zap.SugaredLogger
}

func New() *GoGoTools {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.InfoLevel)
	logger := zap.Must(cfg.Build())
	return &GoGoTools{Log: logger.Sugar()}
}

func (ggt *GoGoTools) Verbose() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.DebugLevel)
	logger := zap.Must(cfg.Build())
	ggt.Log = logger.Sugar()
}
