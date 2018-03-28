package zenrpc_mw

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/semrush/zenrpc"
)

// This middleware is compilation of go-kit/kit/tracing/opentracing functions.
func Tracing(tracer opentracing.Tracer) zenrpc.MiddlewareFunc {
	return func(invoke zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			operationName := zenrpc.NamespaceFromContext(ctx) + "." + method

			// Step one
			// Try to join to a trace propagated in `req`.
			if req, ok := zenrpc.RequestFromContext(ctx); ok {
				var span opentracing.Span
				wireContext, _ := tracer.Extract(
					opentracing.TextMap,
					opentracing.HTTPHeadersCarrier(req.Header),
				)
				span = tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
				ext.HTTPMethod.Set(span, req.Method)
				ext.HTTPUrl.Set(span, req.URL.String())
				ctx = opentracing.ContextWithSpan(ctx, span)
			}

			// Step two
			// Trace server
			serverSpan := opentracing.SpanFromContext(ctx)
			if serverSpan == nil {
				// All we can do is create a new root span.
				serverSpan = tracer.StartSpan(method)
			} else {
				serverSpan.SetOperationName(method)
			}
			defer serverSpan.Finish()
			ext.SpanKindRPCServer.Set(serverSpan)
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)

			return invoke(ctx, method, params)
		}
	}
}
