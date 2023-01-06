package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/app"
	"github.com/orangeseeds/go-api/pkg/storage"
)

// connectDB initializes the driver with username and password authentication,
// uri parameter points to the host address of the neo4j database. It also check the connectivity with
// the host and returns neo4j.Driver, returns error if connection cannot be extablished.
func connectDB(uri string, username string, password string) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		uri, neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		return nil, err
	}
	err = driver.VerifyConnectivity(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB successfully as " + username)
	return driver, nil
}

// listenAndServe serves the application on port :addr
// and exits the application if it encounters any error.
func listenAndServe(server *app.Server, addr string) error {
	log.Printf("serving on port :%v\n", addr)
	err := server.Run(":" + addr)
	if err != nil {
		return err
	}
	return nil
}

// execFlagCommand services each flag given from the commandline.execFlagCommand.
// Returns true if only need to service the flag and not proceed any further,
// returns false if to proceed further after servicing the flag.
func execFlagCommand(config app.ServerConfig) bool {
	var ran = false
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "routes":
			printRouteList(config)
			ran = true
		case "migrate":
			runMigrations(config)
			ran = true

		}
	})
	return ran
}

// printRouteList prints the list of all the routes in the router.
func printRouteList(config app.ServerConfig) {
	server := app.NewServer(config, app.NewRouter(), nil)
	for _, item := range server.Routes().RouteList() {
		fmt.Println(item)
	}
}

// runMigration runs cypher query to set required value constrains in the neo4j server,
func runMigrations(config app.ServerConfig) {
	db, err := connectDB(config.Uri, config.Username, config.Password)
	if err != nil {
		log.Fatalln(err)
	}

	storage := storage.NewStorage(db, config.Username)
	err = storage.RunMigrations()
	if err != nil {
		panic(err)
	}
	log.Println("Migrations ran successfully")
}
