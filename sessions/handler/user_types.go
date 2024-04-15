package handler

import "github.com/VitoNaychev/elysium-challenge/sessions/domain"

type ErrorResponse struct {
	Message string `json:"message"`
}

type SignUpResponse struct {
	JWT  JWTResponse  `json:"jwt"`
	User UserResponse `json:"user"`
}

func userToSignUpResponse(u domain.User) SignUpResponse {
	return SignUpResponse{
		JWT: JWTResponse{
			Token: u.JWTs[0],
		},
		User: UserResponse{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Password:  u.Password,
		},
	}
}

type JWTResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func signUpRequestToUser(r SignUpRequest) domain.User {
	return domain.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Password:  r.Password,
	}
}
