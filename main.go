package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ansonallard/users-service/src/controller"
	"github.com/ansonallard/users-service/src/env"
	"github.com/ansonallard/users-service/src/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

const (
	OPENAPI_SPEC_FILE_PATH = "src/public/openapi.yaml"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config=types.cfg.yaml src/public/openapi.yaml

func ValidationMiddleware(router routers.Router) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Find route
			route, pathParams, err := router.FindRoute(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error finding route: %v", err), http.StatusBadRequest)
				return
			}

			// Create validation context
			ctx := context.Background()
			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    r,
				PathParams: pathParams,
				Route:      route,
				Options: &openapi3filter.Options{
					MultiError: true,
				},
			}

			if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
				http.Error(w, fmt.Sprintf("Error validating request: %v", err), http.StatusBadRequest)
				return
			}

			// Create response recorder
			rr := httptest.NewRecorder()
			next.ServeHTTP(rr, r)

			// Copy response
			response := rr.Result()

			// Validate response
			responseValidationInput := &openapi3filter.ResponseValidationInput{
				RequestValidationInput: requestValidationInput,
				Status:                 response.StatusCode,
				Header:                 response.Header,
				Body:                   response.Body,
				Options: &openapi3filter.Options{
					MultiError: true,
				},
			}

			if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
				http.Error(w, fmt.Sprintf("Error validating response: %v", err), http.StatusInternalServerError)
				return
			}

			// Copy the validated response to the original writer
			for k, v := range rr.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(rr.Code)
			w.Write(rr.Body.Bytes())
		})
	}
}

func main() {
	ctx := context.Background()

	// Load and parse OpenAPI spec
	loader := openapi3.NewLoader()
	openAPISpec, err := loader.LoadFromFile(OPENAPI_SPEC_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}

	// Validate the OpenAPI spec itself
	err = openAPISpec.Validate(ctx)
	if err != nil {
		log.Fatalf("Error validating swagger spec: %v", err)
	}

	// Create router from OpenAPI spec
	router, err := gorillamux.NewRouter(openAPISpec)
	if err != nil {
		log.Fatalf("Error creating router: %v", err)
	}

	httpMultiplexer := http.NewServeMux()

	cont := controller.NewBaseController(controller.BaseControllerOpts{
		OidcController: controller.NewOidcController(service.NewOidcService()),
	})
	validationMiddleware := ValidationMiddleware(router)
	httpMultiplexer.Handle("/", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, _, err := router.FindRoute(r)
		// Shouldn't get here as the validation will return
		// route not found error
		if err != nil {
			fmt.Printf("err %+v", err)
			return
		}

		cont.ProcessRequest(controller.ProcessRequestOpts{
			OpenAPIOperationId: route.Operation.OperationID,
			Request:            r,
			Writer:             w})

	})))

	port := env.GetPort()
	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), httpMultiplexer))

}
