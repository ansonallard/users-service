// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Defines values for OAuth2ErrorSchemaError.
const (
	InvalidClient        OAuth2ErrorSchemaError = "invalid_client"
	InvalidGrant         OAuth2ErrorSchemaError = "invalid_grant"
	InvalidRequest       OAuth2ErrorSchemaError = "invalid_request"
	InvalidScope         OAuth2ErrorSchemaError = "invalid_scope"
	UnauthorizedClient   OAuth2ErrorSchemaError = "unauthorized_client"
	UnsupportedGrantType OAuth2ErrorSchemaError = "unsupported_grant_type"
)

// Defines values for OAuth2TokenRequestGrantType.
const (
	AuthorizationCode OAuth2TokenRequestGrantType = "authorization_code"
	ClientCredentials OAuth2TokenRequestGrantType = "client_credentials"
	Password          OAuth2TokenRequestGrantType = "password"
	RefreshToken      OAuth2TokenRequestGrantType = "refresh_token"
)

// Defines values for OAuth2TokenResponseTokenType.
const (
	Bearer OAuth2TokenResponseTokenType = "Bearer"
)

// CreateUserRequestDto defines model for CreateUserRequestDto.
type CreateUserRequestDto struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

// OAuth2ErrorSchema defines model for OAuth2ErrorSchema.
type OAuth2ErrorSchema struct {
	Error            OAuth2ErrorSchemaError `json:"error"`
	ErrorDescription *string                `json:"error_description,omitempty"`
	ErrorUri         *string                `json:"error_uri,omitempty"`
}

// OAuth2ErrorSchemaError defines model for OAuth2ErrorSchema.Error.
type OAuth2ErrorSchemaError string

// OAuth2TokenRequest defines model for OAuth2TokenRequest.
type OAuth2TokenRequest struct {
	GrantType    OAuth2TokenRequestGrantType `json:"grant_type"`
	RefreshToken *string                     `json:"refresh_token"`
	Scope        *string                     `json:"scope"`
}

// OAuth2TokenRequestGrantType defines model for OAuth2TokenRequest.GrantType.
type OAuth2TokenRequestGrantType string

// OAuth2TokenResponse defines model for OAuth2TokenResponse.
type OAuth2TokenResponse struct {
	AccessToken  string                       `json:"access_token"`
	ExpiresIn    float32                      `json:"expires_in"`
	RefreshToken *string                      `json:"refresh_token,omitempty"`
	Scope        *string                      `json:"scope,omitempty"`
	TokenType    OAuth2TokenResponseTokenType `json:"token_type"`
}

// OAuth2TokenResponseTokenType defines model for OAuth2TokenResponse.TokenType.
type OAuth2TokenResponseTokenType string

// ResetUserPasswordDto defines model for ResetUserPasswordDto.
type ResetUserPasswordDto struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
	Username    string `json:"username"`
}

// BadRequest defines model for BadRequest.
type BadRequest struct {
	Message *string `json:"message,omitempty"`
}

// Conflict defines model for Conflict.
type Conflict struct {
	Message *string `json:"message,omitempty"`
}

// CreateUserResponseDto User created
type CreateUserResponseDto struct {
	Id string `json:"id"`
}

// Forbidden defines model for Forbidden.
type Forbidden struct {
	Message *string `json:"message,omitempty"`
}

// InternalServerError defines model for InternalServerError.
type InternalServerError struct {
	Message *string `json:"message,omitempty"`
}

// OAuthErrorBadRequest defines model for OAuthErrorBadRequest.
type OAuthErrorBadRequest = OAuth2ErrorSchema

// UnAuthorized defines model for UnAuthorized.
type UnAuthorized struct {
	Message *string `json:"message,omitempty"`
}

// OAuth2TokenFormdataRequestBody defines body for OAuth2Token for application/x-www-form-urlencoded ContentType.
type OAuth2TokenFormdataRequestBody = OAuth2TokenRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserRequestDto

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = CreateUserRequestDto

// ResetPasswordJSONRequestBody defines body for ResetPassword for application/json ContentType.
type ResetPasswordJSONRequestBody = ResetUserPasswordDto

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYS3PiOBD+Ky7t3NaAB5LUhlsek5TJkAeQSXamspSw2iCwJa8kB0iK/74liYcNTiCp",
	"mc0e9oal7tbX3V+3WjyjgMcJZ8CURPVnJEAmnEkwH8eYtODvFKTSXwFnCpj5iZMkogFWlLPKUHKm12Qw",
	"gBjrX4ngCQhFrZEYpMR90D/VNAFUR1IJyvpoNnMXK7w3hEChmV4iIANBE20b1TUER8wxzFx0wlkY0eDj",
	"8LRA8lQE4AQLJBqVAKzgVoJozcN3qvibIGJCqN7C0XUGbIgjCesI9DFOYE4kyF1zjpJiv3QIqQCC6j+0",
	"zMMOnmrHzrjoUUKAfVi8VwhmLvKZAsFw1AbxCOKLEFx8GLAFFseCcSyamYuujlI1MF/vLJ9PAkJUR79V",
	"VpVZsbuyYoxXjfW2VdheM7dMK3FBn4B8WLguuXIyMLTA3CdtMFtABve8fl4tizy8BEs55sJ4GFP2FVhf",
	"DVD9D3cdrYtSqTMXg1VTOo+ojv76gUtPR6XvXunw4fdPyN1SRUsj7urozbJy0Vfep2zRFt7ok+IjW3qF",
	"Lvg71Ppczp2bKgK4SamN1MOi0oClsWkh7BFHlHQXLHOXK0FENbdWC32BzXfK8DL9K6mUyTRJuFAwl+wa",
	"fCt1GfAEMrhXITCgujmWPb8klQq6PVbWy5dD1NEhzBR0PkYZ9JlALXw2JdYNOMnxRZ8fCpCDrk2Pi2xg",
	"uoEAAkxRHMlC3/Nq9WfE0ijCvQhQXYkUCjRsHLdLrsUk49XWwLyL4jgIQMqVIzDBcaLRIZg2Br3zgF7R",
	"hn/75H++pL70WWs/OPEP/FFy/+2kcViGaeOJ3Pn0ivqT5rDpXXb+rF2djsY+HdNefKa+t43wIz7f67fO",
	"DyO9ju/OPH/IJ5edL9XmsLnfPPWn4U25HUYXk3Gr0W7CxcVZ9aazF46TJjTC2sH11ehg2vjWxeRGyvF+",
	"gIroOEmoANmleTdqnrcUZmncA1GYv5XbWJI+o5hQVq0NQ8wlCZnkOGQ4kiPCQl5lNRzSkI2qtSHDY0Cv",
	"ZXtjx5y4QdNjwAJEAdXW6JDLVs5YLgBFVGmBBKU7/PWc/W9v8QzG17t3eR6R64+7E7KnuznkBbOXiyQE",
	"qaBqajrwfPg2OeksGNIzn2dcxFihOmrcddD8CtWm7O4K4ECpxF7DlIUmzooqwy+NUZYkiEcaaKCPIKS9",
	"or3y57JnIpcAwwlFdVQre2XPdCw1MKAqXHe0ypK3Cbe9MH/Z27bgGOyOr4kdA1OmByJjXpjf+gLLthBk",
	"owlSHXMyfWVUmZTG43Ep5CIupSICpvsqeesQlWvns5nNZebxU/W8nzy25ftkwZDUTk196Rzs2dOLjC5R",
	"VgpHTaP8+b3KukXgvtRsPjJgnBPOlOARetB7FUOflzNvBzkHO1rOmRN8PeWraW/njL8t3IXj5Cxfrfr2",
	"20z5DnErfuztmrJ3JCo3vxul2nal3Ltpzzvcwa/Mg3b//n67QtFrbDf+VCI9FL/MIjMzO5Q5egZSU6fk",
	"CFCpYNJZ3Dx5Qhn5/xyXfl77yL8hChqHDZhYSvx7ZPzVTBF6bMje48WMMdOFsxisnZAL24G2MaeVM/9r",
	"GFQ4+fzPoOX/AC+zYDb7JwAA//9QMTzKnxQAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
