package app

import (
	"net/http"

	"github.com/orangeseeds/go-api/pkg/util"
)

type Router struct {
	http.ServeMux
	activeGroup string
	routes      map[string]*Route
}

type Route struct {
	Methods []string
	Handler map[string]func(http.ResponseWriter, *http.Request)
	Path    string
}

func NewRouter() *Router {
	return &Router{
		activeGroup: "",
		routes:      map[string]*Route{},
	}
}

func (r *Router) GET(path string, h http.HandlerFunc) {
	r.add("GET", path, h)
}

func (r *Router) POST(path string, h http.HandlerFunc) {
	r.add("POST", path, h)
}

func (r *Router) DELETE(path string, h http.HandlerFunc) {
	r.add("DELETE", path, h)
}

func (r *Router) Group(group string) {
	r.activeGroup = group
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

func (r *Router) add(method, path string, h http.HandlerFunc) {
	if h == nil {
		panic("http: nil handler")
	}
	grouped := r.activeGroup + path

	route, exists := r.routes[path]
	if !exists {
		newRoute := &Route{
			Methods: []string{method},
			Path:    path,
			Handler: map[string]func(http.ResponseWriter, *http.Request){
				method: h,
			},
		}
		r.routes[grouped] = newRoute
		r.Handle(grouped, newRoute)
		return
	}
	route.Methods = append(route.Methods, method)
	route.Handler[method] = h
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := route.Handler[r.Method]
	if !ok {
		util.RespondHTTPErr(w, http.StatusMethodNotAllowed)
		return
	}
	handler(w, r)
}
