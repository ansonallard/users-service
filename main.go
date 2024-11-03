package main

import (
	"ansonallard/users-service/api"
	"ansonallard/users-service/src/operations"
	"ansonallard/users-service/utils"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

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

	// if err = os.Mkdir("./users"); err != nil {

	// }

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
	openAPISpec, err := loader.LoadFromFile("src/public/openapi.yaml")
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

	// Create the server implementation
	server := &Server{}

	// Create validation middleware
	validationMiddleware := ValidationMiddleware(router)

	// Create handler with routes
	mux := http.NewServeMux()

	notImplemented := api.InternalServerError{Message: utils.StrPtr("message")}

	// Register routes with validation

	// TODO: add error handling
	url := openAPISpec.Servers[0].URL

	mux.Handle(utils.ReturnPathWithTrailingSlash(url), validationMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, _, err := router.FindRoute(r)
		if err != nil {
			fmt.Printf("err %+v", err)
			return
		}

		switch route.Operation.OperationID {
		case operations.CREATE_USER:
			server.CreateUser(w, r)
		case operations.LOGIN:
			fallthrough
		case operations.RESET_PASSWORD:
			fallthrough
		case operations.OAUTH_CLIENT_CREDS:
			fallthrough
		default:
			WriteResponse(w, nil, http.StatusNotImplemented, notImplemented)
		}
	})))

	privateKeyFileBytes, err := os.ReadFile("private.pem")
	var privateKey *rsa.PrivateKey
	if err != nil {
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			fmt.Printf("Cannot generate RSA Key\n")
			os.Exit(1)
		}
		if err = encodePrivateKey(privateKey, "private.pem"); err != nil {
			os.Exit(1)
		}
	} else {
		privateKeyPem, _ := pem.Decode(privateKeyFileBytes)
		privateKey, _ = x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{Id: "1234",
	// 	Issuer: "authorization.ansonallard.com", IssuedAt: time.Now().Unix(),
	// 	ExpiresAt: time.Now().Add(time.Minute * 5).Unix()})
	// signedToken, err := token.SignedString(privateKey)

	// if err != nil {
	// 	fmt.Printf("%+v", err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("jwt: %s\n", signedToken)

	publicKeyFileBytes, err := os.ReadFile("public.pem")
	// var publicKey *rsa.PublicKey
	if err != nil {
		if publicKeyFileBytes, err = encodePublicKey(&privateKey.PublicKey, "public.pem"); err != nil {
			os.Exit(1)
		}
	} else {
		publicKeyPem, _ := pem.Decode(publicKeyFileBytes)
		_, _ = x509.ParsePKCS1PublicKey(publicKeyPem.Bytes)
	}

	mux.HandleFunc("/well-known/jwks.json", func(w http.ResponseWriter, r *http.Request) {
		jwk := JWKJson{
			Keys: &[]JWKKey{{Kty: "RSA", Use: "sig", Alg: "RSA256", X5c: []string{string(publicKeyFileBytes)}}},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(jwk)
	})

	log.Println("Server starting on :5000")
	log.Fatal(http.ListenAndServe(":5000", mux))

	// myToken, err := jwt.ParseWithClaims(signedToken, jwt.Claims{}, func(token *jwt.Token) (interface{}, error) {
	// 	// Validate signing method
	// 	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return publicKey, nil
	// })

}

type JWKJson struct {
	Keys *[]JWKKey `json:"keys"`
}

type JWKKey struct {
	X5c []string `json:"x5c"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
}

func encodePublicKey(key *rsa.PublicKey, filename string) ([]byte, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(key)
	publicKeyBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error when creating %s %s\n", filename, err)
		return nil, err
	}
	err = pem.Encode(publicPem, &publicKeyBlock)
	if err != nil {
		fmt.Printf("error when encoding %s %s\n", filename, err)
		return nil, err
	}
	result, _ := os.ReadFile(filename)
	return result, nil
}

func encodePrivateKey(key *rsa.PrivateKey, filename string) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error when creating %s %s\n", filename, err)
		return err
	}
	err = pem.Encode(privatePem, &privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encoding %s %s\n", filename, err)
		return err
	}
	return nil
}
