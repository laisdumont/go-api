package main

import (
	"go-api/db"
	"go-api/handler"
	"go-api/repository"
	"go-api/router"
	"go-api/service"
	"log"
	"net/http"
)

func main() {
	db.Connect()

	repo := repository.NewUserRepository(db.DB)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)
	r := router.SetupRoutes(h)

	log.Println("Servidor ouvindo em http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
