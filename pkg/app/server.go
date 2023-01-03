package app

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/orangeseeds/go-api/pkg/api"
)

type Server struct {
	config      ServerConfig
	router      *Router
	userService api.UserService
}

type ServerConfig struct {
	Uri       string `json:"NEO4J_URI"`
	Username  string `json:"NEO4J_USERNAME"`
	Password  string `json:"NEO4J_PASSWORD"`
	Port      int    `json:"APP_PORT"`
	JwtSecret string `json:"JWT_SECRET"`
}

func NewServer(config ServerConfig, router *Router, userService api.UserService) *Server {
	return &Server{
		config:      config,
		router:      router,
		userService: userService,
	}
}

func (s *Server) Run(addr string) error {
	if !s.validConfig() {
		return errors.New("server configurations not valid")
	}
	r := s.Routes()
	return http.ListenAndServe(addr, r)
}

func (s *Server) validConfig() bool {
	config := s.config
	if len(config.Password) == 0 || len(config.Uri) == 0 || len(config.Username) == 0 || len(config.JwtSecret) == 0 {
		return false
	}
	return true
}

func LoadConfig(path string) ServerConfig {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	config := ServerConfig{}
	if err = json.Unmarshal(file, &config); err != nil {
		panic(err)
	}
	return config
}
