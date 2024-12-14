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

// Defines values for OAuth2AuthorizationRequestResponseType.
const (
	Code OAuth2AuthorizationRequestResponseType = "code"
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

// LoginDto defines model for LoginDto.
type LoginDto struct {
	Password    string  `json:"password"`
	RedirectUri *string `json:"redirect_uri,omitempty"`
	Username    string  `json:"username"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

// OAuth2AuthorizationRequest defines model for OAuth2AuthorizationRequest.
type OAuth2AuthorizationRequest struct {
	ClientId     string                                 `json:"client_id"`
	RedirectUri  *string                                `json:"redirect_uri,omitempty"`
	ResponseType OAuth2AuthorizationRequestResponseType `json:"response_type"`
	Scope        *string                                `json:"scope,omitempty"`
}

// OAuth2AuthorizationRequestResponseType defines model for OAuth2AuthorizationRequest.ResponseType.
type OAuth2AuthorizationRequestResponseType string

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
	Code         *string                     `json:"code,omitempty"`
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

// OAuth2AuthorizationFormdataRequestBody defines body for OAuth2Authorization for application/x-www-form-urlencoded ContentType.
type OAuth2AuthorizationFormdataRequestBody = OAuth2AuthorizationRequest

// OAuth2TokenFormdataRequestBody defines body for OAuth2Token for application/x-www-form-urlencoded ContentType.
type OAuth2TokenFormdataRequestBody = OAuth2TokenRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = LoginDto

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = CreateUserRequestDto

// ResetPasswordJSONRequestBody defines body for ResetPassword for application/json ContentType.
type ResetPasswordJSONRequestBody = ResetUserPasswordDto

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xY33PaOBD+Vzy6vp0BF5LOhbc2bTqmJT+ANLl2coyw1iAwkk+SAyTj//1Gsg02Mb8y",
	"TdOHe7Pllfzt7refVnpEHp+GnAFTEjUfkQAZcibBvHzApAP/RiCVfvM4U8DMIw7DgHpYUc5qY8mZHpPe",
	"CKZYP4WChyAUTRaZgpR4CPpRLUJATSSVoGyI4tjORvhgDJ5CsR4iID1BQ702amoIlkgxxDY65cwPqPd6",
	"eDogeSQ8sLwMiUYlACu4liA6afg+Kn4QREwI1Z9wcJkD6+NAwjoC/RvLM38kyF5zjpJyv3QIqQCCmj+0",
	"zd0enmrHzrgYUEKAvVq8VwhiG7lMgWA46IK4B/FJCC5eDViGxUrAWAma2EYX7yM1Mm/PLJ83AnzURH/U",
	"VpVZS77Kmlm8blbvJhN218w105O4oA9AXi1c51xZORjaIPVJL5gvIIM7rZ+tZVGEF2IpZ1wYD6eUfQU2",
	"VCPU/MteR2ujSOrMTSGZpnQeURP98wNXHt5XvjuVk7s/3yB7RxUtF7FXv35aVjb6yoeUvag7AggV4Kl+",
	"JKg297mYYoWaSL//Fu5nqnhgDBSfJMpT6oK7h9Sldna6VBnApKIyappSyBXtAWi9gAJT/VIFPjxJ2T7c",
	"T748ImDRVLvkcQI5R1YzpMdD2B2S4sJ2Dvbm6OT15okuQCbDGULK7nFAST+TIHs5kvwrNzAU2LxHDC+1",
	"YWUVMRmFIRcKUssMczY9cbksGAZUvyBBj5us0oxsD1vi5eYQ9TTBcsRZ44ZOWhmCnFu5COI8G/tmcq7M",
	"NDBfgBz1E1Yvc+gJIMAUxYEsDUpxWvMRsSgI8CAA1FQigm2c2mG5FqycVzsj9ixlwJ4HUq4cgTmehhod",
	"gkVrNPjs0Qvacq8f3Lfn1JUu6xx7p+47dxLefjttnVRh0XogNy69oO68PW47572/GxcfJzOXzuhgeqa+",
	"d43xPf58NOx8Pgn0OL45c9wxn5/3PtXb4/Zx+6O78K+qXT/4Mp91Wt02fPlyVr/qHfmzsA0tv/Hu8mLy",
	"btH61sfkSsrZsVdW5jAPqQDZp0U3Go6zNGbRdACiNH8rt7EkQ0YxoazeGPuYS+IzybHPcCAnhPm8zhrY",
	"pz6b1BtjhmeADlCQVD6f0PQDYAGihGprdChkq7BYIQBlVOmABKX7gsuU/YfvpAxml/tvpjwgl6/XSeT/",
	"bheQl3TsNpLgRYKqhZHm9MhmctLLGDIwr2fZZtO66aG08dJLJV9XAEdKhUnzRplv4qyoMvzSGGVFgrin",
	"ngZ6D0ImjZ1TfVt1TORCYDikqIkaVafqGMVSIwOqxrWi1QqypseHYLSy2Ckm6mAVdmTL1TyfAlPJXPM3",
	"YZ51G1C2iaMkyCDVB04WW/reeWU2m1X0flyJRABMyy05tCMv7R/iOM7t5CYSDaf+1ONuZAoE2WgEmIAw",
	"ll+5t4xT0fzUqH3W4VM2tAqBtdL9YgV/e69hQB45ziZHl/hrpScbM/ntcydrbcFDqcvgvQmCdcqZEjxA",
	"d/pbSpyl4IVcbmaMIf1+TOmlUvQLGVJoEEqYUU8y8BNPicUNtuRMlvHuN86/0Z3NmU/OjRa2tJ2VKuN6",
	"yleHy70zfli4l0e8uCjtulV6muY9YlV+n7Rvmp6RnMIVgZnU2D2pcDVz5Jzs4Vfuzuz49nb3hLILn/04",
	"Uwt0UjYzx+TMoszSDbNaWBVLgIoEk1bWphRJZOxfiD+ltx97ccn5uRzeJhZJwMTS4teR8aWZInSPmW/6",
	"yhljWlErO4VZPheJ6uxiTqew/MswqLRN/p9By6vGzSyI4/8CAAD//6VFvjsCGQAA",
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
