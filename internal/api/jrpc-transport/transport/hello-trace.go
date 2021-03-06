// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"context"
	"github.com/KatharsisTL/transport-generator-example/internal/api/service"
	"github.com/opentracing/opentracing-go"
)

type traceHello struct {
	next service.Hello
}

func traceMiddlewareHello(next service.Hello) service.Hello {
	return &traceHello{next: next}
}

func (svc traceHello) Hello(ctx context.Context, name string) (resp string, err error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("method", "Hello")
	return svc.next.Hello(ctx, name)
}
