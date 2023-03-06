package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/orangeseeds/NewsAPI/pkg/api"
	"github.com/orangeseeds/NewsAPI/pkg/app"
	"github.com/orangeseeds/NewsAPI/pkg/storage"
)

var (
	addr      = flag.String("addr", "8080", "port number")
	migrate   = flag.Bool("migrate", false, "run migration")
	routeList = flag.Bool("routes", false, "display route list")
)

func main() {
	// test comment
	flag.Parse()
	config, err := app.LoadConfig("./config.json")
	if err != nil {
		panic("error loading config" + err.Error())
	}
	if execFlagCommand(config) {
		return
	}
	db, err := connectDB(config.Uri, config.Username, config.Password)
	if err != nil {
		log.Fatalln(err)
	}

	storage := storage.NewStorage(db, config.Username)
	router := app.NewRouter()
	userService := api.NewUserService(storage)

	server := app.NewServer(config, router, userService)

	// Shutting down the server gracefully.
	errChan := make(chan error, 1)
	signalChan := make(chan os.Signal, 1)
	go func() {
		errChan <- listenAndServe(server, *addr)
	}()

	signal.Notify(signalChan, os.Interrupt)

	select {
	case err := <-errChan:
		log.Printf("failed to serve: %s\n", err.Error())
	case <-signalChan:
		log.Println("exiting the server...")
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer func() {
		err = db.Close(ctx)
		if err != nil {
			log.Printf("error closing db connection: %s\n", err.Error())
		}
		cancel()
	}()
	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("error shutting down the server: %s\n", err.Error())
	}
}
