package middleware

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/vx416/gox/log"
)

type (
	bodyDumpResponseWriter struct {
		io.Writer
		http.ResponseWriter
		status int
	}
)

func Logging(logger log.Logger, reqDump, respDump bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			reqBody, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

			respBody := new(bytes.Buffer)
			mw := io.MultiWriter(respBody, resp)
			w := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: resp}
			resp = w

			logger = logger.Field("request_id", req.Header.Get(HeaderXRequestID)).
				Field("method", req.Method).Field("path", req.URL.RawPath)
			ctx := logger.Attach(req.Context())
			req = req.WithContext(ctx)

			logFields := make(map[string]interface{})
			logFields["uri"] = req.URL.Host
			logFields["bytes_in"] = req.Header.Get("Content-Length")

			if reqDump {
				logFields["req_dump"] = string(reqBody)
			}

			start := time.Now()
			defer func() {
				end := time.Now()
				logger = logger.Fields(logFields).Field("status", http.StatusText(w.status)).
					Field("latency", end.Sub(start).String()).
					Field("bytes_out", respBody.Len())

				if respDump {
					logger = logger.Field("resp_dump", respBody.String())
				}
				logger.Info("access log")
			}()

			next.ServeHTTP(resp, req)
		})
	}
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
