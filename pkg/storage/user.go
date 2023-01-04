package storage

import (
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/util"
)

func (s *storage) CreateUser(request api.AuthUserRequest) error {
	hashPassword, _ := util.Hash(request.Password)

	query := `
		CREATE (u:User { 
			user_id: apoc.create.uuid(),
			username: $username,
			email: $email,
			password: $password,
			password_reset_token: $password_reset_token
		})`
	params := map[string]any{
		"username":             request.Username,
		"password":             hashPassword,
		"email":                request.Email,
		"password_reset_token": util.RandomString(64),
	}

	_, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindUserByEmail(email string) (*api.User, error) {
	query := `
		MATCH (u:User {email: $email}) 
		RETURN 
			u.username as username, 
			u.password as password
		`
	params := map[string]any{
		"email": email,
	}

	res, err := s.RunCypher(query, params, ReadMode)
	if err != nil {
		return nil, err
	}

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, err
	}

	hashedPassword, _ := record.Get("password")
	username, _ := record.Get("username")

	return &api.User{
		Username: username.(string),
		Email:    email,
		Password: hashedPassword.(string),
	}, nil
}

func (s *storage) FindUserByEmailAndPassword(email string, password string) (*api.User, error) {
	user, err := s.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !util.PasswordsMatch(user.Password, password) {
		return nil, nil
	}
	return &api.User{
		Username: user.Username,
		Email:    user.Email,
	}, nil

}

func (s *storage) SetUserReadArticle(userEmail string, articleId string) (*api.Article, error) {
	query := `
		MATCH (a:Article {article_id: $article_id}) 
		MATCH (u:User {email: $email})
		MERGE (u)-[:Read]->(a)
		RETURN 
			a.article_id as article_id
		`
	params := map[string]any{
		"article_id": articleId,
		"email":      userEmail,
	}

	res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, err
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}

func (s *storage) SetUserBookmarksArticle(userEmail string, articleId string) (*api.Article, error) {
	query := `
		MATCH (a:Article {article_id: $article_id}) 
		MATCH (u:User {email: $email})
		MERGE (u)-[:Bookmarks]->(a)
		RETURN 
			a.article_id as article_id
		`
	params := map[string]any{
		"article_id": articleId,
		"email":      userEmail,
	}

	res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, err
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}
