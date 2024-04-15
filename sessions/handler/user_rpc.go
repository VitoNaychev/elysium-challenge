package handler

import (
	"context"

	"github.com/VitoNaychev/elysium-challenge/rpc/sessions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRPCHandler struct {
	sessions.UnimplementedUserServer

	userService IUserService
}

func NewUserRPCHandler(userService IUserService) *UserRPCHandler {
	return &UserRPCHandler{
		userService: userService,
	}
}

func (u *UserRPCHandler) Authenticate(ctx context.Context, r *sessions.AuthenticateRequest) (*sessions.AuthenticateResponse, error) {
	id, err := u.userService.Authenticate(r.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &sessions.AuthenticateResponse{Id: (int32)(id)}, nil
}
