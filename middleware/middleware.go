package middleware

import "net/http"

type Middleware func(next func(resp http.ResponseWriter, req *http.Request)) func(resp http.ResponseWriter, req *http.Request)
