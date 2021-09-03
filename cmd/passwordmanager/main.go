package main

import (
	"log"
	"net/http"
	"os"

	openapi "github.com/iskorotkov/passwordmanager/go"
	"github.com/iskorotkov/passwordmanager/internal/database/postgres"
	"github.com/iskorotkov/passwordmanager/internal/services"
)

func main() {
	log.Printf("Server started")

	connString := os.Getenv("CONN_STRING")
	db, err := postgres.New(connString)
	if err != nil {
		log.Fatal(err)
	}

	passwordService := services.NewPasswordService(db)
	passwordController := openapi.NewDefaultApiController(passwordService)

	router := openapi.NewRouter(passwordController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
