package main

import (
	"crabigateur-api/pkg/api"
	"crabigateur-api/pkg/app"
	"crabigateur-api/pkg/repository"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\\n", err)
		os.Exit(1)
	}
}

func run() error {
	var connectionString string

	flag.StringVar(&connectionString, "dsn", "host=localhost port=5555 user=crabi password=gateur dbname=crabigateur sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	db, err := setupDatabase(connectionString)
	if err != nil {
		return err
	}

	storage := repository.NewStorage(db)

	router := gin.Default()
	router.Use(cors.Default())
	
	userService := api.NewUserService(storage)

	server := app.NewServer(router, userService)

	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}

func setupDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// ping the DB to ensure that it is connected
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}