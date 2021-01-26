package middleware

import (
	"net/http"

	"github.com/rs/xid"
)

func XidGent() string {
	return xid.New().String()
}

func RequestID(genID func() string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			reqID := req.Header.Get(HeaderXRequestID)
			if reqID == "" {
				resp.Header().Set(HeaderXRequestID, genID())
			}

			next.ServeHTTP(resp, req)
		})
	}
}
