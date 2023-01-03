package app

import (
	"net/http"

	"github.com/orangeseeds/go-api/pkg/helpers"
)

type Router struct {
	http.ServeMux
	activeGroup string
}

type Route struct {
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
	Path    string
}

func (r *Router) GET(path string, h http.HandlerFunc) {
	r.add("GET", path, h)
}

func (r *Router) POST(path string, h http.HandlerFunc) {
	r.add("POST", path, h)
}

func (r *Router) add(method, path string, h http.HandlerFunc) {
	if h == nil {
		panic("http: nil handler")
	}
	grouped := r.activeGroup + path
	route := &Route{
		Method:  method,
		Path:    path,
		Handler: h,
	}
	r.Handle(grouped, route)
}

func (r *Router) Group(group string) {
	r.activeGroup = group
}

func NewRouter() *Router {
	return &Router{
		activeGroup: "",
	}
}

func With(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	if len(middlewares) > 0 {
		for _, m := range middlewares {
			fn = m(fn)
		}
		return fn
	}
	return http.HandlerFunc(fn)
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if route.Method != r.Method {
		helpers.RespondHTTPErr(w, http.StatusMethodNotAllowed)
		return
	}
	route.Handler(w, r)
}
