package methods

import (
	"fmt"
	"net/http"

	middleware "github.com/heinswanhtet/blogora-api/middlewares"
)

type HTTPMethod string

// Allowed HTTP methods
const (
	GET     HTTPMethod = http.MethodGet
	HEAD    HTTPMethod = http.MethodHead
	POST    HTTPMethod = http.MethodPost
	PUT     HTTPMethod = http.MethodPut
	PATCH   HTTPMethod = http.MethodPatch
	DELETE  HTTPMethod = http.MethodDelete
	CONNECT HTTPMethod = http.MethodConnect
	OPTIONS HTTPMethod = http.MethodOptions
	TRACE   HTTPMethod = http.MethodTrace
)

type CustomMux struct {
	*http.ServeMux
}

func NewCustomMux() *CustomMux {
	return &CustomMux{ServeMux: http.NewServeMux()}
}

func (mux *CustomMux) Attach(method HTTPMethod, path string, handler func(http.ResponseWriter, *http.Request), middlewares ...middleware.Middleware) {
	mux.Handle(fmt.Sprintf("%s %s", method, path), middleware.Chain(http.HandlerFunc(handler), middlewares...))
}

func (mux *CustomMux) Use(path string, router *CustomMux, middlewares ...middleware.Middleware) {
	var path1 string

	if len(path) > 0 {
		path1 = path[:len(path)-1]
	}

	if len(path1) > 0 {
		mux.Handle(path, middleware.Chain(
			http.StripPrefix(path1, router),
			middlewares...,
		))
	} else {
		mux.Handle(path, middleware.Chain(
			router,
			middlewares...,
		))
	}
}
