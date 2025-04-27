package http

import (
	"net/http"
	"strings"

	"github.com/tinello/golang-openapi/core"
)

func NewOperationHandler(op operation, provider *core.Provider, generateId func() string, applicationVersion string) *operationHandler {
	return &operationHandler{
		op:                 op,
		provider:           provider,
		generateId:         generateId,
		applicationVersion: applicationVersion,
	}
}

type operationHandler struct {
	op                 operation
	provider           *core.Provider
	applicationVersion string
	generateId         func() string
}

func (h *operationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.op.Execute(response, request, *h.provider)
}

func removePort(host string) string {
	return strings.Split(host, ":")[0]
}

type operation interface {
	Execute(response http.ResponseWriter, request *http.Request, provider core.Provider)
	GetId() string
}
