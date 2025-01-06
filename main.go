package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	OPENAPI_SPEC_FILE_PATH = "public/tenants.openapi.yaml"
)

const (
	OAUTH_TOKEN_ROUTE     = "/oauth/token"
	OAUTH_AUTHORIZE_ROUTE = "/oauth/authorize"
	USERS_LOGIN           = "/users/login"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config=types.cfg.yaml public/tenants.openapi.yaml

type responseWriter struct {
	gin.ResponseWriter
	body []byte
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func ValidationMiddleware(router routers.Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		route, pathParams, err := router.FindRoute(c.Request)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Error finding route: %v", err)})
			c.Abort()
			return
		}

		ctx := context.Background()

		// Validate Request
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    c.Request,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				MultiError: true,
				AuthenticationFunc: func(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		}

		if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error validating request: %v", err)})
			c.Abort()
			return
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           []byte{},
		}
		c.Writer = writer

		// Process the request
		c.Next()

		// Validate Response
		responseValidationInput := &openapi3filter.ResponseValidationInput{
			RequestValidationInput: requestValidationInput,
			Status:                 writer.Status(),
			Header:                 writer.Header(),
			Body:                   io.NopCloser(bytes.NewReader(writer.body)),
			Options: &openapi3filter.Options{
				MultiError: true,
			},
		}

		// Third party package does not validate non application/json response types
		// Do not validate responses for text/html respones types
		if route.Operation.OperationID == "LoginPage" {
			responseValidationInput.Options.ExcludeResponseBody = true
		}

		if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
			log.Printf("Error validating response: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error validating response: %v", err)})
			return
		}
	}
}

func configureDb(ctx context.Context, mongoClient *mongo.Client) {
	usersServiceDb := mongoClient.Database("users-service")
	tenantsCollection := usersServiceDb.Collection("tenants")
	if tenantsCollection == nil {
		usersServiceDb.CreateCollection(ctx, "tenants")
		tenantsCollection = usersServiceDb.Collection("tenants")
	}
	// tenantsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{{}})
}

func main() {
	ctx := context.Background()

	mongoClient, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

	configureDb(ctx, mongoClient)

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

	// Create Gin router
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.New()
	ginRouter.Use(gin.Recovery())

	ginRouter.Use(ValidationMiddleware(router))
	cont := controller.NewOidcController(service.NewOidcService())
	tenantsService := service.NewTenantService(mongoClient)
	tenantsControllers := controller.NewTenantsContorller(tenantsService)
	usersService := service.NewUsersService(&tenantsService, mongoClient)
	usersController := controller.NewUsersController(&usersService)

	// Get current working directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Unable to get the current file path")
	}

	htmlFilePath := filepath.Join(filepath.Dir(filename), filepath.Join("html", "login.html"))

	if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
		log.Fatalf("HTML file does not exist at path: %s", htmlFilePath)
	}

	ginRouter.LoadHTMLFiles(htmlFilePath)

	ginRouter.Any("/*path", func(c *gin.Context) {
		route, pathParams, err := router.FindRoute(c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error finding route: %v", err)})
			return
		}

		var operationErr error

		switch route.Operation.OperationID {
		case "OAuth2Authorize":
			err := cont.OAuth2Authorize(c)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		case "OAuth2TokenOptions":
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			c.JSON(http.StatusOK, myStruct{})
		case "OAuth2Token":
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			operationErr = cont.OAuth2Token(c)
		case "LoginPage":
			redirectURI := c.Query("redirect_uri")
			parsedURI, err := url.Parse(redirectURI)
			if err != nil || (parsedURI.Scheme == "" && redirectURI != "") {
				log.Printf("Invalid redirect_uri: %s", redirectURI)
				redirectURI = "" // Reset to empty if invalid
			}
			fmt.Println("redirectURI: " + redirectURI)
			hostName := parsedURI.Scheme + "://" + parsedURI.Host
			c.Header("Access-Control-Allow-Origin", hostName)
			c.HTML(http.StatusOK, filepath.Base(htmlFilePath), gin.H{
				"RedirectURI":                 redirectURI,
				"Access-Control-Allow-Origin": redirectURI,
			})
		case "login":
			defer c.Request.Body.Close()

			var body api.UserLoginRequest
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			fmt.Println(body)
			url, err := url.Parse(*body.RedirectUri)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			q := url.Query()
			q.Add("code", "abc123")
			url.RawQuery = q.Encode()
			newUrl := url.String()

			// Add CORS headers explicitly
			// ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
			// // ctx.Header("Access-Control-Allow-Origin", "https://oauth.pstmn.io/v1/callback")
			// ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			// ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
			// ctx.Header("Access-Control-Expose-Headers", "Location") // Expose Location header for the redirect

			fmt.Println("newUrl: " + newUrl)

			c.JSON(http.StatusOK, api.UserLoginResponse{RedirectUrl: newUrl})
		case "createTenant":
			operationErr = tenantsControllers.CreateTenant(ctx, c)
		case "createUser":
			operationErr = usersController.CreateUser(ctx, c, pathParams)
		default:
			fmt.Println(route.Operation.OperationID)
			c.JSON(http.StatusNoContent, myStruct{})
		}

		if operationErr != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	port := env.GetPort()
	log.Printf("Server starting on :%s", port)
	ginRouter.Run(fmt.Sprintf(":%s", port))
}

type myStruct struct{}
