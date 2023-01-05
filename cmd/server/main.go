package main

import (
	"context"
	"flag"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/app"
	"github.com/orangeseeds/go-api/pkg/storage"
)

func main() {
	var (
		addr    = flag.String("addr", "8080", "port number")
		migrate = flag.Bool("migrate", false, "run migration")
	)
	flag.Parse()

	config := app.LoadConfig("./config.json")
	db, err := connectDB(config.Uri, config.Username, config.Password)
	if err != nil {
		log.Fatalln(err)
	}

	storage := storage.NewStorage(db, config.Username)
	router := app.NewRouter()
	userService := api.NewUserService(storage)

	runMigrations(storage, *migrate)

	server := app.NewServer(config, router, userService)
	listenAndServe(server, *addr)
}

// connectDB initializes the driver with username and password authentication,
// uri param points to the host address of the neo4j database. It also check the connectivity with
// the host and returns neo4j.Driver
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

// runMigration sets required value constrains in the neo4j server,
// run is set to true when "-migrate" flag is provided when running the application
func runMigrations(storage storage.Storage, run bool) {
	if !run {
		return
	}
	err := storage.RunMigrations()
	if err != nil {
		panic(err)
	}
	log.Fatalln("Migrations ran successfully")

}

// listenAndServe serves the application on port :addr
// and exits the application if it encounters any error
func listenAndServe(server *app.Server, addr string) {
	log.Printf("serving on port :%v\n", addr)
	err := server.Run(":" + addr)
	if err != nil {
		log.Fatalf("error running the server: %v", err)
	}

}
