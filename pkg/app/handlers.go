package app

import (
	"encoding/json"
	"net/http"

	"github.com/orangeseeds/go-api/pkg/api"
)

func (s *Server) ApiStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": "The API is running smoothly",
		})
	}
}

func (s *Server) CreateNewUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		u := api.NewUserRequest{
			Username: "usernameTest",
			Password: "passwordTest",
			Email:    "emailTest",
		}
		_ = s.userService.Create(u)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    u,
			"message": "new user created.",
		})
	}
}

func (s *Server) AllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := s.userService.All()
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": err,
			})
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    users,
			"message": "new user created.",
		})

	}
}