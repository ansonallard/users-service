package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ansonallard/users-service/internal/controller"
	"github.com/ansonallard/users-service/internal/env"
	"github.com/ansonallard/users-service/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gin-gonic/gin"
)

const (
	OPENAPI_SPEC_FILE_PATH = "public/openapi.yaml"
)

const (
	OAUTH_TOKEN_ROUTE     = "/oauth/token"
	OAUTH_AUTHORIZE_ROUTE = "/oauth/authorize"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config=types.cfg.yaml public/openapi.yaml

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

	r := gin.Default()

	cont := controller.NewOidcController(service.NewOidcService())

	r.POST(OAUTH_TOKEN_ROUTE, func(ctx *gin.Context) {
		err := cont.OAuth2Token(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})

	r.GET(OAUTH_AUTHORIZE_ROUTE, func(ctx *gin.Context) {
		err := cont.OAuth2Authorize(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	port := env.GetPort()
	log.Printf("Server starting on :%s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
