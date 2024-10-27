package errors

import "fmt"

// Define a custom error type
type UserNotFoundError struct{}

type UserAlreadyExistsError struct {
	Email string
}

var ErrUserNotFound = &UserNotFoundError{}

var ErrUserAlreadyExist = &UserAlreadyExistsError{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("User with email %s already exists", e.Email)
}
