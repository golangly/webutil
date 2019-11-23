package webutil

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"

	"github.com/golangly/log"
)

func RequestLogger(next http.Handler) http.Handler {
	logResponseBodyValue := strings.ToLower(os.Getenv("LOG_HTTP_RESPONSE_BODY"))
	logResponseBody := logResponseBodyValue == "1" || logResponseBodyValue == "yes" || logResponseBodyValue == "true" || logResponseBodyValue == "y"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Wrap the response to enable access to the final status code & response contents
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Tee the response to an additional buffer so we can look inside & log it
		buf := bytes.Buffer{}
		ww.Tee(&buf)

		// Invoke target handler, but time the request
		start := time.Now()
		next.ServeHTTP(ww, r)
		duration := time.Since(start)

		// Prepare log context
		logger := log.
			With("rid", w.Header().Get(RequestIDHeaderName)).
			With("remoteAddr", r.RemoteAddr).
			With("proto", r.Proto).
			With("method", r.Method).
			With("uri", r.RequestURI).
			With("host", r.Host).
			With("requestHeaders", r.Header).
			With("responseHeaders", w.Header()).
			With("elapsed", duration).
			With("bytesWritten", ww.BytesWritten()).
			With("status", ww.Status())

		// Add HTTP response bytes, if requested to
		if logResponseBody {
			logger = logger.With("bytesOut", buf.Bytes())
		}

		// Log it
		logger.Info("HTTP request completed")
	})
}
