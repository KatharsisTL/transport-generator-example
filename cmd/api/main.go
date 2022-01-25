package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/KatharsisTL/transport-generator-example/internal/api/config"
	"github.com/KatharsisTL/transport-generator-example/internal/api/jrpc-transport/transport"
	"github.com/KatharsisTL/transport-generator-example/internal/api/service"
	"github.com/KatharsisTL/transport-generator-example/pkg/sig"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

func main() {
	serverConfig := config.Server()
	logger := newLogger(serverConfig.LogLevel)
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return sig.Listen(ctx)
	})

	svc := service.New(ctx, &serverConfig, logger)

	jrpcServer := transport.New(
		logger,
		transport.Hello(transport.NewHello(logger, svc)),
	).WithLog(logger).WithMetrics().WithTrace()
	jrpcServer.ServeHealth(serverConfig.HealthBind, http.StatusOK)
	transport.ServeMetrics(logger, ":8000")

	g.Go(makeServerRunner(ctx, logger, jrpcServer.Fiber(), serverConfig.Bind))
	if err := g.Wait(); err != nil {
		if !errors.Is(err, sig.ErrShutdownSignalReceived) {
			log.Error().Err(err).Msg("failed to wait error group")
		}

		if err := svc.Stop(); err != nil {
			log.Error().Err(err).Msg("failed to stopped service")
		}
	}
}

func newLogger(logLevel string) zerolog.Logger {
	ll, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level")
	}

	return zerolog.New(os.Stdout).Level(ll).With().Timestamp().Logger()
}

func makeServerRunner(ctx context.Context, logger zerolog.Logger, server *fiber.App, bindAddr string) func() error {
	return func() error {
		errCh := make(chan error)
		go func() {
			errCh <- server.Listen(bindAddr)
		}()
		logger.Info().Msg("http server has been started")
		select {
		case err := <-errCh:
			logger.Error().Err(err).Msg("http server failed")
			return err
		case <-ctx.Done():
			logger.Info().Msg("http server is stopping")
			err := server.Shutdown()
			if err == nil || err == http.ErrServerClosed {
				logger.Info().Msg("http server has been stopped")
				return nil
			}
			logger.Error().Err(err).Msg("failed to stop http server")
			return fmt.Errorf("shutdowning http server: %w", err)
		}
	}
}
