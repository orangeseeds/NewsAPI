package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/orangeseeds/NewsAPI/pkg/api"
	"github.com/orangeseeds/NewsAPI/pkg/util"
	"github.com/orangeseeds/NewsAPI/pkg/util/validator"
)

func (s *Server) handleApiStatus(w http.ResponseWriter, r *http.Request) {
	util.Respond(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"version": "1.0",
		},
		"message": "The API is running smoothly",
		"success": true,
	})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	u := api.CreateUserRequest{}
	ok := validateReqBody(w, r, &u)
	if !ok {
		return
	}

	userId, err := s.userService.Create(u)
	if err != nil {
		util.RespondErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := util.AuthToken(s.config.JwtSecret, userId, time.Minute*50)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
		return
	}
	u.Password = ""
	util.Respond(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     u,
		"message":  "new user created.",
		"jwtToken": token,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	u := api.LoginRequest{}
	ok := validateReqBody(w, r, &u)
	if !ok {
		return
	}

	user, err := s.userService.Login(u.Email, u.Password)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err.Error())
		return
	}
	token, err := util.AuthToken(s.config.JwtSecret, user.UserId, time.Minute*50)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
		return
	}

	user.UserId = ""
	util.Respond(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     user,
		"message":  "logged in as " + user.Username,
		"jwtToken": token,
	})
}

func (s *Server) handleReadArticle(w http.ResponseWriter, r *http.Request) {
	a := api.ReadArticleRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	article, err := s.userService.Read(userId, a.ArticleId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": fmt.Sprintf("added article to read list."),
	})
}

func (s *Server) handleBookmarkArticle(w http.ResponseWriter, r *http.Request) {
	a := api.BookmarkArticleRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	article, err := s.userService.Bookmark(userId, a.ArticleId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "added the article to bookmarked list.",
	})
}
func (s *Server) handleUnBookmarkArticle(w http.ResponseWriter, r *http.Request) {
	a := api.BookmarkArticleRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	article, err := s.userService.UnBookmark(userId, a.ArticleId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "removed the article from bookmarked list.",
	})
}

func (s *Server) handleFollowSource(w http.ResponseWriter, r *http.Request) {
	a := api.FollowSourceRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	source, err := s.userService.FollowSource(userId, a.SourceId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    source,
		"message": "added the source to followed list.",
	})
}

func (s *Server) handleUnFollowSource(w http.ResponseWriter, r *http.Request) {
	a := api.FollowSourceRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	source, err := s.userService.UnFollowSource(userId, a.SourceId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    source,
		"message": "removed the source from followed list.",
	})
}

func (s *Server) handleUnFollowCategory(w http.ResponseWriter, r *http.Request) {
	a := api.FollowCategoryRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	category, err := s.userService.UnFollowCategory(userId, a.CategoryId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    category,
		"message": "removed the category from followed list.",
	})
}

func (s *Server) handleFollowCategory(w http.ResponseWriter, r *http.Request) {
	a := api.FollowCategoryRequest{}
	ok := validateReqBody(w, r, &a)
	if !ok {
		return
	}

	userId := r.Context().Value("user").(string)
	category, err := s.userService.FollowCategory(userId, a.CategoryId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    category,
		"message": "added the category to followed list.",
	})
}

func validateReqBody[T any](w http.ResponseWriter, r *http.Request, v *T) bool {
	err := util.DecodeBody(r, v)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return false
	}

	errs := validator.ValidateStruct(*v)
	if len(errs) != 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
		return false
	}
	return true
}
