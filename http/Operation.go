package http

import (
	"net"
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
		cachedServerIp:     "",
	}
}

type operationHandler struct {
	op                 operation
	provider           *core.Provider
	applicationVersion string
	cachedServerIp     string
	generateId         func() string
}

func (h *operationHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	h.op.Execute(response, request, *h.provider)
}

func (h *operationHandler) serverIpFor(request *http.Request) string {
	if h.cachedServerIp == "" {
		hostIps, err := net.LookupHost(removePort(request.Host))
		if err == nil && len(hostIps) > 0 {
			h.cachedServerIp = hostIps[0]
		} else {
			h.cachedServerIp = getLocalServerIp()
		}
	}
	return h.cachedServerIp
}

func removePort(host string) string {
	return strings.Split(host, ":")[0]
}

func getLocalServerIp() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "0.0.0.0"
}

type operation interface {
	Execute(response http.ResponseWriter, request *http.Request, provider core.Provider)
	GetId() string
}
