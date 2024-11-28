package service

import (
	"github.com/ansonallard/users-service/src/api"
	"github.com/ansonallard/users-service/src/errors"
	"github.com/ansonallard/users-service/utils"
)

type OidcService struct{}

func NewOidcService() *OidcService {
	return &OidcService{}
}

// func (s *OidcService) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	requestData := api.CreateUserRequestDto{}
// 	err := json.NewDecoder(r.Body).Decode(&requestData)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	response := api.CreateUserResponseDto{
// 		Id: "123",
// 	}

// 	message := "error"
// 	switch {
// 	case requestData.Username == "bad":
// 		WriteResponse(w, nil, http.StatusBadRequest, api.BadRequest{Message: &message})
// 	case requestData.Username == "unauthorized":
// 		WriteResponse(w, nil, http.StatusUnauthorized, api.UnAuthorized{Message: &message})
// 	case requestData.Username == "forbidden":
// 		WriteResponse(w, nil, http.StatusForbidden, api.Forbidden{Message: &message})
// 	case requestData.Username == "conflict":
// 		WriteResponse(w, nil, http.StatusConflict, api.Conflict{Message: &message})
// 	case requestData.Username == "error":
// 		WriteResponse(w, nil, http.StatusInternalServerError, api.InternalServerError{Message: &message})
// 	default:
// 		WriteResponse(w, nil, http.StatusCreated, response)
// 	}
// }

type Input struct {
	GrantType    api.OAuth2TokenRequestGrantType
	Scope        *string
	ClientId     string
	ClientSecret string
}

func (s *OidcService) Oauth2Token(request Input) (response *api.OAuth2TokenResponse, err error) {
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
		fallthrough
	case api.Password:
		fallthrough
	default:
		return nil, &errors.OAuth2Error{OAuth2Error: api.InvalidGrant}
	}
}
