package service

import "errors"

var (
	ErrEmailNotFound = errors.New("user with this email doesn't exist")
	ErrWrongPassword = errors.New("wrong password for user with this email")
)
