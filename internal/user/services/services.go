package services

import "github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"

type UserService struct {
	Repository repository.UserRepository
}

// func NewUserService(repository repository.UserRepository) *UserService {
// 	return &UserService{
// 		repository: repository,
// 	}
// }
