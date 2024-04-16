package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
)

var (
	ErrEmptyBody    = errors.New("request body is empty")
	ErrMissingToken = errors.New("missing token in request")
)

type UserService interface {
	Create(*domain.User) error
	Login(string, string) (string, error)
	Logout(string) error
	Authenticate(string) (int, error)
}

type UserHTTPHandler struct {
	userService UserService

	http.Handler
}

func NewUserHTTPHandler(userService UserService) *UserHTTPHandler {
	userHandler := UserHTTPHandler{
		userService: userService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/user/signup", userHandler.SignUp)
	mux.HandleFunc("/user/login", userHandler.Login)
	mux.HandleFunc("/user/logout", userHandler.Logout)

	userHandler.Handler = mux

	return &userHandler
}

func (u *UserHTTPHandler) SignUp(w http.ResponseWriter, r *http.Request) {
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

func (u *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		writeErrorResponse(w, http.StatusBadRequest, ErrEmptyBody)
		return
	}

	var request LoginRequest
	json.NewDecoder(r.Body).Decode(&request)

	jwt, err := u.userService.Login(request.Email, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailNotFound) || errors.Is(err, service.ErrWrongPassword) {
			writeErrorResponse(w, http.StatusUnauthorized, err)
			return
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	response := JWTResponse{
		Token: jwt,
	}
	json.NewEncoder(w).Encode(response)
}

func (u *UserHTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	jwt := r.Header.Get("Token")
	if jwt == "" {
		writeErrorResponse(w, http.StatusBadRequest, ErrMissingToken)
		return
	}

	err := u.userService.Logout(jwt)
	if err != nil {
		if errors.Is(err, service.ErrInvalidJWT) {
			writeErrorResponse(w, http.StatusUnauthorized, err)
			return
		} else if errors.Is(err, service.ErrUserNotFound) {
			writeErrorResponse(w, http.StatusNotFound, err)
			return
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
}
