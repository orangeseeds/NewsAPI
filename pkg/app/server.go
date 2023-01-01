package app

import (
	"net/http"

	"github.com/orangeseeds/go-api/pkg/api"
)

type Server struct {
	router      *http.ServeMux
	userService api.UserService
}

func NewServer(router *http.ServeMux, userService api.UserService) *Server {
	return &Server{
		router:      router,
		userService: userService,
	}
}

func (s *Server) Run(addr string) error {
	r := s.Routes()
	return http.ListenAndServe(addr, r)
}
