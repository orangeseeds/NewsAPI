package storage

import (
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/util"
)

func (s *storage) CreateUser(request api.CreateUserRequest) (string, error) {
	hashPassword, _ := util.Hash(request.Password)

	query := `
		CREATE (u:User { 
			user_id: apoc.create.uuid(),
			username: $username,
			email: $email,
			password: $password,
			password_reset_token: $password_reset_token
		}) RETURN u.user_id as user_id`
	params := map[string]any{
		"username":             request.Username,
		"password":             hashPassword,
		"email":                request.Email,
		"password_reset_token": util.RandomString(64),
	}

	session, res, err := s.RunCypher(query, params, WriteMode)

	if err != nil {
		if GetType(err) == InvalidConstraint {
			return "", NewError(GetType(err), "invalid constraint email.")
		}
		return "", NewError(NeoError, err.Error())
	}
	defer session.Close(s.ctx)
	record, err := res.Single(s.ctx)
	if err != nil {
		return "", NewError(NeoError, err.Error())
	}

	userId, _ := record.Get("user_id")

	return userId.(string), nil
}

func (s *storage) FindUserByEmail(email string) (*api.User, error) {
	query := `
		MATCH (u:User {email: $email}) 
		RETURN 
			u.user_id as user_id,
			u.username as username, 
			u.password as password
		`
	params := map[string]any{
		"email": email,
	}

	session, res, err := s.RunCypher(query, params, ReadMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no account with the given email.")
	}

	hashedPassword, _ := record.Get("password")
	username, _ := record.Get("username")
	userId, _ := record.Get("user_id")

	return &api.User{
		UserId:   userId.(string),
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
		return nil, NewError(NoMatchFound, "email or password didn't match.")
	}
	return &api.User{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
	}, nil

}

func (s *storage) SetUserReadArticle(userId string, articleId string) (*api.Article, error) {
	query := `
		MATCH (a:Article {article_id: $article_id}) 
		MATCH (u:User {user_id: $user_id})
		MERGE (u)-[:Read]->(a)
		RETURN 
			a.article_id as article_id
		`
	params := map[string]any{
		"article_id": articleId,
		"user_id":    userId,
	}
	session, res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no records found with given info.")
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}

func (s *storage) SetUserBookmarksArticle(userId string, articleId string) (*api.Article, error) {
	query := `
		MATCH (a:Article {article_id: $article_id}) 
		MATCH (u:User {user_id: $user_id})
		MERGE (u)-[:Bookmarks]->(a)
		RETURN 
			a.article_id as article_id
		`
	params := map[string]any{
		"article_id": articleId,
		"user_id":    userId,
	}

	session, res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no records found with given info.")
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}
func (s *storage) DelUserBookmarksArticle(userId string, articleId string) (*api.Article, error) {
	query := `
		MATCH (a:Article {article_id: $article_id}) 
		MATCH (u:User {user_id: $user_id})
		MATCH (u)-[b:Bookmarks]->(a)
		DELETE b
		RETURN 
			a.article_id as article_id
		`
	params := map[string]any{
		"article_id": articleId,
		"user_id":    userId,
	}

	session, res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no records found with given info.")
	}

	article_id, _ := record.Get("article_id")
	return &api.Article{
		ArticleId: article_id.(string),
	}, nil
}

func (s *storage) SetUserFollowsSource(userId string, sourceId string) (*api.Source, error) {
	query := `
		MATCH (s:Source {source_id: $source_id}) 
		MATCH (u:User {user_id: $user_id})
		MERGE (u)-[:Follows]->(s)
		RETURN 
			s.source_id as source_id
		`
	params := map[string]any{
		"source_id": sourceId,
		"user_id":   userId,
	}

	session, res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no records found with given info.")
	}

	source_id, _ := record.Get("source_id")
	return &api.Source{
		SourceId: source_id.(string),
	}, nil
}

func (s *storage) DelUserFollowsSource(userId string, sourceId string) (*api.Source, error) {
	query := `
		MATCH (s:Source {source_id: $source_id}) 
		MATCH (u:User {user_id: $user_id})
		MATCH (u)-[f:Follows]->(s)
		DELETE f
		RETURN 
			s.source_id as source_id
		`
	params := map[string]any{
		"source_id": sourceId,
		"user_id":   userId,
	}

	session, res, err := s.RunCypher(query, params, WriteMode)
	if err != nil {
		return nil, err
	}
	defer session.Close(s.ctx)

	record, err := res.Single(s.ctx)
	if err != nil {
		return nil, NewError(NoMatchFound, "no records found with given info.")
	}

	source_id, _ := record.Get("source_id")
	return &api.Source{
		SourceId: source_id.(string),
	}, nil
}
