package middleware

import (
	"net/http"

	log "github.com/likearthian/go-logger"
	"github.com/likearthian/go-logger/level"
	"github.com/ua-parser/uap-go/uaparser"
)

func MakeHttpTransportLoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqid := r.Header.Get("X-Request-Id")

			ua := uaparser.NewFromSaved()
			cl := ua.Parse(r.Header.Get("User-Agent"))
			_ = level.Info(logger).Log(
				"event", "incoming request",
				"request-id", reqid,
				"uri", r.RequestURI,
				"method", r.Method,
				"headers", r.Header,
				"origin", r.Header.Get("X-Forwarded-For"),
				"protocol", r.Proto,
				"user-agent", cl.UserAgent.ToString(),
				"device", cl.Device.ToString(),
				"os", cl.Os.ToString(),
			)
			next.ServeHTTP(w, r)
		})
	}
}
