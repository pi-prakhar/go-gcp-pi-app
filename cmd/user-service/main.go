package main

import (
	"log"
	"os"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/constants"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/handlers"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/router"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/services"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/database"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/utils"
)

func main() {

	var config Config
	var err error
	var loader utils.Loader[Config]

	// err = godotenv.Load("../../env/.env.local")
	// if err != nil {
	// 	fmt.Println("Error : loading .env file:", err)
	// }

	configFilePath := os.Getenv(constants.CONFIG_FILE_PATH)
	if configFilePath == "" {
		log.Fatalf("Error : Failed to find Config file in path in env")
	}

	loader, err = utils.NewConfigLoader[Config](configFilePath, constants.CONFIG_FILE_TYPE, true)
	if err != nil {
		log.Fatalf("Error : Failed to create config loader : %s", err.Error())
	}

	config, err = loader.Load()
	if err != nil {
		log.Fatalf("Error : Failed to load config : %s", err.Error())
	}

	// Create new Database strategy
	gcpPostgresDatabase := database.NewGCPPostgresStrategy(config.GCP)

	// Create a new database instance with the GCP PostgreSQL strategy
	db, err := database.NewDatabase("gcp-postgres", config.Database, gcpPostgresDatabase)
	if err != nil {
		log.Fatalf("Error : Failed to create database instance : %s", err.Error())
	}
	defer db.Close()

	// get repo instance
	repository := repository.GCPPostgresqlRepository{DB: db.GetDBConnectionPool()}

	// get services instance
	userService := services.UserService{Repository: &repository}

	// get handlers instance
	userHandler := handlers.UserHandler{Service: &userService}
	// get routers instance
	userRouter := router.NewRouter(&userHandler)

	// start the server
	err = userRouter.Engine.Run(config.Service.Port)
	if err != nil {
		log.Fatalf("Error : Starting the server : %s", err.Error())
	}

}
