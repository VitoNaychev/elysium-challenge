package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/VitoNaychev/elysium-challenge/sessions/handler"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
)

type StubUserService struct {
	dummyUserID int
	dummyJWT    string
	dummyErr    error

	spyUser      domain.User
	spyLogoutJWT string
}

func (s *StubUserService) Create(user *domain.User) error {
	s.spyUser = *user

	user.ID = s.dummyUserID
	user.JWTs = []string{s.dummyJWT}

	return s.dummyErr
}

func (s *StubUserService) Login(email, password string) (string, error) {
	return s.dummyJWT, s.dummyErr
}

func (s *StubUserService) Logout(jwt string) error {
	s.spyLogoutJWT = jwt
	return s.dummyErr
}

func TestSignUpHandler(t *testing.T) {
	t.Run("creates new user", func(t *testing.T) {
		wantUser := domain.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}
		signUpRequest := handler.SignUpRequest{
			FirstName: wantUser.FirstName,
			LastName:  wantUser.LastName,
			Email:     wantUser.Email,
			Password:  wantUser.Password,
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(signUpRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/signup", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{}
		userHandler := handler.NewUserHandler(userService)
		userHandler.SignUp(response, request)

		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, userService.spyUser, wantUser)
	})

	t.Run("writes JWT and new user to response", func(t *testing.T) {
		dummyUserID := 10
		dummyJWT := "sampleJWT"

		wantResponse := handler.SignUpResponse{
			JWT: handler.JWTResponse{
				Token: dummyJWT,
			},
			User: handler.UserResponse{
				ID:        dummyUserID,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@example.com",
				Password:  "samplepassword",
			},
		}
		signUpRequest := handler.SignUpRequest{
			FirstName: wantResponse.User.FirstName,
			LastName:  wantResponse.User.LastName,
			Email:     wantResponse.User.Email,
			Password:  wantResponse.User.Password,
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(signUpRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/signup", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyUserID: dummyUserID,
			dummyJWT:    dummyJWT,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.SignUp(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse handler.SignUpResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("returns Bad Request on missing request body", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/user/signup", nil)
		response := httptest.NewRecorder()

		userService := &StubUserService{}
		userHandler := handler.NewUserHandler(userService)

		userHandler.SignUp(response, request)
		assert.Equal(t, response.Code, http.StatusBadRequest)

		var errorResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assert.Equal(t, errorResponse.Message, handler.ErrEmptyBody.Error())
	})

	t.Run("returns Internal Server Error on error from UserService", func(t *testing.T) {
		dummyError := errors.New("dummy error")

		signUpRequest := handler.SignUpRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Password:  "samplepassword",
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(signUpRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/signup", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyError,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.SignUp(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var errorResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assert.Equal(t, errorResponse.Message, dummyError.Error())
	})
}

func TestLoginHandler(t *testing.T) {
	t.Run("returns new JWT on correct email and password", func(t *testing.T) {
		wantResponse := handler.JWTResponse{
			Token: "sampleToken",
		}

		loginRequest := handler.LoginRequest{
			Email:    "johndoe@example.com",
			Password: "samplepassword",
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(loginRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/login", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyJWT: wantResponse.Token,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Login(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		var gotResponse handler.JWTResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse, wantResponse)
	})

	t.Run("returns Unauthorized on ErrEmailNotFound", func(t *testing.T) {
		dummyError := service.ErrEmailNotFound

		loginRequest := handler.LoginRequest{
			Email:    "johndoe@example.com",
			Password: "samplepassword",
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(loginRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/login", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyError,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Login(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse.Message, dummyError.Error())
	})

	t.Run("returns Unauthorized on ErrWrongPassowrd", func(t *testing.T) {
		dummyError := service.ErrWrongPassword

		loginRequest := handler.LoginRequest{
			Email:    "johndoe@example.com",
			Password: "samplepassword",
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(loginRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/login", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyError,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Login(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse.Message, dummyError.Error())
	})

	t.Run("returns Internal Server Error on unknown error from UserService", func(t *testing.T) {
		dummyError := errors.New("dummy error")

		loginRequest := handler.LoginRequest{
			Email:    "johndoe@example.com",
			Password: "samplepassword",
		}

		reqBody := bytes.NewBuffer([]byte{})
		json.NewEncoder(reqBody).Encode(loginRequest)

		request, _ := http.NewRequest(http.MethodPost, "/user/login", reqBody)
		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyError,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Login(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse.Message, dummyError.Error())
	})
}

func TestLogoutHandler(t *testing.T) {
	t.Run("calls UserService.Logout", func(t *testing.T) {
		wantJWT := "sampleToken"

		request, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		request.Header.Add("Token", wantJWT)

		response := httptest.NewRecorder()

		userService := &StubUserService{}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Logout(response, request)
		assert.Equal(t, response.Code, http.StatusOK)

		assert.Equal(t, userService.spyLogoutJWT, wantJWT)
	})

	t.Run("returns Bad Request on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		response := httptest.NewRecorder()

		userService := &StubUserService{}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Logout(response, request)
		assert.Equal(t, response.Code, http.StatusBadRequest)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse.Message, handler.ErrMissingToken.Error())
	})

	t.Run("returns Unathorized on invalid JWT", func(t *testing.T) {
		dummyJWT := "sampleToken"
		dummyErr := service.ErrInvalidJWT

		request, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		request.Header.Add("Token", dummyJWT)

		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyErr,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Logout(response, request)
		assert.Equal(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Not Found on ErrUserNotFound", func(t *testing.T) {
		dummyJWT := "sampleToken"
		dummyErr := service.ErrUserNotFound

		request, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		request.Header.Add("Token", dummyJWT)

		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyErr,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Logout(response, request)
		assert.Equal(t, response.Code, http.StatusNotFound)
	})

	t.Run("returns Internal Server Error on unknown error from UserService", func(t *testing.T) {
		dummyJWT := "sampleToken"
		dummyError := errors.New("dummy error")

		request, _ := http.NewRequest(http.MethodPost, "/user/logout", nil)
		request.Header.Add("Token", dummyJWT)

		response := httptest.NewRecorder()

		userService := &StubUserService{
			dummyErr: dummyError,
		}
		userHandler := handler.NewUserHandler(userService)

		userHandler.Logout(response, request)
		assert.Equal(t, response.Code, http.StatusInternalServerError)

		var gotResponse handler.ErrorResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		assert.Equal(t, gotResponse.Message, dummyError.Error())
	})
}
