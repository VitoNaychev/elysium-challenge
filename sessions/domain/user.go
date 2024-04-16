package domain

import (
	"errors"
)

var ErrJWTNotFound = errors.New("jwt not found is user JWTs array")

type User struct {
	ID        int
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
	Password  string
	JWTs      []string
}

func (u *User) InvalidateJWT(jwt string) error {
	jwtIdx := -1

	for i := range u.JWTs {
		if u.JWTs[i] == jwt {
			jwtIdx = i
			break
		}
	}

	if jwtIdx == -1 {
		return ErrJWTNotFound
	}

	u.JWTs = append(u.JWTs[:jwtIdx], u.JWTs[jwtIdx+1:]...)
	return nil
}
