package http_infra

import (
	"encoding/json"
	"log"
	"net/http"

	sys_errors "github.com/tinello/golang-openapi/core/system/errors"
)

func WriteJsonOkResponse(response http.ResponseWriter, body []byte) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(body)
}

func WriteJsonDomainErrorResponse(response http.ResponseWriter, err error) {
	message := "Internal Server Error"
	switch err.(type) {
	case sys_errors.DomainError:
		message = err.Error()
	default:
		log.Printf("[ERROR]: %s", err.Error())
	}
	WriteJsonErrorResponse(response, http.StatusInternalServerError, message)
}

func WriteJsonErrorResponse(response http.ResponseWriter, statusCode int, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(encodeJsonErrorResponse(message))
}

func encodeJsonErrorResponse(message string) []byte {
	jsonError, _ := json.Marshal(map[string]string{"message": message})
	return jsonError
}

func EncodeJsonSuccessResponse() []byte {
	response, _ := json.Marshal(map[string]interface{}{
		"success": true,
	})
	return response
}
