package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/util"
)

func (s *Server) handleApiStatus(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	util.Respond(w, http.StatusOK, map[string]any{
		"data": "The API is running smoothly",
		"user": user,
	})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := api.AuthUserRequest{}
	err := util.DecodeBody(r, &u)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if u.Username == "" || u.Password == "" || u.Email == "" {
		util.RespondErr(w, http.StatusBadRequest, "username, password and email are required")
		return
	}

	err = s.userService.Create(u)
	if err != nil {
		neo4jError, ok := err.(*neo4j.Neo4jError)
		if ok && neo4jError.Title() == "ConstraintValidationFailed" {
			util.RespondErr(w, http.StatusBadRequest, errors.New("an account already exists with the given email address"))
			return
		}
		util.RespondHTTPErr(w, http.StatusInternalServerError)
		return
	}
	token, err := util.NewToken(s.config.JwtSecret, map[string]any{
		"user": u.Email,
		"exp":  time.Now().Add(time.Minute * 50).Unix(),
	})
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
	u := api.AuthUserRequest{}
	err := util.DecodeBody(r, &u)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if u.Password == "" || u.Email == "" {
		util.RespondErr(w, http.StatusBadRequest, "password and email are required")
		return
	}

	user, err := s.userService.Login(u.Email, u.Password)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, errors.New("no account with email "+u.Email))
		return
	} else if user == nil {
		util.RespondErr(w, http.StatusBadRequest, errors.New("email or password didnot match"))
		return
	}
	token, err := util.NewToken(s.config.JwtSecret, map[string]any{
		"user": u.Email,
		"exp":  time.Now().Add(time.Minute * 50).Unix(),
	})
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     user,
		"message":  "logged in as " + user.Username,
		"jwtToken": token,
	})
}

func (s *Server) handleReadArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a := api.SaveArticleRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if a.ArticleId == "" {
		util.RespondErr(w, http.StatusBadRequest, "article_id required")
		return
	}
	email := r.Context().Value("user").(string)
	article, err := s.userService.Read(email, a.ArticleId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "read the article: " + article.ArticleId,
	})
}

func (s *Server) handleBookmarkArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a := api.SaveArticleRequest{}
	err := util.DecodeBody(r, &a)
	if err != nil {
		util.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if a.ArticleId == "" {
		util.RespondErr(w, http.StatusBadRequest, "article_id required")
		return
	}
	email := r.Context().Value("user").(string)
	article, err := s.userService.Bookmark(email, a.ArticleId)
	if err != nil {
		util.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	util.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "bookmarked the article: " + article.ArticleId,
	})
}