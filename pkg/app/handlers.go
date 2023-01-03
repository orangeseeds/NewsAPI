package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/helpers"
	h "github.com/orangeseeds/go-api/pkg/helpers"
)

func (s *Server) handleApiStatus(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	h.Respond(w, http.StatusOK, map[string]any{
		"data": "The API is running smoothly",
		"user": user,
	})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := api.AuthUserRequest{}
	err := h.DecodeBody(r, &u)
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if u.Username == "" || u.Password == "" || u.Email == "" {
		h.RespondErr(w, http.StatusBadRequest, "username, password and email are required")
		return
	}

	err = s.userService.Create(u)
	if err != nil {
		neo4jError, ok := err.(*neo4j.Neo4jError)
		if ok && neo4jError.Title() == "ConstraintValidationFailed" {
			h.RespondErr(w, http.StatusBadRequest, errors.New("an account already exists with the given email address"))
			return
		}
		h.RespondHTTPErr(w, http.StatusInternalServerError)
		return
	}
	token, err := helpers.NewToken(s.config.JwtSecret, map[string]any{
		"user": u.Email,
		"exp":  time.Now().Add(time.Minute * 50).Unix(),
	})
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}
	u.Password = ""
	h.Respond(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     u,
		"message":  "new user created.",
		"jwtToken": token,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	u := api.AuthUserRequest{}
	err := h.DecodeBody(r, &u)
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if u.Password == "" || u.Email == "" {
		h.RespondErr(w, http.StatusBadRequest, "password and email are required")
		return
	}

	user, err := s.userService.Login(u.Email, u.Password)
	if err != nil {
		h.RespondErr(w, http.StatusBadRequest, errors.New("no account with email "+u.Email))
		return
	} else if user == nil {
		h.RespondErr(w, http.StatusBadRequest, errors.New("email or password didnot match"))
		return
	}
	token, err := helpers.NewToken(s.config.JwtSecret, map[string]any{
		"user": u.Email,
		"exp":  time.Now().Add(time.Minute * 50).Unix(),
	})
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	h.Respond(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     user,
		"message":  "logged in as " + user.Username,
		"jwtToken": token,
	})
}

func (s *Server) handleReadArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a := api.SaveArticleRequest{}
	err := h.DecodeBody(r, &a)
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if a.ArticleId == "" {
		h.RespondErr(w, http.StatusBadRequest, "article_id required")
		return
	}
	email := r.Context().Value("user").(string)
	article, err := s.userService.Read(email, a.ArticleId)
	if err != nil {
		h.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	h.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "read the article: " + article.ArticleId,
	})
}

func (s *Server) handleBookmarkArticle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	a := api.SaveArticleRequest{}
	err := h.DecodeBody(r, &a)
	if err != nil {
		h.RespondHTTPErr(w, http.StatusInternalServerError)
	}

	if a.ArticleId == "" {
		h.RespondErr(w, http.StatusBadRequest, "article_id required")
		return
	}
	email := r.Context().Value("user").(string)
	article, err := s.userService.Bookmark(email, a.ArticleId)
	if err != nil {
		h.RespondErr(w, http.StatusBadRequest, err)
		return
	}

	h.Respond(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    article,
		"message": "bookmarked the article: " + article.ArticleId,
	})
}