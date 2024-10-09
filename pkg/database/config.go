package database

import "time"

type Config struct {
	Type         string
	Host         string
	Port         int
	Database     string
	User         string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	SSLMode      string
	Options      map[string]string // Additional database-specific options
}

type GCPConfig struct {
	ProjectID       string
	Region          string
	InstanceName    string
	UseIAMAuth      bool
	CredentialsFile string
	UsePrivateIP    bool
}
