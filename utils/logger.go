package utils

import (
	"context"

	"github.com/CrowderSoup/emoji-match/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewZapLogger returns a new zap logger
func NewZapLogger(config *config.Config) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if config.Environment == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	return logger, err
}

// ProvideSugaredLogger provides a sugared logger to all who want one
func ProvideSugaredLogger(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

// InvokeLogger manages the loggers Lifecycle
func InvokeLogger(lc fx.Lifecycle, logger *zap.Logger) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Sync()
				return nil
			},
		},
	)
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		NewZapLogger,
		ProvideSugaredLogger,
		NewTraceExporter,
		NewTracerProvider,
	),
)
