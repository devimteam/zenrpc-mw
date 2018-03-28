package zenrpc_mw

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/semrush/zenrpc"
)

func RequestCounter(counter metrics.Counter) zenrpc.MiddlewareFunc {
	return func(invoke zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			r := invoke(ctx, method, params)

			if namespace := zenrpc.NamespaceFromContext(ctx); namespace != "" {
				method = namespace + "." + method
			}
			code := ""
			if r.Error != nil {
				code = strconv.Itoa(r.Error.Code)
			}
			counter.With("method", method, "code", code).Add(1)
			return r
		}
	}
}

func RequestDuration(histogram metrics.Histogram) zenrpc.MiddlewareFunc {
	return func(invoke zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			begin := time.Now()
			r := invoke(ctx, method, params)

			if namespace := zenrpc.NamespaceFromContext(ctx); namespace != "" {
				method = namespace + "." + method
			}
			code := ""
			if r.Error != nil {
				code = strconv.Itoa(r.Error.Code)
			}
			histogram.With("method", method, "code", code).Observe(time.Since(begin).Seconds())
			return r
		}
	}
}
