package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

// User represents our example data structure
type User struct {
	ID string `json:"id"`
}

func main() {
	// Load and parse OpenAPI spec
	openapiSpec, err := openapi3.NewLoader().LoadFromFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// Validate the OpenAPI spec itself
	err = openapiSpec.Validate(context.Background())
	if err != nil {
		log.Fatalf("Error validating OpenAPI spec: %v", err)
	}

	// Create a router from the OpenAPI spec
	router, err := gorillamux.NewRouter(openapiSpec)
	if err != nil {
		log.Fatalf("Error creating router: %v", err)
	}

	// Example handler
	getUserHandler := func(w http.ResponseWriter, r *http.Request) {
		user := User{
			ID: "1",
		}
		json.NewEncoder(w).Encode(user)
	}

	middleware := middleware{Router: router}
	// Register routes with validation middleware
	http.HandleFunc("/v1/users", middleware.validationMiddleware(getUserHandler))

	port := 5000
	log.Printf("Server starting on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type middleware struct {
	Router routers.Router
}

// Create a validation middleware
func (m *middleware) validationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find route
		route, pathParams, err := m.Router.FindRoute(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error finding route: %v", err), http.StatusBadRequest)
			return
		}

		// Validate request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
		}

		if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
			http.Error(w, fmt.Sprintf("Error validating request: %v", err), http.StatusBadRequest)
			return
		}

		// Create a response recorder to capture the response
		rr := httptest.NewRecorder()
		next.ServeHTTP(rr, r)

		// Validate response
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 rr.Code,
			Header:                 rr.Header(),
		}

		if err := openapi3filter.ValidateResponse(r.Context(), responseValidationInput); err != nil {
			http.Error(w, fmt.Sprintf("Error validating response: %v", err), http.StatusInternalServerError)
			return
		}

		// Copy the validated response to the original response writer
		for k, v := range rr.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(rr.Code)
		w.Write(rr.Body.Bytes())
	}
}
