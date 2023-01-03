package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/helpers"
)

func (s *storage) CreateUser(request api.AuthUserRequest) error {
	session := s.neoDB.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})
	defer session.Close()

	hashPassword, _ := helpers.Hash(request.Password)

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
		"password_reset_token": helpers.RandomString(64),
	}

	_, err := session.Run(query, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindUserByEmail(email string) (*api.User, error) {
	session := s.neoDB.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
		MATCH (u:User {email: $email}) 
		RETURN 
			u.username as username, 
			u.password as password
		`
	params := map[string]any{
		"email": email,
	}

	res, err := session.Run(query, params)
	if err != nil {
		return nil, err
	}

	record, err := res.Single()
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

	if !helpers.PasswordsMatch(user.Password, password) {
		return nil, nil
	}
	return &api.User{
		Username: user.Username,
		Email:    user.Email,
	}, nil

}

func (s *storage) SetUserReadArticle(userEmail string, articleId string) (*api.Article, error) {
	session := s.neoDB.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})
	defer session.Close()

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

	res, err := session.Run(query, params)
	if err != nil {
		return nil, err
	}

	record, err := res.Single()
	if err != nil {
		return nil, err
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}

func (s *storage) SetUserBookmarksArticle(userEmail string, articleId string) (*api.Article, error) {
	session := s.neoDB.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})
	defer session.Close()

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

	res, err := session.Run(query, params)
	if err != nil {
		return nil, err
	}

	record, err := res.Single()
	if err != nil {
		return nil, err
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}
