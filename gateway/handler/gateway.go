package handler

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type Route struct {
	Name         string `mapstructure:"name"`
	Context      string `mapstructure:"context"`
	Target       string `mapstructure:"target"`
	Authenticate bool   `mapstructure:"authenticate,omitempty"`
}

type GatewayConfig struct {
	ListenPort       string  `mapstructure:"listenPort"`
	AuthenticateAddr string  `mapstructure:"authenticateAddr"`
	Routes           []Route `mapstructure:"routes"`
}

func LoadGatewayConfig(path string) (GatewayConfig, error) {
	viper.SetConfigType("yaml")

	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		return GatewayConfig{}, nil
	}

	gatewayConfig := GatewayConfig{}
	viper.UnmarshalKey("gateway", &gatewayConfig)

	return gatewayConfig, nil
}

type Gateway struct {
	http.Handler
}

func NewGateway(config GatewayConfig) *Gateway {
	var (
		proxy http.Handler
		err   error
	)

	mux := http.NewServeMux()

	for _, route := range config.Routes {
		if route.Authenticate {
			proxy, err = NewAuthProxy(route.Target, config.AuthenticateAddr)
		} else {
			proxy, err = NewProxy(route.Target)
		}
		if err != nil {
			log.Println("NewProxy error: ", err)
		}

		log.Printf("Mapping '%v' | %v ---> %v", route.Name, route.Context, route.Target)
		mux.HandleFunc(route.Context, proxy.ServeHTTP)
	}

	return &Gateway{
		Handler: mux,
	}
}
