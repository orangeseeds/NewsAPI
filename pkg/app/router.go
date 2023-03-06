package app

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/orangeseeds/NewsAPI/pkg/util"
	"golang.org/x/exp/slices"
)

type Router struct {
	http.ServeMux
	activeGroup string
	routes      map[string]*Route
}

type Route struct {
	Methods  []string
	Handlers map[string]func(http.ResponseWriter, *http.Request)
	Path     string
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

func (r *Router) RouteList() []string {
	routes := []string{}
	for path := range r.routes {
		routes = append(routes, path)
	}
	sort.Strings(routes)
	return routes
}

func With(fn http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	if len(middlewares) == 0 {
		return http.HandlerFunc(fn)
	}
	for _, m := range middlewares {
		fn = m(fn)
	}
	return fn
}

func (r *Router) add(method, path string, h http.HandlerFunc) {
	if h == nil {
		panic("http: nil handler")
	}
	fullPath := r.activeGroup + path

	route, exists := r.routes[fullPath]
	if !exists {
		newRoute := &Route{
			Methods: []string{method},
			Path:    fullPath,
			Handlers: map[string]func(http.ResponseWriter, *http.Request){
				method: h,
			},
		}
		r.routes[fullPath] = newRoute
		r.Handle(fullPath, newRoute)
		return
	}

	if slices.Contains(route.Methods, method) {
		panic(fmt.Sprintf("multiple routes defined for %s %s", method, fullPath))
	}
	route.Methods = append(route.Methods, method)
	route.Handlers[method] = h
	return
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := route.Handlers[r.Method]
	if !ok {
		util.RespondHTTPErr(w, http.StatusMethodNotAllowed)
		return
	}
	handler(w, r)
}
