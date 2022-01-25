// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"strings"
	"sync"
)

func (http *httpHello) serveHello(ctx *fiber.Ctx) (err error) {
	return http.serveMethod(ctx, "hello", http.hello)
}
func (http *httpHello) hello(span opentracing.Span, ctx *fiber.Ctx, requestBase baseJsonRPC) (responseBase *baseJsonRPC) {
	var err error
	var request requestHelloHello
	if requestBase.Params != nil {
		if err = json.Unmarshal(requestBase.Params, &request); err != nil {
			ext.Error.Set(span, true)
			span.SetTag("msg", "request body could not be decoded: "+err.Error())
			return makeErrorResponseJsonRPC(requestBase.ID, parseError, "request body could not be decoded: "+err.Error(), nil)
		}
	}
	if requestBase.Version != Version {
		ext.Error.Set(span, true)
		span.SetTag("msg", "incorrect protocol version: "+requestBase.Version)
		return makeErrorResponseJsonRPC(requestBase.ID, parseError, "incorrect protocol version: "+requestBase.Version, nil)
	}
	methodContext := opentracing.ContextWithSpan(ctx.Context(), span)

	var response responseHelloHello
	response.Resp, err = http.svc.Hello(methodContext, request.Name)
	if err != nil {
		if http.errorHandler != nil {
			err = http.errorHandler(err)
		}
		ext.Error.Set(span, true)
		span.SetTag("msg", err)
		span.SetTag("errData", toString(err))
		return makeErrorResponseJsonRPC(requestBase.ID, internalError, err.Error(), err)
	}
	responseBase = &baseJsonRPC{
		ID:      requestBase.ID,
		Version: Version,
	}
	if responseBase.Result, err = json.Marshal(response); err != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "response body could not be encoded: "+err.Error())
		return makeErrorResponseJsonRPC(requestBase.ID, parseError, "response body could not be encoded: "+err.Error(), nil)
	}
	return
}
func (http *httpHello) serveBatch(ctx *fiber.Ctx) (err error) {
	batchSpan := extractSpan(http.log, fmt.Sprintf("jsonRPC:%s", ctx.Path()), ctx)
	defer injectSpan(http.log, batchSpan, ctx)
	defer batchSpan.Finish()
	methodHTTP := ctx.Method()
	if methodHTTP != fiber.MethodPost {
		ext.Error.Set(batchSpan, true)
		batchSpan.SetTag("msg", "only POST method supported")
		ctx.Response().SetStatusCode(fiber.StatusMethodNotAllowed)
		if _, err = ctx.WriteString("only POST method supported"); err != nil {
			return
		}
		return
	}
	if value := ctx.Context().Value(CtxCancelRequest); value != nil {
		return
	}
	var single bool
	var requests []baseJsonRPC
	if err = json.Unmarshal(ctx.Body(), &requests); err != nil {
		var request baseJsonRPC
		if err = json.Unmarshal(ctx.Body(), &request); err != nil {
			ext.Error.Set(batchSpan, true)
			batchSpan.SetTag("msg", "request body could not be decoded: "+err.Error())
			return sendResponse(http.log, ctx, makeErrorResponseJsonRPC([]byte("\"0\""), parseError, "request body could not be decoded: "+err.Error(), nil))
		}
		single = true
		requests = append(requests, request)
	}
	responses := make(jsonrpcResponses, 0, len(requests))
	var wg sync.WaitGroup
	for _, request := range requests {
		methodNameOrigin := request.Method
		method := strings.ToLower(request.Method)
		span := opentracing.StartSpan(request.Method, opentracing.ChildOf(batchSpan.Context()))
		span.SetTag("batch", true)
		switch method {
		case "hello":
			wg.Add(1)
			func(request baseJsonRPC) {
				responses.append(http.hello(span, ctx, request))
				wg.Done()
			}(request)
		default:
			ext.Error.Set(span, true)
			span.SetTag("msg", "invalid method '"+methodNameOrigin+"'")
			responses.append(makeErrorResponseJsonRPC(request.ID, methodNotFoundError, "invalid method '"+methodNameOrigin+"'", nil))
		}
		span.Finish()
	}
	wg.Wait()
	if single {
		return sendResponse(http.log, ctx, responses[0])
	}
	return sendResponse(http.log, ctx, responses)
}
func (http *httpHello) serveMethod(ctx *fiber.Ctx, methodName string, methodHandler methodTraceJsonRPC) (err error) {
	span := extractSpan(http.log, fmt.Sprintf("jsonRPC:%s", ctx.Path()), ctx)
	defer injectSpan(http.log, span, ctx)
	defer span.Finish()
	methodHTTP := ctx.Method()
	if methodHTTP != fiber.MethodPost {
		ext.Error.Set(span, true)
		span.SetTag("msg", "only POST method supported")
		ctx.Response().SetStatusCode(fiber.StatusMethodNotAllowed)
		if _, err = ctx.WriteString("only POST method supported"); err != nil {
			return
		}
	}
	if value := ctx.Context().Value(CtxCancelRequest); value != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "request canceled")
		return
	}
	var request baseJsonRPC
	var response *baseJsonRPC
	if err = json.Unmarshal(ctx.Body(), &request); err != nil {
		ext.Error.Set(span, true)
		span.SetTag("msg", "request body could not be decoded: "+err.Error())
		return sendResponse(http.log, ctx, makeErrorResponseJsonRPC([]byte("\"0\""), parseError, "request body could not be decoded: "+err.Error(), nil))
	}
	methodNameOrigin := request.Method
	method := strings.ToLower(request.Method)
	if method != "" && method != methodName {
		ext.Error.Set(span, true)
		span.SetTag("msg", "invalid method "+methodNameOrigin)
		return sendResponse(http.log, ctx, makeErrorResponseJsonRPC(request.ID, methodNotFoundError, "invalid method "+methodNameOrigin, nil))
	}
	response = methodHandler(span, ctx, request)
	if response != nil {
		return sendResponse(http.log, ctx, response)
	}
	return
}
