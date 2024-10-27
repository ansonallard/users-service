package main

import (
	"ansonallard/users-service/api"
	"ansonallard/users-service/utils"
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

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config=types.cfg.yaml src/public/openapi.yaml

// ServerInterface defines the interface that the generated code expects
// type ServerInterface interface {
// 	GetUser(w http.ResponseWriter, r *http.Request, userID string)
// }

// Server implements ServerInterface
type Server struct{}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	requestData := api.CreateUserRequestDto{}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Fatal(err)
	}
	response := api.CreateUserResponseDto{
		Id: "123",
	}

	message := "error"
	switch {
	case requestData.Username == "bad":
		WriteResponse(w, nil, http.StatusBadRequest, api.BadRequest{Message: &message})
	case requestData.Username == "unauthorized":
		WriteResponse(w, nil, http.StatusUnauthorized, api.UnAuthorized{Message: &message})
	case requestData.Username == "forbidden":
		WriteResponse(w, nil, http.StatusForbidden, api.Forbidden{Message: &message})
	case requestData.Username == "conflict":
		WriteResponse(w, nil, http.StatusConflict, api.Conflict{Message: &message})
	case requestData.Username == "error":
		WriteResponse(w, nil, http.StatusInternalServerError, api.InternalServerError{Message: &message})
	default:
		WriteResponse(w, nil, http.StatusCreated, response)
	}
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
	swagger, err := loader.LoadFromFile("src/public/openapi.yaml")
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

	notImplemented := api.InternalServerError{Message: utils.StrPtr("message")}

	// Register routes with validation
	mux.Handle("/v1/users", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// userID := r.URL.Path[len("/users/"):]
		// server.GetUser(w, r, userID)
		server.CreateUser(w, r)
	})))
	mux.Handle("/v1/users/login", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, nil, http.StatusNotImplemented, notImplemented)
	})))
	mux.Handle("/v1/users/resetPassword", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, nil, http.StatusNotImplemented, notImplemented)
	})))
	mux.Handle("/v1/oauth/token", validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, nil, http.StatusNotImplemented, notImplemented)
	})))

	log.Println("Server starting on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))
}
