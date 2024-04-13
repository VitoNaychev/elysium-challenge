package repository

import "errors"

var (
	ErrNotFound = errors.New("didn't find object in repository")
)
