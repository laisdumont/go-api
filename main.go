package main

import (
	"go-api/db"
	"go-api/handler"
	"go-api/repository"
	"go-api/router"
	"go-api/service"
	"log"
	"net/http"
	"os"
)

func main() {
	db.Connect()

	repo := repository.NewUserRepository(db.DB)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)
	r := router.SetupRoutes(h)

	host := os.Getenv("APP_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	address := host + ":" + port

	log.Printf("Servidor ouvindo em http://%s\n", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
