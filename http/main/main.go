package main

import (
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/tinello/golang-openapi/core"
	http_delivery "github.com/tinello/golang-openapi/http"
)

func main() {
	provider := core.GetProviderInstance()

	server := &http.Server{
		Handler: http_delivery.NewServer(&provider, generateId),
	}

	listenAddress := ":" + core.MustGetEnv("PORT")
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(
		"%s:%s server started on port %s...\n",
		"golang-openapi",
		http_delivery.GetApplicationVersion(),
		listenAddress)
	log.Fatalln(server.Serve(listener))
}

func generateId() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
