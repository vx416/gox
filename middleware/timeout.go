package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(to time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			done := make(chan struct{})

			ctx, cancel := context.WithTimeout(req.Context(), to)
			defer cancel()
			req = req.WithContext(ctx)

			go func() {
				defer func() {
					close(done)
				}()
				next.ServeHTTP(resp, req)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				resp.WriteHeader(http.StatusRequestTimeout)
			}
		})
	}
}
