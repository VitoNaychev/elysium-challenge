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

func (u *UserService) CreateUser(user *domain.User) error {
	u.repo.Create(user)

	jwt, _ := crypto.GenerateJWT(u.jwtConfig, user.ID)
	user.JWTs = []string{jwt}

	u.repo.Update(user)

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
		return -1, NewUserServiceError("couldn't verigy JWT", err)
	}

	_, err = u.repo.GetByID(id)
	if err != nil {
		return -1, ErrUserNotFound
	}

	return id, nil
}
