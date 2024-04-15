package service

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

func (u *UserServiceError) Wrap(err error) error {
	u.err = err
	return u
}

var (
	ErrUserNotFound  = &UserServiceError{msg: "user doesn't exist"}
	ErrEmailNotFound = &UserServiceError{msg: "user with this email doesn't exist"}
	ErrWrongPassword = &UserServiceError{msg: "wrong password for user with this email"}

	ErrInvalidJWT = &UserServiceError{msg: "invalid JWT"}
)
