package main

import (
	"log"
	"net/http"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/config"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/handlers"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/router"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/services"
)

func main() {
	//time.Sleep(30 * time.Minute)
	// err := godotenv.Load("../../env/.env.local")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	var authConfig *models.Config = config.LoadAuthConfig()

	authService := services.NewGoogleAuthService(authConfig)
	authHandler := handlers.NewAuthHandler(authService)
	router := router.NewRouter(authHandler, authConfig)

	srv := http.Server{
		Addr:    authConfig.Service.Port,
		Handler: router.Mux,
	}

	log.Printf("Server started at %s", authConfig.Service.Port)
	log.Fatal(srv.ListenAndServe())
}
