package config

import "time"

type Config struct {
	Service Service
	Auth    Auth
}

type Service struct {
	Mode         string
	Port         string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

type Auth struct {
	ServiceHost string
	Google      Google
	JWT         JWT
}

type Google struct {
	ClientId     string
	ClientSecret string
}

type JWT struct {
	Key string
}
