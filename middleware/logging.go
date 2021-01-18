package middleware

import "net/http"

type GenID func() string

func Logging() Middleware {
	return func(next func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request) {

		return func(resp http.ResponseWriter, req *http.Request) {

		}
	}
}
