package main

import (
	"log"
	"net/http"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/app"
	"github.com/orangeseeds/go-api/pkg/repository"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	neoDriver, err := connectDB()
	if err != nil {
		log.Fatalln(err)
	}
	// Storage
	storage := repository.NewStorage(neoDriver)
	// Services
	userService := api.NewUserService(storage)

	router := http.NewServeMux()
	server := app.NewServer(router, userService)
	err = server.Run(":8080")
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func connectDB() (*neo4j.Driver, error) {

	config := config{
		Uri:      "neo4j+ssc://776b369c.databases.neo4j.io",
		Password: "g41MQzUv5GlG1LlcH-9J-VXkkIaupqdfBNoTgesuTpo",
		Username: "neo4j",
	}

	driver, err := neo4j.NewDriver(config.Uri,
		neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		return nil, err
	}
	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to db successfully as " + config.Username)
	return &driver, nil
}

type config struct {
	Uri      string
	Password string
	Username string
}