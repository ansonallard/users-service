package controller

import (
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"

	"github.com/ansonallard/users-service/src/api"
	"github.com/ansonallard/users-service/src/errors"
	"github.com/ansonallard/users-service/src/service"
	"github.com/ansonallard/users-service/utils"
	"github.com/gin-gonic/gin"
)

type OidcController struct {
	oidcService *service.OidcService
}

func NewOidcController(oidcService *service.OidcService) *OidcController {
	return &OidcController{oidcService: oidcService}
}

func (c *OidcController) OAuth2Token(g *gin.Context) error {

	r := g.Request
	clientId, clientSecret, err := c.parseBasicAuthorizationHeader(r)
	if err != nil {
		return c.errorHanding(g, err)
	}
	// fmt.Printf("clientId: %s, clientSecret: %s\n", *clientId, *clientSecret)
	// body := []byte{}
	// n, err := io.ReadAll(r.Body)
	r.ParseForm()
	defer r.Body.Close()

	response, err := c.oidcService.Oauth2Token(service.Input{
		GrantType:    api.OAuth2TokenRequestGrantType(r.FormValue("grant_type")),
		Scope:        utils.ToAddress(r.FormValue("scope")),
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
	})

	if err != nil {
		return c.errorHanding(g, err)
	}
	g.JSON(http.StatusOK, response)
	return nil
}

func (c *OidcController) errorHanding(g *gin.Context, err error) error {
	if _, ok := err.(*errors.OAuth2Error); ok {
		g.JSON(http.StatusBadRequest, err)
		return nil
	}
	return err
}

func (c *OidcController) parseBasicAuthorizationHeader(r *http.Request) (clientId, clientSecret *string, err error) {
	basicHeaderList := r.Header["Authorization"]
	if len(basicHeaderList) != 1 {
		return nil, nil, &errors.OAuth2Error{OAuth2Error: api.InvalidRequest}
	}
	basicHeader := basicHeaderList[0]
	basicAuthorizationRegex := regexp.MustCompile("^Basic (?P<EncodedHeader>.*)$")

	results := basicAuthorizationRegex.FindStringSubmatch(basicHeader)
	if len(results) != 2 {
		return nil, nil, &errors.OAuth2Error{OAuth2Error: api.InvalidRequest}
	}
	decodedHeaderBytes, _ := base64.StdEncoding.DecodeString(results[1])
	headerStr := strings.Split(string(decodedHeaderBytes), ":")
	if len(headerStr) != 2 {
		return nil, nil, &errors.OAuth2Error{OAuth2Error: api.InvalidRequest}
	}
	clientId, clientSecret = &headerStr[0], &headerStr[1]
	return clientId, clientSecret, nil
}
