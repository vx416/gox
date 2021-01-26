package middleware

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type (
	JWTTokenKey struct{}
)

func JWTToken(verifyFunc func(token *jwt.Token) (interface{}, error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			tokenStr := extractToken(req.Header, "Bearer", HeaderAuth)
			if tokenStr == "" {
				return
			}
			token, err := jwt.Parse(tokenStr, verifyFunc)
			if err != nil || !token.Valid {
				return
			}
			ctx := context.WithValue(req.Context(), JWTTokenKey{}, token)
			req = req.WithContext(ctx)
			next.ServeHTTP(resp, req)
		})
	}
}

func extractToken(header http.Header, schema, headerKey string) string {
	tokenHeader := header.Get(headerKey)
	if tokenHeader == "" {
		return ""
	}
	if len(schema)+1 > len(tokenHeader) || tokenHeader[:len(schema)] != schema {
		return ""
	}

	return tokenHeader[len(schema)+1:]
}
