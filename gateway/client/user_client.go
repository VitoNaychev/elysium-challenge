package client

import (
	"sync"

	"github.com/VitoNaychev/elysium-challenge/rpc/sessions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	once sync.Once

	conn *grpc.ClientConn
	err  error
)

func NewUserClient(serverURL string) (sessions.UserClient, error) {
	once.Do(func() {
		conn, err = grpc.NewClient(serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	if err != nil {
		return nil, err
	}

	return sessions.NewUserClient(conn), nil
}
