package route

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler

type Router interface {
	http.Handler
	Routes

	Use(middlewares ...MiddlewareFunc)

	Methods(method ...string) *Route

}

type Routes interface {
	Routes() []Route
	Middlewares() []MiddlewareFunc
}

// Route describes the details of a routing handler.
// Handlers map key is an HTTP method
type Route struct {
	SubRoutes Routes
	Handlers  map[string]http.Handler
	Pattern   string
}