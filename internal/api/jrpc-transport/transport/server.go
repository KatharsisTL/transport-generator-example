// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"io"
)

const maxRequestBodySize = 104857600
const headerRequestID = "X-Request-Id"

type Server struct {
	log zerolog.Logger

	httpAfter  []Handler
	httpBefore []Handler

	config fiber.Config

	srvHTTP   *fiber.App
	srvHealth *fiber.App

	reporterCloser io.Closer
	httpHello      *httpHello
}

func New(log zerolog.Logger, options ...Option) (srv *Server) {

	srv = &Server{
		config: fiber.Config{
			BodyLimit:             maxRequestBodySize,
			DisableStartupMessage: true,
		},
		log: log,
	}
	for _, option := range options {
		option(srv)
	}
	srv.srvHTTP = fiber.New(srv.config)
	srv.srvHTTP.Post("/", srv.serveBatch)
	for _, option := range options {
		option(srv)
	}
	return
}

func (srv *Server) Fiber() *fiber.App {
	return srv.srvHTTP
}

func (srv *Server) WithLog(log zerolog.Logger) *Server {
	if srv.httpHello != nil {
		srv.httpHello = srv.Hello().WithLog(log)
	}
	return srv
}

func (srv *Server) ServeHealth(address string, response interface{}) {
	srv.srvHealth = fiber.New(fiber.Config{DisableStartupMessage: true})
	srv.srvHealth.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(response)
	})
	go func() {
		err := srv.srvHealth.Listen(address)
		ExitOnError(srv.log, err, "serve health on "+address)
	}()
}

func sendResponse(log zerolog.Logger, ctx *fiber.Ctx, resp interface{}) (err error) {
	ctx.Response().Header.SetContentType("application/json")
	if err = json.NewEncoder(ctx).Encode(resp); err != nil {
		log.Error().Err(err).Str("body", string(ctx.Body())).Msg("response write error")
	}
	return
}

func (srv *Server) Shutdown() {
	if srv.srvHTTP != nil {
		_ = srv.srvHTTP.Shutdown()
	}
	if srv.srvHealth != nil {
		_ = srv.srvHealth.Shutdown()
	}
	if srvMetrics != nil {
		_ = srvMetrics.Shutdown()
	}
}

func (srv *Server) WithTrace() *Server {
	if srv.httpHello != nil {
		srv.httpHello = srv.Hello().WithTrace()
	}
	return srv
}

func (srv *Server) WithMetrics() *Server {
	if srv.httpHello != nil {
		srv.httpHello = srv.Hello().WithMetrics()
	}
	return srv
}

func (srv Server) Hello() *httpHello {
	return srv.httpHello
}
