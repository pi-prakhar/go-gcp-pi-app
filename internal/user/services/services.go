package services

import "github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"

type UserService struct {
	repository repository.UserRepository
}
