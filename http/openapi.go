package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	http_infra "github.com/tinello/golang-openapi/http/infrastructure"
)

var (
	operationNotFound = errors.New("Not found")
)

func NewOpenApiRouter(contract []byte, operations map[string]http.Handler) http.Handler {
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(contract)
	if err != nil {
		log.Fatalln("Failed to load swagger from file:", err)
	}

	errValidate := doc.Validate(loader.Context)
	if errValidate != nil {
		log.Fatalln("Failed to validate swagger from file:", errValidate)
	}

	router, errRouter := gorillamux.NewRouter(doc)
	if errRouter != nil {
		log.Fatalln("Failed to create router:", errRouter)
	}

	return &openApiRouter{
		router:     router,
		operations: operations,
	}
}

type openApiRouter struct {
	router     routers.Router
	operations map[string]http.Handler
}

func (r *openApiRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, err := r.validateRequestAndGetOperationHandler(request)
	switch err {
	case nil:
	case operationNotFound:
		http_infra.WriteJsonErrorResponse(response, http.StatusNotFound, err.Error())
		return
	default:
		http_infra.WriteJsonErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	handler.ServeHTTP(response, request)
}

func (r *openApiRouter) validateRequestAndGetOperationHandler(request *http.Request) (http.Handler, error) {
	route, pathParams, err := r.router.FindRoute(request)
	if err != nil {
		return nil, operationNotFound
	}

	handler, found := r.operations[route.Operation.OperationID]
	if !found {
		return nil, operationNotFound
	}

	err = openapi3filter.ValidateRequest(
		request.Context(),
		&openapi3filter.RequestValidationInput{
			Request:    request,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				AuthenticationFunc: func(_ context.Context, _ *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		},
	)

	if err != nil {
		return nil, errors.New(strings.Replace(err.Error(), "\n", " ", -1))
	}

	return handler, nil
}
