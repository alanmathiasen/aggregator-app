package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alanmathiasen/aggregator-api/auth"
	"github.com/alanmathiasen/aggregator-api/db"
	"github.com/alanmathiasen/aggregator-api/router"
	"github.com/alanmathiasen/aggregator-api/services"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	Config Config
	Models services.Models
}

func (app *Application) Serve() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	fmt.Println("Server running on port", port);
	srv := &http.Server {
		Addr: fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}
	return srv.ListenAndServe()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	auth.InitStore()
	
	cfg := Config {
		Port : os.Getenv(("PORT")),
	}

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)

	
	defer dbConn.DB.Close()

	if err != nil {
		fmt.Println("Couldn't connect to DB")
	}

	app := &Application {
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	log.Fatal(app.Serve())
	
}