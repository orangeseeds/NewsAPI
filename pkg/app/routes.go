package app

func (s *Server) Routes() *Router {
	router := s.router
	m := NewMiddleware(s.config)

	var (
		logger Middleware = m.Logger
		cors   Middleware = m.CORS
		auth   Middleware = m.Auth
	)
	// r.Group("/api")
	router.GET("/status/", With(s.handleApiStatus, cors, logger))

	router.POST("/login/", With(s.handleLogin, cors, logger))
	router.POST("/register/", With(s.handleRegister, cors, logger))

	router.POST("/read/", With(s.handleReadArticle, cors, auth, logger))
	router.POST("/bookmark/", With(s.handleBookmarkArticle, cors, auth, logger))
	return router
}
