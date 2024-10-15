package services

import (
	"context"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"
)

type UserService struct {
	Repository repository.UserRepository
}

// func NewUserService(repository repository.UserRepository) *UserService {
// 	return &UserService{
// 		repository: repository,
// 	}
// }

func (s *UserService) CreateUser(ctx context.Context, user models.GoogleUser) error {
	return s.Repository.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error) {
	return s.Repository.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.GoogleUser, error) {
	return s.Repository.GetUsers(ctx)
}
