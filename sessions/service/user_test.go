package service_test

import (
	"testing"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/crypto"
	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/VitoNaychev/elysium-challenge/sessions/repository"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
	"github.com/joho/godotenv"
)

type StubUserRepo struct {
	users []domain.User

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

func (s *StubUserRepo) GetByEmail(email string) (domain.User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}

	return domain.User{}, repository.ErrNotFound
}

func (s *StubUserRepo) GetByID(id int) (domain.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return domain.User{}, repository.ErrNotFound
}

func TestCreateUser(t *testing.T) {
	godotenv.Load("../test.env")

	jwtConfig, err := crypto.InitJWTConfigFromEnv()
	assert.RequireNoError(t, err)

	t.Run("stores new user", func(t *testing.T) {
		wantUserID := 10
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{nextUserID: wantUserID}
		userService := service.NewUserService(jwtConfig, repo)

		dirtyUser := wantUser
		err := userService.Create(&dirtyUser)
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
		userService := service.NewUserService(jwtConfig, repo)

		dirtyUser := wantUser
		err := userService.Create(&dirtyUser)
		assert.RequireNoError(t, err)

		if len(dirtyUser.JWTs) != 1 {
			t.Fatalf("got %v JWTs but want one", len(dirtyUser.JWTs))
		}

		assertValidJWT(t, jwtConfig, dirtyUser.JWTs[0], wantUserID)
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
		userService := service.NewUserService(jwtConfig, repo)

		dirtyUser := wantUser
		err := userService.Create(&dirtyUser)
		assert.RequireNoError(t, err)

		// check that the JWT was persisted in the repository
		if len(repo.spyUpdateUser.JWTs) != 1 {
			t.Fatalf("got %v JWTs but want one", len(repo.spyUpdateUser.JWTs))
		}

		// check that all other data is unchanged
		wantUser.ID = wantUserID
		wantUser.JWTs = repo.spyUpdateUser.JWTs
		assert.Equal(t, repo.spyUpdateUser, wantUser)
	})
}

func TestLoginUser(t *testing.T) {
	godotenv.Load("../test.env")

	jwtConfig, err := crypto.InitJWTConfigFromEnv()
	assert.RequireNoError(t, err)

	t.Run("return ErrEmailNotFound on user with no such email", func(t *testing.T) {
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		_, err := userService.Login("missingemail@example.com", wantUser.Password)
		assert.Equal(t, err, (error)(service.ErrEmailNotFound))
	})

	t.Run("return ErrWrongPassword on wrong password", func(t *testing.T) {
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		_, err := userService.Login(wantUser.Email, "wrongpassword")
		assert.Equal(t, err, (error)(service.ErrWrongPassword))
	})

	t.Run("generates JWT", func(t *testing.T) {
		wantUser := domain.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		jwt, err := userService.Login(wantUser.Email, wantUser.Password)
		assert.RequireNoError(t, err)

		assertValidJWT(t, jwtConfig, jwt, wantUser.ID)
	})

	t.Run("updates JWT array", func(t *testing.T) {
		wantUser := domain.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
			JWTs:      []string{"testJWT"},
		}

		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		jwt, err := userService.Login(wantUser.Email, wantUser.Password)
		assert.RequireNoError(t, err)

		if len(repo.spyUpdateUser.JWTs) != 2 {
			t.Fatalf("didn't update JWT array before update call")
		}

		assert.Equal(t, repo.spyUpdateUser.JWTs[1], jwt)
	})
}

func TestAuthenticate(t *testing.T) {
	godotenv.Load("../test.env")

	jwtConfig, err := crypto.InitJWTConfigFromEnv()
	assert.RequireNoError(t, err)

	t.Run("returns UserServiceError on invalid JWT", func(t *testing.T) {
		invalidJWT := "invalidJWT"

		repo := &StubUserRepo{}
		userService := service.NewUserService(jwtConfig, repo)

		_, err := userService.Authenticate(invalidJWT)
		assert.ErrorType[*service.UserServiceError](t, err)
	})

	t.Run("return ErrUserNotFound on missing user", func(t *testing.T) {
		unknownUserID := 15

		repo := &StubUserRepo{}
		userService := service.NewUserService(jwtConfig, repo)

		jwt, err := crypto.GenerateJWT(jwtConfig, unknownUserID)
		assert.RequireNoError(t, err)

		_, err = userService.Authenticate(jwt)
		assert.Equal(t, err, (error)(service.ErrUserNotFound))
	})

	t.Run("returns user ID on valid JWT", func(t *testing.T) {
		wantUser := domain.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		jwt, err := crypto.GenerateJWT(jwtConfig, wantUser.ID)
		assert.RequireNoError(t, err)

		wantUser.JWTs = []string{jwt}
		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		gotUserID, err := userService.Authenticate(jwt)
		assert.RequireNoError(t, err)

		assert.Equal(t, gotUserID, wantUser.ID)
	})
}

func TestLogout(t *testing.T) {
	godotenv.Load("../test.env")

	jwtConfig, err := crypto.InitJWTConfigFromEnv()
	assert.RequireNoError(t, err)

	t.Run("returns UserServiceError on invalid JWT", func(t *testing.T) {
		invalidJWT := "invalidJWT"

		repo := &StubUserRepo{}
		userService := service.NewUserService(jwtConfig, repo)

		err := userService.Logout(invalidJWT)
		assert.ErrorType[*service.UserServiceError](t, err)
	})

	t.Run("return ErrUserNotFound on missing user", func(t *testing.T) {
		unknownUserID := 15

		repo := &StubUserRepo{}
		userService := service.NewUserService(jwtConfig, repo)

		jwt, err := crypto.GenerateJWT(jwtConfig, unknownUserID)
		assert.RequireNoError(t, err)

		err = userService.Logout(jwt)
		assert.Equal(t, err, (error)(service.ErrUserNotFound))
	})

	t.Run("returns UserServiceError on missing JWTs array entry", func(t *testing.T) {
		wantUser := domain.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		jwt, err := crypto.GenerateJWT(jwtConfig, wantUser.ID)
		assert.RequireNoError(t, err)

		wantUser.JWTs = []string{"sampleEntry"}
		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		err = userService.Logout(jwt)
		assert.ErrorType[*service.UserServiceError](t, err)
	})

	t.Run("removes JWTs array entry on valid JWT", func(t *testing.T) {
		wantUser := domain.User{
			ID:        10,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		jwt, err := crypto.GenerateJWT(jwtConfig, wantUser.ID)
		assert.RequireNoError(t, err)

		wantUser.JWTs = []string{"sampleEntry", jwt}
		repo := &StubUserRepo{
			users: []domain.User{wantUser},
		}
		userService := service.NewUserService(jwtConfig, repo)

		err = userService.Logout(jwt)
		assert.RequireNoError(t, err)

		wantUser.JWTs = []string{"sampleEntry"}
		assert.Equal(t, repo.spyUpdateUser, wantUser)
	})
}

func assertValidJWT(t testing.TB, jwtConfig crypto.JWTConfig, jwt string, wantUserID int) {
	t.Helper()

	gotUserID, err := crypto.VerifyJWT(jwtConfig, jwt)
	assert.RequireNoError(t, err)

	assert.Equal(t, gotUserID, wantUserID)
}
