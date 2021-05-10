package middleware

import (
	"net/http"

	"github.com/rs/xid"
)

func HttpRequestIDInjectorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqid := r.Header.Get("X-Request-Id")

		if reqid == "" {
			reqid = xid.New().String()
		}

		r.Header.Set("X-Request-Id", reqid)

		// inject header
		// w.Header().Set("Access-Control-Allow-Headers", "Authorization,Origin,Accept,datetime,signature")
		// w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
