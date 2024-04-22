package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/VitoNaychev/elysium-challenge/gateway/client"
	"github.com/VitoNaychev/elysium-challenge/rpc/sessions"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthProxy struct {
	proxy  *httputil.ReverseProxy
	client sessions.UserClient
}

func NewAuthProxy(targetURL, authURL string) (*AuthProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	client, err := client.NewUserClient(authURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return &AuthProxy{
		proxy:  proxy,
		client: client,
	}, nil
}

func (a *AuthProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := a.client.Authenticate(context.Background(), &sessions.AuthenticateRequest{Token: ""})
	if err != nil {
		statusError, ok := status.FromError(err)
		if ok && statusError.Code() == codes.Unauthenticated {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
		}
		return
	}

	a.proxy.ServeHTTP(w, r)
}
