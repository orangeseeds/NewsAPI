package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/orangeseeds/NewsAPI/pkg/util"
)

type middlewares struct {
	config ServerConfig
}

type Middleware func(fn http.HandlerFunc) http.HandlerFunc

func NewMiddlewareContainer(config ServerConfig) *middlewares {
	return &middlewares{
		config: config,
	}
}

func (m *middlewares) Auth(fn http.HandlerFunc) http.HandlerFunc {
	withJwt := func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		secret := m.config.JwtSecret
		jwtClaims, err := util.GetJWTClaims(secret, token)
		if err != nil {
			util.RespondHTTPErr(w, http.StatusUnauthorized)
			return
		}
		user, ok := jwtClaims["user"]
		if !ok {
			util.RespondHTTPErr(w, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(context.Background(), "user", user)
		fn(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(withJwt)
}

func (m *middlewares) Logger(fn http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		fn(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s", method, uri, duration.Round(time.Millisecond))
	}
	return http.HandlerFunc(logFn)
}

func (m *middlewares) CORS(fn http.HandlerFunc) http.HandlerFunc {
	corsFn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
	return http.HandlerFunc(corsFn)
}
