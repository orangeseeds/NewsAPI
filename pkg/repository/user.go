package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
)

func (s *storage) CreateUser(request api.NewUserRequest) error {
	session := s.neoDB.NewSession(neo4j.SessionConfig{
		AccessMode:   neo4j.AccessModeWrite,
		DatabaseName: "neo4j",
	})
	defer session.Close()

	query := `
		CREATE (u:User { 
			user_id: $user_id,
			username: $username,
			email: $email,
			password: $password
		})`
	params := map[string]any{
		"user_id":  "",
		"username": request.Username,
		"password": request.Password,
		"email":    request.Email,
	}

	_, err := session.Run(query, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindUserByEmailAndPassword(email string, password string) (*api.User, error) {

	session := s.neoDB.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	res, err := session.Run(`
		MATCH 
			(u:User {email: $email}) 
		RETURN 
			u.username as username, 
			u.password as password,
			u.email as email
		`, nil)
	if err != nil {
		return nil, err
	}

	record, err := res.Single()
	if err != nil {
		return nil, err
	}

	username, _ := record.Get("username")
	password, _ := record.Get("password")
	email, _ := record.Get("email")

	return &api.User{
		Username: username.(string),
		Email:    email,
	}, nil

}