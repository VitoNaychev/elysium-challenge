package service_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/crypto"
	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
	"github.com/joho/godotenv"
)

type StubUserRepo struct {
	nextUserID    int
	spyCreateUser domain.User
	spyUpdateUser domain.User
}

func (s *StubUserRepo) Create(user *domain.User) error {
	user.ID = s.nextUserID
	s.spyCreateUser = *user

	return nil
}

func (s *StubUserRepo) Update(user *domain.User) error {
	s.spyUpdateUser = *user

	return nil
}

func TestCreateUser(t *testing.T) {
	godotenv.Load("../test.env")
	secret := []byte(os.Getenv("SECRET"))
	expiresAt, err := time.ParseDuration(os.Getenv("EXPIRES_AT"))
	assert.RequireNoError(t, err)

	jwtConfig := crypto.JWTConfig{
		Secret:    secret,
		ExpiresAt: expiresAt,
	}

	t.Run("stores new user", func(t *testing.T) {
		wantUserID := 10
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{nextUserID: wantUserID}
		service := service.NewUserService(jwtConfig, repo)

		// as we're interested in only obeserving whether service.CreateUser
		// creates a user using our repo, we use a dirtyUser as the argument
		// to prevent the function from overriding wantUser.JWTs and causing
		// the test to fail
		dirtyUser := wantUser
		err := service.CreateUser(&dirtyUser)
		assert.RequireNoError(t, err)

		wantUser.ID = wantUserID
		assert.Equal(t, repo.spyCreateUser, wantUser)
	})

	t.Run("generates JWT", func(t *testing.T) {
		wantUserID := 10
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{nextUserID: wantUserID}
		service := service.NewUserService(jwtConfig, repo)

		err := service.CreateUser(&wantUser)
		assert.RequireNoError(t, err)

		if len(wantUser.JWTs) != 1 {
			t.Fatalf("got %v JWTs but want one", len(wantUser.JWTs))
		}

		token, err := crypto.VerifyJWT(jwtConfig, wantUser.JWTs[0])
		assert.RequireNoError(t, err)

		gotSubject, err := token.Claims.GetSubject()
		assert.RequireNoError(t, err)

		gotUserID, err := strconv.Atoi(gotSubject)
		assert.RequireNoError(t, err)

		assert.Equal(t, gotUserID, wantUserID)
	})

	t.Run("updates JWT array", func(t *testing.T) {
		wantUserID := 10
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{nextUserID: wantUserID}
		service := service.NewUserService(jwtConfig, repo)

		err := service.CreateUser(&wantUser)
		assert.RequireNoError(t, err)

		assert.Equal(t, repo.spyUpdateUser, wantUser)
	})
}
