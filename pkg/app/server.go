package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/orangeseeds/NewsAPI/pkg/api"
	"github.com/orangeseeds/NewsAPI/pkg/util/validator"
)

type Server struct {
	http.Server
	config      ServerConfig
	router      *Router
	userService api.UserService
}

type ServerConfig struct {
	Uri       string `json:"NEO4J_URI" validate:"required"`
	Username  string `json:"NEO4J_USERNAME" validate:"required"`
	Password  string `json:"NEO4J_PASSWORD" validate:"required"`
	Port      string `json:"APP_PORT" validate:"required"`
	JwtSecret string `json:"JWT_SECRET" validate:"required"`
}

func NewServer(config ServerConfig, router *Router, userService api.UserService) *Server {
	return &Server{
		config:      config,
		router:      router,
		userService: userService,
	}
}

func (s *Server) Run(addr string) error {
	s.Routes()

	s.Addr = addr
	s.Handler = s.router
	return s.ListenAndServe()
}

func LoadConfig(path string) (ServerConfig, error) {
	var config ServerConfig
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	errs := validator.ValidateStruct(config)
	if len(errs) > 0 {
		return config, fmt.Errorf("%v", errs)
	}
	return config, nil
}
