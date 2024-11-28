package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/ansonallard/users-service/src/api"
	"github.com/ansonallard/users-service/src/operations"
	"github.com/ansonallard/users-service/src/service"
	"github.com/ansonallard/users-service/utils"
)

type OidcController struct {
	oidcService *service.OidcService
}

func NewOidcController(oidcService *service.OidcService) *OidcController {
	return &OidcController{oidcService: oidcService}
}

func (c *OidcController) OAuth2Token(r *http.Request) (*api.OAuth2TokenResponse, error) {

	clientId, clientSecret, _ := c.parseBasicAuthorizationHeader(r)
	fmt.Printf("clientId: %s, clientSecret: %s\n", *clientId, *clientSecret)
	// body := []byte{}
	// n, err := io.ReadAll(r.Body)
	r.ParseForm()
	defer r.Body.Close()

	return c.oidcService.Oauth2Token(service.Input{
		GrantType:    api.OAuth2TokenRequestGrantType(r.FormValue("grant_type")),
		Scope:        utils.ToAddress(r.FormValue("scope")),
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
	})

}

func (c *OidcController) parseBasicAuthorizationHeader(r *http.Request) (clientId, clientSecret *string, err error) {
	basicHeaderList := r.Header["Authorization"]
	if len(basicHeaderList) != 1 {
		return nil, nil, fmt.Errorf(operations.Invalid_request)
	}
	basicHeader := basicHeaderList[0]
	basicAuthorizationRegex := regexp.MustCompile("^Basic (?P<EncodedHeader>.*)$")

	results := basicAuthorizationRegex.FindStringSubmatch(basicHeader)
	if len(results) != 2 {
		return nil, nil, fmt.Errorf(operations.Invalid_request)
	}
	decodedHeaderBytes, _ := base64.StdEncoding.DecodeString(results[1])
	headerStr := strings.Split(string(decodedHeaderBytes), ":")
	if len(headerStr) != 2 {
		return nil, nil, fmt.Errorf(operations.Invalid_request)
	}
	clientId, clientSecret = &headerStr[0], &headerStr[1]
	return clientId, clientSecret, nil
}
