package webutil

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0
const RequestIDHeaderName = "X-Request-ID"

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeaderName)
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Set(RequestIDHeaderName, id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestIDKey, id)))
	})
}
