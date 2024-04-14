package service

import "errors"

type UserServiceError struct {
	msg string
	err error
}

func NewUserServiceError(msg string, err error) *UserServiceError {
	return &UserServiceError{
		msg: msg,
		err: err,
	}
}

func (u *UserServiceError) Error() string {
	return u.msg
}

func (u *UserServiceError) Unwrap() error {
	return u.err
}

var (
	ErrUserNotFound  = &UserServiceError{msg: "user doesn't exist"}
	ErrEmailNotFound = errors.New("user with this email doesn't exist")
	ErrWrongPassword = errors.New("wrong password for user with this email")
)
