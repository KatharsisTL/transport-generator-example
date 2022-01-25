// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

type ServiceRoute interface {
	SetRoutes(route *fiber.App)
}

type Option func(srv *Server)
type Handler = fiber.Handler
type ErrorHandler func(err error) error

func Service(svc ServiceRoute) Option {
	return func(srv *Server) {
		if srv.srvHTTP != nil {
			svc.SetRoutes(srv.Fiber())
		}
	}
}

func Hello(svc *httpHello) Option {
	return func(srv *Server) {
		if srv.srvHTTP != nil {
			srv.httpHello = svc
			svc.SetRoutes(srv.Fiber())
		}
	}
}

func MaxBodySize(max int) Option {
	return func(srv *Server) {
		srv.config.BodyLimit = max
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(srv *Server) {
		srv.config.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(srv *Server) {
		srv.config.WriteTimeout = timeout
	}
}

func Use(args ...interface{}) Option {
	return func(srv *Server) {
		if srv.srvHTTP != nil {
			srv.srvHTTP.Use(args...)
		}
	}
}
