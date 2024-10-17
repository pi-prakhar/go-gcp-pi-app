package errors

// Define a custom error type
type UserNotFoundError struct{}

var ErrUserNotFound = &UserNotFoundError{}

func (e *UserNotFoundError) Error() string {
	return "User not found"
}
