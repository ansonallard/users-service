package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/controller"
	"github.com/ansonallard/users-service/internal/env"
	"github.com/ansonallard/users-service/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	OPENAPI_SPEC_FILE_PATH = "public/openapi.yaml"
)

const (
	OAUTH_TOKEN_ROUTE     = "/oauth/token"
	OAUTH_AUTHORIZE_ROUTE = "/oauth/authorize"
	USERS_LOGIN           = "/users/login"
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

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get the current file path")
	}

	htmlFilePath := filepath.Join(filepath.Dir(filename), filepath.Join("html", "login.html"))

	if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
		log.Fatalf("HTML file does not exist at path: %s", htmlFilePath)
	}

	r.LoadHTMLFiles(htmlFilePath)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, myStruct{})
	})

	r.GET("/login", func(ctx *gin.Context) {

		redirectURI := ctx.Query("redirect_uri")
		parsedURI, err := url.Parse(redirectURI)
		if err != nil || (parsedURI.Scheme == "" && redirectURI != "") {
			log.Printf("Invalid redirect_uri: %s", redirectURI)
			redirectURI = "" // Reset to empty if invalid
		}
		fmt.Println("redirectURI: " + redirectURI)
		ctx.HTML(http.StatusOK, filepath.Base(htmlFilePath), gin.H{
			"RedirectURI": redirectURI,
		})
		// ctx.JSON(http.StatusOK, "hello")
	})

	origins := []string{
		"http://localhost:3000", // React development server
		"http://localhost:8080", // Your current server
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
		"https://oauth.pstmn.io/v1/callback",
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Location", // Important for redirect headers
		},
		AllowCredentials: true,
	}))

	r.POST(USERS_LOGIN, func(ctx *gin.Context) {

		r := ctx.Request
		defer r.Body.Close()

		var body api.LoginDto
		if err := ctx.BindJSON(&body); err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		fmt.Println(body)
		url, err := url.Parse(*body.RedirectUri)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		q := url.Query()
		q.Add("code", "abc123")
		url.RawQuery = q.Encode()
		newUrl := url.String()
		// ctx.JSON(http.StatusOK, myStruct{})

		// Add CORS headers explicitly
		ctx.Header("Access-Control-Allow-Origin", "http://localhost:5000")
		// ctx.Header("Access-Control-Allow-Origin", "https://oauth.pstmn.io/v1/callback")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Expose-Headers", "Location") // Expose Location header for the redirect

		fmt.Println("newUrl: " + newUrl)

		ctx.Redirect(http.StatusFound, newUrl)
	})

	port := env.GetPort()
	log.Printf("Server starting on :%s", port)
	r.Run(fmt.Sprintf(":%s", port))
}

type myStruct struct{}
