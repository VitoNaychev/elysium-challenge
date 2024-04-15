package handler_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/VitoNaychev/elysium-challenge/assert"
	"github.com/VitoNaychev/elysium-challenge/rpc/sessions"
	"github.com/VitoNaychev/elysium-challenge/sessions/handler"
	"github.com/VitoNaychev/elysium-challenge/sessions/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthenticateRPC(t *testing.T) {
	t.Run("returns user ID on valid JWT", func(t *testing.T) {
		wantUserID := 10
		dummyJWT := "sampleToken"

		request := &sessions.AuthenticateRequest{Token: dummyJWT}

		userService := StubUserService{
			dummyUserID: wantUserID,
			dummyJWT:    dummyJWT,
		}
		userHandler := handler.NewUserRPCHandler(&userService)

		response, err := userHandler.Authenticate(context.Background(), request)
		assert.RequireNoError(t, err)

		assert.Equal(t, response.Id, int32(wantUserID))
	})

	t.Run("returns Unauthenticated on invalid JWT", func(t *testing.T) {
		dummyErr := service.ErrInvalidJWT
		dummyJWT := "sampleToken"

		request := &sessions.AuthenticateRequest{Token: dummyJWT}

		userService := StubUserService{
			dummyErr: dummyErr,
			dummyJWT: dummyJWT,
		}
		userHandler := handler.NewUserRPCHandler(&userService)

		_, err := userHandler.Authenticate(context.Background(), request)

		statusError, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected status.Error, got %v", reflect.TypeOf(err))
		}

		assert.Equal(t, statusError.Code(), codes.Unauthenticated)
	})
}
