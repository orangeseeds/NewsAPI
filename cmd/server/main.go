package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/orangeseeds/go-api/pkg/api"
	"github.com/orangeseeds/go-api/pkg/app"
	"github.com/orangeseeds/go-api/pkg/repository"
)

func main() {
	var (
		addr          = flag.String("addr", ":8080", "port number")
		runMigrations = flag.Bool("migrate", false, "run migration")
	)

	flag.Parse()
	// f, err := os.OpenFile("./testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	err := run(*runMigrations, *addr)
	if err != nil {
		log.Fatalln(err)
	}
}

func run(runMigrations bool, addr string) error {
	config := app.LoadConfig("./config.json")

	neoDriver, err := connectDB(config.Uri, config.Username, config.Password)
	if err != nil {
		log.Fatalln(err)
	}
	// Storage
	storage := repository.NewStorage(neoDriver)

	// run migrations
	if runMigrations {
		err = storage.RunMigrations()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Migrations ran successfully")
		return nil
	}

	// Services
	userService := api.NewUserService(storage)

	router := app.NewRouter()
	server := app.NewServer(config, router, userService)

	log.Printf("serving on :%v", config.Port)
	err = server.Run(":" + fmt.Sprintf("%v", config.Port))
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func connectDB(uri string, username string, password string) (neo4j.Driver, error) {

	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB successfully as " + username)
	return driver, nil
}
