package operations

import (
	"encoding/json"
	"net/http"

	"github.com/tinello/golang-openapi/core"
	sys_domain "github.com/tinello/golang-openapi/core/system/domain"
	http_infra "github.com/tinello/golang-openapi/http/infrastructure"
)

func NewServiceInfo(applicationVersion string) *serviceInfo {
	return &serviceInfo{
		applicationVersion: applicationVersion,
	}
}

type serviceInfo struct {
	applicationVersion string
}

func (s *serviceInfo) Execute(response http.ResponseWriter, _ *http.Request, provider core.Provider) {
	result := provider.GetServiceInfo().Execute()
	if !result.Healthy() {
		http_infra.WriteJsonDomainErrorResponse(response, result.Error)
		return
	}
	http_infra.WriteJsonOkResponse(response, encodeJsonServiceInfoResponse(result, s.applicationVersion))
}

func (*serviceInfo) GetId() string {
	return "service_info"
}

func encodeJsonServiceInfoResponse(serviceInfo *sys_domain.ServiceInfo, version string) []byte {
	data := jsonServiceInfoResponse{
		Name:    "golang-openapi",
		Version: version,
		Healthy: serviceInfo.Healthy(),
	}
	response, _ := json.Marshal(data)
	return response
}

type jsonServiceInfoResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Healthy bool   `json:"healthy"`
}

type mapResult map[string]interface{}
