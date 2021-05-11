package ctxutil

import (
	"context"

	"github.com/rs/xid"
)

type (
	RequestIDKey struct{}
)

func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, reqID)
}

func GetReqID(ctx context.Context) (context.Context, string) {
	val := ctx.Value(RequestIDKey{})
	if val == nil {
		reqID := xid.New().String()
		ctx = WithReqID(ctx, reqID)
		return ctx, reqID
	}

	if reqID, ok := val.(string); ok {
		return ctx, reqID
	}

	reqID := xid.New().String()
	ctx = WithReqID(ctx, reqID)
	return ctx, reqID
}
