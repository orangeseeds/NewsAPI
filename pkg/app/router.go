package app

import "net/http"

func (s *Server) Routes() *http.ServeMux {
	r := s.router

	r.HandleFunc("/", s.ApiStatus())
	r.HandleFunc("/test", s.CreateNewUser())
	r.HandleFunc("/all", s.AllUsers())
	return r
}