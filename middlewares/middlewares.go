package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
