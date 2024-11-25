package controller

import (
	"ansonallard/users-service/src/operations"
	"ansonallard/users-service/utils"
	"net/http"
)

type BaseController struct {
	oidcController *OidcController
}

type ProcessRequestOpts struct {
	OpenAPIOperationId string
	Writer             http.ResponseWriter
	Request            *http.Request
}

type BaseControllerOpts struct {
	OidcController *OidcController
}

func NewBaseController(opts BaseControllerOpts) *BaseController {
	return &BaseController{oidcController: opts.OidcController}
}

func (c *BaseController) ProcessRequest(opts ProcessRequestOpts) {
	var (
		err      error
		response any
	)
	switch opts.OpenAPIOperationId {
	case operations.OAUTH2_TOKEN:
		response, err = c.oidcController.OAuth2Token(opts.Request)
	default:
		utils.WriteResponse(opts.Writer, nil, http.StatusNotImplemented, nil)
		return
	}

	if err != nil {
		utils.WriteResponse(opts.Writer, nil, http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(opts.Writer, nil, http.StatusOK, response)
}
