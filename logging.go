package zenrpc_mw

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/semrush/zenrpc"
)

func Logger(logger log.Logger) zenrpc.MiddlewareFunc {
	return func(invoke zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			begin, ip := time.Now(), "<nil>"
			if req, ok := zenrpc.RequestFromContext(ctx); ok && req != nil {
				ip = req.RemoteAddr
			}

			r := invoke(ctx, method, params)
			logger.Log(
				"ip", ip,
				"method", zenrpc.NamespaceFromContext(ctx)+"."+method,
				"duration", time.Since(begin),
				"err", r.Error,
			)
			return r
		}
	}
}
