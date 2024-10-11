package utils

import (
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/database"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (*database.Config, *database.GCPConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, nil, err
	}

	var dbConfig database.Config
	if err := v.UnmarshalKey("database", &dbConfig); err != nil {
		return nil, nil, err
	}

	var gcpConfig database.GCPConfig
	if err := v.UnmarshalKey("gcp", &gcpConfig); err != nil {
		return nil, nil, err
	}

	return &dbConfig, &gcpConfig, nil
}
