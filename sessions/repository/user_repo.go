package repository

import "github.com/VitoNaychev/elysium-challenge/sessions/domain"

type UserRepo interface {
	Create(*domain.User) error
	Update(*domain.User) error
	GetByEmail(string) (domain.User, error)
	GetByID(int) (domain.User, error)
}
