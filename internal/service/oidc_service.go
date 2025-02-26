package service

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/constants"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/keys"
	"github.com/ansonallard/users-service/internal/utils"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OidcServiceConfig struct {
	MongoClient *mongo.Client
	PrivateKey  *rsa.PrivateKey
}

type OidcService interface {
	OAuth2Authorize(request api.OAuth2AuthorizationRequest) (redirectUrl *url.URL, err error)
	Oauth2Token(request OAuth2TokenInput) (response *api.OAuth2TokenResponse, err error)
}

type oidcService struct {
	db         *mongo.Client
	privateKey *rsa.PrivateKey
}

func NewOidcService(config *OidcServiceConfig) (OidcService, error) {
	if config.MongoClient == nil {
		return nil, fmt.Errorf("db not provided")
	}
	if config.PrivateKey == nil {
		return nil, fmt.Errorf("private key not provided")
	}
	return &oidcService{
		db:         config.MongoClient,
		privateKey: config.PrivateKey,
	}, nil
}

type OAuth2TokenInput struct {
	GrantType    api.OAuth2TokenRequestGrantType
	Scope        *string
	Code         string
	ClientId     string
	ClientSecret string
}

func (s *oidcService) OAuth2Authorize(request api.OAuth2AuthorizationRequest) (redirectUrl *url.URL, err error) {
	url, err := url.Parse(*request.RedirectUri)
	if err != nil {
		return nil, &errors.OAuth2Error{OAuth2Error: api.InvalidRequest}
	}
	q := url.Query()
	q.Add("code", "abc123")
	url.RawQuery = q.Encode()
	return url, nil
}

func (s *oidcService) Oauth2Token(request OAuth2TokenInput) (response *api.OAuth2TokenResponse, err error) {
	switch request.GrantType {
	case api.RefreshToken:
		fallthrough
	case api.ClientCredentials:
		response := api.OAuth2TokenResponse{
			AccessToken:  "1234",
			ExpiresIn:    300,
			TokenType:    api.Bearer,
			RefreshToken: utils.ToAddress("abcd"),
		}
		if request.Scope != nil {
			response.Scope = request.Scope
		}
		return &response, nil
	case api.AuthorizationCode:
		encryptionKey, err := keys.ReadKeyFromFile(constants.AUTHORIZATION_ENCRYPTION_FILENAME)
		if err != nil {
			return nil, err
		}

		result, err := keys.Decrypt(request.Code, encryptionKey)
		if err != nil {
			return nil, &errors.OAuth2Error{
				OAuth2Error: api.InvalidGrant,
			}
		}
		var authorizationData AuthorizationCode
		if err = json.Unmarshal(result, &authorizationData); err != nil {
			return nil, &errors.OAuth2Error{
				OAuth2Error: api.InvalidGrant,
			}
		}

		// TODO: Validate client id and redirect uri

		if authorizationData.Exp.Before(time.Now().UTC()) {
			return nil, &errors.OAuth2Error{
				OAuth2Error: api.InvalidGrant,
			}
		}

		kid := "1"
		token := jwt.NewWithClaims(
			jwt.SigningMethodRS256,
			jwt.StandardClaims{
				Subject:   "anson-user",
				Id:        "1234",
				Issuer:    "authorization.ansonallard.com",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			},
		)
		token.Header["kid"] = kid

		signedToken, err := token.SignedString(s.privateKey)
		if err != nil {
			fmt.Printf("%+v", err)
			return nil, err
		}
		fmt.Printf("jwt: %s\n", signedToken)

		response := api.OAuth2TokenResponse{
			AccessToken: signedToken,
			ExpiresIn:   300,
			TokenType:   api.Bearer,
			// TODO: implement this
			RefreshToken: utils.ToAddress("abcd"),
		}
		if request.Scope != nil {
			response.Scope = request.Scope
		}
		return &response, nil
	case api.Password:
		fallthrough
	default:
		return nil, &errors.OAuth2Error{OAuth2Error: api.InvalidGrant}
	}
}
