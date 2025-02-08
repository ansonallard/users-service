package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/errors"
	"github.com/ansonallard/users-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service *service.UsersService
}

func NewUsersController(usersService *service.UsersService) UsersController {
	return UsersController{service: usersService}
}

func (u *UsersController) CreateUser(ctx context.Context, g *gin.Context, pathParams api.PathParams) error {

	var request api.CreateOrUpdateUserJSONRequestBody
	err := g.BindJSON(&request)
	if err != nil {
		g.JSON(http.StatusBadRequest, fmt.Errorf("error: %+v", err))
	}
	incomingUser := service.IncomingUserRequest{
		Username: pathParams["username"],
		Password: request.Password,
		TenantId: pathParams["id"],
	}
	err = u.service.Create(ctx, incomingUser)
	if err != nil {
		switch err.(type) {
		case errors.UserExistsError:
			g.JSON(http.StatusConflict, err)
		case errors.TenantNotFoundError:
			g.JSON(http.StatusNotFound, err)
		default:
			return err
		}
	}
	return nil
}

func (u *UsersController) Login(ctx context.Context, g *gin.Context, pathParams api.PathParams) error {
	defer g.Request.Body.Close()

	var body api.UserLoginRequest
	if err := g.BindJSON(&body); err != nil {
		g.AbortWithStatus(http.StatusBadRequest)
	}
	url, err := url.Parse(*body.RedirectUri)
	if err != nil {
		g.AbortWithStatus(http.StatusInternalServerError)
	}

	tenantId := pathParams["id"]

	// TODO: Pass application id for redirect uri validation
	result, err := u.service.Login(ctx, service.LoginInput{
		Username:    body.Username,
		Password:    body.Password,
		RedirectUri: *body.RedirectUri,
		TenantId:    tenantId,
	})
	if err != nil {
		switch err.(type) {
		case errors.UserNotFoundError:
			g.AbortWithStatus(http.StatusUnauthorized)
			return nil
		case errors.NotAuthorizedError:
			g.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
			return nil
		default:
			return err
		}
	}
	q := url.Query()
	q.Add("code", result.Code)
	url.RawQuery = q.Encode()

	// Add CORS headers explicitly
	// ctx.Header("Access-Control-Allow-Origin", ctx.GetHeader("Origin"))
	// // ctx.Header("Access-Control-Allow-Origin", "https://oauth.pstmn.io/v1/callback")
	// ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	// ctx.Header("Access-Control-Expose-Headers", "Location") // Expose Location header for the redirect

	g.JSON(http.StatusOK, api.UserLoginResponse{RedirectUrl: url.String()})
	return nil
}
