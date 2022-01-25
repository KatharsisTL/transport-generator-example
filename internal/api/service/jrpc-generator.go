// @tg title=`hello-api service`
// @tg version=0.0.1
// @tg description=`инициализация сервиса`
// @tg servers=`http://localhost:9000`
//go:generate tg transport --services . --out ../jrpc-transport/transport --outSwagger ../jrpc-transport/swagger.yaml
package service

import (
	"context"
)

// @tg jsonRPC-server log trace metrics
type Hello interface {
	// @tg desc=`возвращает hello world`
	Hello(ctx context.Context, name string) (resp string, err error)
}
