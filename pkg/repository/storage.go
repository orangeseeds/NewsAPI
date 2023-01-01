package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
)

type Storage interface {
	CreateUser(request api.NewUserRequest) error
	FindUserByEmailAndPassword(email string, password string) (*api.User, error)
}
type storage struct {
	neoDB neo4j.Driver
}

func NewStorage(neo neo4j.Driver) Storage {
	return &storage{
		neoDB: neo,
	}
}