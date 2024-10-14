package models

import "github.com/pi-prakhar/go-gcp-pi-app/pkg/database"

type Config struct {
	Service  Service
	Database database.Config
	GCP      database.GCPConfig
}

type Service struct {
	Mode string
	Port string
}
