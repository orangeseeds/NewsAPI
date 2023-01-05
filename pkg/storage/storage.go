package storage

import (
	"context"
	"errors"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
)

const (
	WriteMode = neo4j.AccessModeWrite
	ReadMode  = neo4j.AccessModeRead
)

type Storage interface {
	RunMigrations() error

	// UserService Methods
	CreateUser(api.CreateUserRequest) (string, error)
	FindUserByEmailAndPassword(string, string) (*api.User, error)
	SetUserReadArticle(string, string) (*api.Article, error)
	SetUserBookmarksArticle(string, string) (*api.Article, error)
	DelUserBookmarksArticle(string, string) (*api.Article, error)
	SetUserFollowsSource(string, string) (*api.Source, error)
	DelUserFollowsSource(string, string) (*api.Source, error)
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

func (s storage) RunCypher(cypher string, params map[string]any, mode neo4j.AccessMode) (neo4j.SessionWithContext, neo4j.ResultWithContext, error) {
	dbName, ok := s.ctx.Value("dbName").(string)
	if !ok {
		return nil, nil, errors.New("need a dbName in config")
	}

	session := s.neoDB.NewSession(s.ctx, neo4j.SessionConfig{
		AccessMode:   mode,
		DatabaseName: dbName,
	})

	res, err := session.Run(s.ctx, cypher, params)
	if err != nil {
		defer session.Close(s.ctx)
		neo4jError, ok := err.(*neo4j.Neo4jError)
		if ok && neo4jError.Title() == "ConstraintValidationFailed" {
			return nil, nil, NewError(InvalidConstraint, err.Error())
		}
		return nil, nil, NewError(NeoError, err.Error())
	}

	return session, res, nil
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

		`CREATE CONSTRAINT unique_creator_id IF NOT EXISTS
		FOR (x:Creator) REQUIRE x.creator_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_source_id IF NOT EXISTS
		FOR (x:Source) REQUIRE x.source_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_category_id IF NOT EXISTS
		FOR (x:Category) REQUIRE x.category_id IS UNIQUE`,

		`CREATE CONSTRAINT unique_user_email IF NOT EXISTS
		FOR (x:User) REQUIRE x.email IS UNIQUE`,

		`CREATE CONSTRAINT unique_user_reset_token IF NOT EXISTS
		FOR (x:User) REQUIRE x.password_reset_token IS UNIQUE`,
	}

	for _, item := range constraints {
		_, err := session.Run(s.ctx, item, nil)
		if err != nil {
			return err
		}
		log.Printf("cypher: %s ,successfully completed.\n", item)
	}
	return nil

}
