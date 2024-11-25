package service

import (
	"ansonallard/users-service/src/api"
	"fmt"
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
	GrantType    string
	Scope        string
	ClientId     string
	ClientSecret string
}

func (s *OidcService) Oauth2ClientCredentials(request Input) (response *api.OAuth2TokenResponse, err error) {
	switch request.GrantType {
	case "client_credentials":
		return &api.OAuth2TokenResponse{
			AccessToken: "1234",
			ExpiresIn:   300,
			TokenType:   api.Bearer,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported grant type '%s'", request.GrantType)
	}
}
