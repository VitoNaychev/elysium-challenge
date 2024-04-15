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
)

type StubUserService struct {
	dummyUserID int
	dummyJWT    string
	dummyErr    error

	spyUser domain.User
}

func (s *StubUserService) Create(user *domain.User) error {
	s.spyUser = *user

	user.ID = s.dummyUserID
	user.JWTs = []string{s.dummyJWT}

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
