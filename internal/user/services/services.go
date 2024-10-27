package services

import (
	"context"
	"fmt"
	"log"

	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/repository"
	errors "github.com/pi-prakhar/go-gcp-pi-app/pkg/error"
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
	// check if user does not already exists in database
	oldUser, err := s.Repository.GetUserByEmail(ctx, user.Email)
	fmt.Println("yes1")
	if err != errors.ErrUserNotFound && err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("yes")
	if oldUser != nil {

		return &errors.UserAlreadyExistsError{Email: user.Email}
	}

	return s.Repository.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.GoogleUser, error) {
	return s.Repository.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUsers(ctx context.Context) ([]*models.GoogleUser, error) {
	return s.Repository.GetUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.GoogleUser) (*models.GoogleUser, error) {
	var err error
	var currUser *models.GoogleUser

	log.Println("UPDATE : user : ", user)
	currUser, err = s.Repository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	log.Println("UPDATE : currUser : ", currUser)

	mergeGoogleUser(currUser, user)

	err = s.Repository.UpdateUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, email string) error {
	return s.Repository.DeleteUserByEmail(ctx, &email)
}

func mergeGoogleUser(currUser *models.GoogleUser, updatedUser *models.GoogleUser) {
	updatedUser.ID = currUser.ID
	// TODO : FIX the issue with verified email
	if updatedUser.VerifiedMail != currUser.VerifiedMail {
		updatedUser.VerifiedMail = currUser.VerifiedMail
	}
	if updatedUser.Name == "" {
		updatedUser.Name = currUser.Name
	}
	if updatedUser.GivenName == "" {
		updatedUser.GivenName = currUser.GivenName
	}
	if updatedUser.FamilyName == "" {
		updatedUser.FamilyName = currUser.FamilyName
	}
	if updatedUser.Picture == "" {
		updatedUser.Picture = currUser.Picture
	}
	if updatedUser.Locale == "" {
		updatedUser.Locale = currUser.Locale
	}
}
