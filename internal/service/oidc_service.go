package service

import (
	"net/url"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/utils"
)

type OidcService struct{}

func NewOidcService() *OidcService {
	return &OidcService{}
}

type OAuth2TokenInput struct {
	GrantType    api.OAuth2TokenRequestGrantType
	Scope        *string
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
