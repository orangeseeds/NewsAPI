package storage

import (
	"context"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
)

const (
	WriteMode = neo4j.AccessModeWrite
	ReadMode  = neo4j.AccessModeRead
)

type Storage interface {
	RunMigrations() error
	CreateUser(request api.AuthUserRequest) error
	FindUserByEmailAndPassword(email string, password string) (*api.User, error)
	SetUserReadArticle(string, string) (*api.Article, error)
	SetUserBookmarksArticle(string, string) (*api.Article, error)
}
type storage struct {
	ctx   context.Context
	neoDB neo4j.DriverWithContext
}

func NewStorage(neo neo4j.DriverWithContext, name string) Storage {
	ctx := context.WithValue(context.Background(), "dbName", name)
	return &storage{
		ctx:   ctx,
		neoDB: neo,
	}
}

func (s storage) RunCypher(cypher string, params map[string]any, mode neo4j.AccessMode) (neo4j.ResultWithContext, error) {
	dbName, ok := s.ctx.Value("dbName").(string)
	if !ok {
		return nil, errors.New("need a dbName in config")
	}

	session := s.neoDB.NewSession(s.ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: dbName,
	})
	defer session.Close(s.ctx)

	res, err := session.Run(s.ctx, cypher, params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *storage) RunMigrations() error {
	session := s.neoDB.NewSession(
		s.ctx,
		neo4j.SessionConfig{
			AccessMode:   neo4j.AccessModeWrite,
			DatabaseName: "neo4j",
		},
	)
	defer session.Close(s.ctx)

	// Syntax for defining constraints:
	// CREATE CONSTRAINT [constraint_name] [IF NOT EXISTS]
	// FOR (n:LabelName)
	// REQUIRE (n.propertyName_1, ..., n.propertyName_n) IS UNIQUE
	// [OPTIONS "{" option: value[, ...] "}"]

	// Constraints to set unique ids for each Nodel Labels
	constraints := []string{
		`CREATE CONSTRAINT unique_user_id IF NOT EXISTS
		FOR (x:User) REQUIRE x.user_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_article_id IF NOT EXISTS
		FOR (x:Article) REQUIRE x.article_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_author_id IF NOT EXISTS
		FOR (x:Author) REQUIRE x.author_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_source_id IF NOT EXISTS
		FOR (x:Source) REQUIRE x.source_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_topic_id IF NOT EXISTS
		FOR (x:Topic) REQUIRE x.topic_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_user_email IF NOT EXISTS
		FOR (x:User) REQUIRE x.email IS UNIQUE`,

		`CREATE CONSTRAINT unique_user_reset_token IF NOT EXISTS
		FOR (x:User) REQUIRE x.password_reset_token IS UNIQUE`,
	}

	for i := range constraints {
		_, err := session.Run(s.ctx, constraints[i], nil)
		if err != nil {
			return err
		}
	}
	return nil

}
