package service

import (
	"github.com/VitoNaychev/elysium-challenge/crypto"
	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/VitoNaychev/elysium-challenge/sessions/repository"
)

type UserService struct {
	jwtConfig crypto.JWTConfig
	repo      repository.UserRepo
}

func NewUserService(jwtConfig crypto.JWTConfig, repo repository.UserRepo) *UserService {
	return &UserService{
		jwtConfig: jwtConfig,
		repo:      repo,
	}
}

func (u *UserService) Create(user *domain.User) error {
	err := u.repo.Create(user)
	if err != nil {
		return NewUserServiceError("couldn't create user", err)
	}

	jwt, err := crypto.GenerateJWT(u.jwtConfig, user.ID)
	if err != nil {
		return NewUserServiceError("couldn't generate JWT", err)
	}

	user.JWTs = []string{jwt}

	err = u.repo.Update(user)
	if err != nil {
		return NewUserServiceError("couldn't update user", err)
	}

	return nil
}

func (u *UserService) Login(email, password string) (string, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return "", ErrEmailNotFound
	}

	if user.Password != password {
		return "", ErrWrongPassword
	}

	jwt, _ := crypto.GenerateJWT(u.jwtConfig, user.ID)
	user.JWTs = append(user.JWTs, jwt)

	u.repo.Update(&user)

	return jwt, nil
}

func (u *UserService) Authenticate(jwt string) (int, error) {
	id, err := crypto.VerifyJWT(u.jwtConfig, jwt)
	if err != nil {
		return -1, ErrInvalidJWT.Wrap(err)
	}

	_, err = u.repo.GetByID(id)
	if err != nil {
		return -1, ErrUserNotFound
	}

	return id, nil
}

func (u *UserService) Logout(jwt string) error {
	id, err := crypto.VerifyJWT(u.jwtConfig, jwt)
	if err != nil {
		return ErrInvalidJWT.Wrap(err)
	}

	user, err := u.repo.GetByID(id)
	if err != nil {
		return ErrUserNotFound
	}

	err = user.InvalidateJWT(jwt)
	if err != nil {
		return ErrInvalidJWT.Wrap(err)
	}

	err = u.repo.Update(&user)
	if err != nil {
		return NewUserServiceError("couldn't update user", err)
	}

	return nil
}
