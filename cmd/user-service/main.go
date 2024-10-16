package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/config"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/handlers"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/router"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/services"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/database"
)

func main() {
	// time.Sleep(30 * time.Minute)
	err := godotenv.Load("../../env/.env.local")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var userConfig *models.Config = config.LoadUserConfig()

	gcpPostgresDatabase := database.NewGCPPostgresStrategy(userConfig.GCP)
	db, err := database.NewDatabase("gcp-postgres", userConfig.Database, gcpPostgresDatabase)
	if err != nil {
		log.Fatalf("Error : Failed to create database instance : %s", err.Error())
	}
	defer db.Close()

	repository := repository.GCPPostgresqlRepository{DB: db.GetDBConnectionPool()}
	userService := services.UserService{Repository: &repository}
	userHandler := handlers.UserHandler{Service: &userService}
	userRouter := router.NewRouter(&userHandler)

	err = userRouter.Engine.Run(userConfig.Service.Port)
	if err != nil {
		log.Fatalf("Error : Starting the server : %s", err.Error())
	}

}
