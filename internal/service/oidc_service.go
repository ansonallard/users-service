package service

import (
	"encoding/json"
	"net/url"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/constants"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/keys"
	"github.com/ansonallard/users-service/internal/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OidcService struct {
	db *mongo.Client
}

func NewOidcService(mongoClient *mongo.Client) *OidcService {
	return &OidcService{
		db: mongoClient,
	}
}

type OAuth2TokenInput struct {
	GrantType    api.OAuth2TokenRequestGrantType
	Scope        *string
	Code         string
	ClientId     string
	ClientSecret string
}

func (s *OidcService) OAuth2Authorize(request api.OAuth2AuthorizationRequest) (redirectUrl *url.URL, err error) {
	url, err := url.Parse(*request.RedirectUri)
	if err != nil {
		return nil, &errors.OAuth2Error{OAuth2Error: api.InvalidRequest}
	}
	q := url.Query()
	q.Add("code", "abc123")
	url.RawQuery = q.Encode()
	return url, nil
}

func (s *OidcService) Oauth2Token(request OAuth2TokenInput) (response *api.OAuth2TokenResponse, err error) {
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
			return nil, err
		}
		var authorizationData AuthorizationCode
		if err = json.Unmarshal(result, &authorizationData); err != nil {
			return nil, err
		}
		response := api.OAuth2TokenResponse{
			AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			ExpiresIn:    300,
			TokenType:    api.Bearer,
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
