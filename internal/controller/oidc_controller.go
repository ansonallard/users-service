package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/service"
	"github.com/ansonallard/users-service/internal/utils"
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
	r.ParseForm()
	defer r.Body.Close()

	response, err := c.oidcService.Oauth2Token(service.OAuth2TokenInput{
		GrantType:    api.OAuth2TokenRequestGrantType(r.FormValue("grant_type")),
		Scope:        utils.ToAddress(r.FormValue("scope")),
		Code:         r.FormValue("code"),
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
	})

	if err != nil {
		return c.errorHanding(g, err)
	}
	g.JSON(http.StatusOK, response)
	return nil
}

func (c *OidcController) OAuth2Authorize(g *gin.Context) error {
	r := g.Request

	r.ParseForm()
	defer r.Body.Close()

	redirectUri := r.FormValue("redirect_uri")
	clientId := r.FormValue("client_id")
	url, err := c.oidcService.OAuth2Authorize(api.OAuth2AuthorizationRequest{
		ResponseType: api.OAuth2AuthorizationRequestResponseType(r.FormValue("response_type")),
		ClientId:     clientId,
		RedirectUri:  &redirectUri,
		Scope:        utils.ToAddress(r.FormValue("scope")),
	})
	if err != nil {
		return c.errorHanding(g, err)
	}

	var redirect_url string
	// redirect_url = url.String()
	rurl, _ := url.Parse("http://localhost:5000/login")
	q := rurl.Query()
	q.Add("client_id", clientId)
	q.Add("redirect_uri", redirectUri)

	corsUrl, _ := url.Parse(redirectUri)
	rurl.RawQuery = q.Encode()
	redirect_url = rurl.String()

	hostName := corsUrl.Scheme + "://" + corsUrl.Host
	fmt.Println(hostName)
	g.Header("Access-Control-Allow-Origin", hostName)
	g.Redirect(http.StatusFound, redirect_url)
	return nil
}

type myStruct struct{}

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
