package service

import (
	"context"
	"github.com/KatharsisTL/transport-generator-example/internal/api/config"
	"github.com/rs/zerolog"
)

type service struct {
	ctx    context.Context
	logger zerolog.Logger
	config *config.ServiceConfig
}

func New(ctx context.Context, cfg *config.ServiceConfig, log zerolog.Logger) *service {
	return &service{
		ctx:    ctx,
		logger: log,
		config: cfg,
	}
}

func (s *service) Stop() error {
	s.logger.Info().Msg("service has been stopped")

	return nil
}
