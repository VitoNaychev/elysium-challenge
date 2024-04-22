package main

import (
	"log"
	"net/http"

	"github.com/VitoNaychev/elysium-challenge/gateway/handler"
)

func main() {
	config, err := handler.LoadGatewayConfig("../config/default.yml")
	if err != nil {
		log.Fatal("LoadGatewayConfig error: ", err)
	}

	gateway := handler.NewGateway(config)

	log.Printf("Listening on %v...", config.ListenPort)
	log.Fatal(http.ListenAndServe(config.ListenPort, gateway))
}
