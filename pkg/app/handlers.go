package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/util"
	"github.com/orangeseeds/go-api/pkg/util/validator"
)

func (s *Server) handleApiStatus(w http.ResponseWriter, r *http.Request) {
	util.Respond(w, http.StatusOK, map[string]any{
		"data": "The API is running smoothly",
	})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := api.CreateUserRequest{}
	err := util.DecodeBody(r, &u)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(u)
	if len(errs) != 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	u := api.LoginRequest{}
	err := util.DecodeBody(r, &u)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(u)
	if len(errs) > 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	a := api.ReadArticleRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(a)
	if len(errs) != 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	a := api.BookmarkArticleRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(a)
	if len(errs) > 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	a := api.BookmarkArticleRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(a)
	if len(errs) > 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	a := api.FollowSourceRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(a)
	if len(errs) > 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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
	defer r.Body.Close()
	a := api.FollowSourceRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	errs := validator.ValidateStruct(a)
	if len(errs) > 0 {
		util.RespondErr(w, http.StatusBadRequest, errs)
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