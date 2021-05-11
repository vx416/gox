package interceptor

import (
	"context"
	"time"

	"github.com/rs/xid"
	"github.com/vx416/gox/ctxutil"
	"github.com/vx416/gox/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerLogging(logger log.Logger, reqDump, respDump bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(map[string]string{
				RequestIDKey: xid.New().String(),
			})
		}
		reqID := ""
		vals := md.Get(RequestIDKey)
		if len(vals) == 1 {
			reqID = vals[0]
		} else {
			reqID = xid.New().String()
		}

		logFields := map[string]interface{}{
			"grpc_func":  info.FullMethod,
			"request_id": reqID,
		}
		if reqDump {
			logFields["request_dump"] = dumpStruct(req)
		}

		logger = logger.Fields(logFields)
		ctx = logger.Attach(ctx)
		start := time.Now()
		resp, err = handler(ctx, req)
		if respDump {
			logger = logger.Field("response_dump", dumpStruct(resp))
		}
		logger = logger.Field("latency", time.Since(start).String())
		if err != nil {
			logger.Err(err).Error("request failed")
		} else {
			logger.Info("request success")
		}

		return resp, err
	}
}

func ClientRequestID(
	ctx context.Context, method string,
	req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	ctx, reqID := ctxutil.GetReqID(ctx)
	md := metadata.New(map[string]string{
		RequestIDKey: reqID,
	})
	ctx = metadata.NewIncomingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func dumpStruct(st interface{}) string {
	return ""
}
