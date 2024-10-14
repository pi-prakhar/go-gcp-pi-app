package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/config"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/constants"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/handlers"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/router"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/services"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/utils"
)

func main() {
	//time.Sleep(30 * time.Minute)
	var authConfig config.Config
	var err error
	var loader utils.Loader[config.Config]

	configFilePath := os.Getenv(constants.AUTH_CONFIG_FILE_PATH)
	if configFilePath == "" {
		log.Fatalf("Error : Failed to find Config file in path in env")
	}

	loader, err = utils.NewConfigLoader[config.Config](configFilePath, constants.AUTH_CONFIG_FILE_TYPE, true)
	if err != nil {
		log.Fatalf("Error : Failed to create config loader : %s", err.Error())
	}

	authConfig, err = loader.Load()
	if err != nil {
		log.Fatalf("Error : Failed to load config : %s", err.Error())
	}

	authService := services.NewGoogleAuthService(&config.Config{})
	authHandler := handlers.NewAuthHandler(authService)
	router := router.NewRouter(authHandler)

	srv := http.Server{
		Addr:    authConfig.Service.Port,
		Handler: router.Mux,
	}

	log.Printf("Server started at %s", authConfig.Service.Port)
	log.Fatal(srv.ListenAndServe())
}
