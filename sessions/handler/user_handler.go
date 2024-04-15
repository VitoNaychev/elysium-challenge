package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
)

var (
	ErrEmptyBody = errors.New("request body is empty")
)

type IUserService interface {
	Create(*domain.User) error
}

type UserHandler struct {
	userService IUserService
}

func NewUserHandler(userService IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeErrorResponse(w, http.StatusBadRequest, ErrEmptyBody)

		return
	}

	var request SignUpRequest
	json.NewDecoder(r.Body).Decode(&request)

	user := signUpRequestToUser(request)

	err := u.userService.Create(&user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response := userToSignUpResponse(user)
	json.NewEncoder(w).Encode(response)
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
}
