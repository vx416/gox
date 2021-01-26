package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/vx416/gox/log"
)

func Recovery() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("err: %+v", r)
					}
					ctx := req.Context()
					logger := log.Ctx(ctx)
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, false)
					msg := fmt.Sprintf("[PANIC RECOVER] %s \nerr:%+v", err, stack[:length])
					logger.Error(msg)
				}
			}()

			next.ServeHTTP(resp, req)
		})
	}
}
