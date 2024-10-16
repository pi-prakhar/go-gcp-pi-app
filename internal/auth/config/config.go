package config

import (
	"log"
	"os"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/constants"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/utils"
)

func LoadAuthConfig() *models.Config {
	var err error
	var loader utils.Loader[models.Config]
	var authConfig models.Config

	configFilePath := os.Getenv(constants.AUTH_CONFIG_FILE_PATH)
	if configFilePath == "" {
		log.Fatalf("Error : Failed to find Config file in path in env")
	}

	loader, err = utils.NewConfigLoader[models.Config](configFilePath, constants.AUTH_CONFIG_FILE_TYPE, true)
	if err != nil {
		log.Fatalf("Error : Failed to create config loader : %s", err.Error())
	}

	authConfig, err = loader.Load()
	if err != nil {
		log.Fatalf("Error : Failed to load config : %s", err.Error())
	}
	return &authConfig
}
