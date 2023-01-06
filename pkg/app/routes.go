package app

// TODO:
// ✔ => completed
// * => doesn't require Authentication
//
// [GET]*    status =>  returns information about the API. ✔
//
// [POST]*   login/ => returns JWTtoken for the current user credentials ✔
// [POST]*   register/ => creates new user with the provided creds and returns JWTtoken ✔
// [POST]    logout/ => invalidates all the valid JWT token, of the requesting user
//
// [GET]     profile => returns user profile information.
// [POST]    profile/update/ => updates profile information with given info
// [GET]     followList/ => returns list of Category|Authors|Sources, that User has followed
// [GET]     blockList/ => returns list of Category|Authors|Sources, that User has blocked
// [GET]     bookmarks/ => returns list of currently bookmarked articles
//
// [POST]    bookmark/ => adds article with given id to bookmarked list ✔
// [DELETE]  bookmark/ => removes article with given id from bookmark list✔
// [POST]    read/ => stores a record that the user has read the article with the given id✔
//
// [POST]    followSource/   => adds the source with given id to followed sources list✔
// [DELETE]  followSource/ => removes the the source with given id from followed sources list✔
//
// [POST]    followCategory/ => adds the category with given id to followed categorys list✔
// [DELETE]  followCategory/ => removes the the category with given id from followed categorys list✔
//
// [POST]    blockSource/ => adds the source with given id to blocked sources list
// [DELETE]  blockSource/ => removes the the source with given id from blocked sources list
//
// [POST]    blockAuthor/ => adds the author with given id to blocked authors list
// [DELETE]  blockAuthor/ => removes the the author with given id from blocked authors list
//
// [GET]*    headlines/ => returns list of most relevant articles
// [GET]*    articles/ => returns list of suggested articles for the user
//
// [GET]*    sources/ => returns list of most relevant sources
// [GET]*    source/?id= => returns list of articles published by the given source
// [GET]*    categories/ => returns list of most relevant categories
// [GET]*    category/?id= => returns list of articles of the given category
//
// [GET]*    search/ => returns list of matching results, based on the search_query[article,author,topic,source]
func (s *Server) Routes() *Router {
	router := s.router
	m := NewMiddlewareContainer(s.config)

	var (
		logger Middleware = m.Logger
		cors   Middleware = m.CORS
		auth   Middleware = m.Auth
	)

	router.GET("/status/", With(s.handleApiStatus, cors, logger))

	router.POST("/login/", With(s.handleLogin, cors, logger))
	router.POST("/register/", With(s.handleRegister, cors, logger))

	router.POST("/read/", With(s.handleReadArticle, cors, auth, logger))

	router.POST("/bookmark/", With(s.handleBookmarkArticle, cors, auth, logger))
	router.DELETE("/bookmark/", With(s.handleUnBookmarkArticle, cors, auth, logger))

	router.DELETE("/followSource/", With(s.handleUnFollowSource, cors, auth, logger))
	router.POST("/followSource/", With(s.handleFollowSource, cors, auth, logger))

	router.DELETE("/followCategory/", With(s.handleUnFollowCategory, cors, auth, logger))
	router.POST("/followCategory/", With(s.handleFollowCategory, cors, auth, logger))
	return router
}
