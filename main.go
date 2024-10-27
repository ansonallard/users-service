package main

import (
	"ansonallard/users-service/api"
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

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config=types.cfg.yaml openapi.yaml

// Create types.cfg.yaml with:
/*
generate:
  models: true
  embedded-spec: true
output: types.gen.go
package: api
*/

// ServerInterface defines the interface that the generated code expects
type ServerInterface interface {
	GetUser(w http.ResponseWriter, r *http.Request, userID string)
}

// Server implements ServerInterface
type Server struct{}

// GetUser implements ServerInterface
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request, userID string) {
	user := api.User{
		Id:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	WriteResponse(w, map[string]string{
		"Content-Type": "application/json",
	}, http.StatusCreated, user)
}

func WriteResponse(w http.ResponseWriter, headers map[string]string, statusCode int, body interface{}) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

// ValidationMiddleware creates middleware that validates requests and responses
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
	swagger, err := loader.LoadFromFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}

	// Validate the OpenAPI spec itself
	err = swagger.Validate(ctx)
	if err != nil {
		log.Fatalf("Error validating swagger spec: %v", err)
	}

	// Create router from OpenAPI spec
	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		log.Fatalf("Error creating router: %v", err)
	}

	// Create the server implementation
	server := &Server{}

	// Create validation middleware
	validationMiddleware := ValidationMiddleware(router)

	// Create handler with routes
	mux := http.NewServeMux()

	// Register routes with validation
	mux.Handle("/users/", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Path[len("/users/"):]
		server.GetUser(w, r, userID)
	})))

	log.Println("Server starting on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))
}
