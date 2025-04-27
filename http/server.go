package http

import (
	"log"
	"net/http"
	"os"

	"github.com/tinello/golang-openapi/core"
	"github.com/tinello/golang-openapi/http/operations"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

//go:generate mule -e -p http contract.yaml

// func NewServer(provider core.Provider, generateId func() string) http.Handler {
func NewServer(provider *core.Provider, generateId func() string) http.Handler {
	contract, err := ContractResource()
	if err != nil {
		log.Fatalln("Failed to load Swagger contract:", err)
	}

	applicationVersion := GetApplicationVersion()
	operationHandlers := map[string]http.Handler{}
	for _, op := range []operation{
		operations.NewServiceInfo("0.0.1"),
	} {
		operationHandlers[op.GetId()] = NewOperationHandler(op, provider, generateId, applicationVersion)
	}

	return h2cHandler(Logger{NewOpenApiRouter(contract, operationHandlers)})
}

func GetApplicationVersion() string {
	if value, ok := os.LookupEnv("VERSION"); ok {
		return value
	}
	return "local-version"
}

func h2cHandler(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}
