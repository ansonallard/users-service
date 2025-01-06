package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ansonallard/users-service/internal/api"
	"github.com/ansonallard/users-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UsersController struct {
	service *service.UsersService
}

func NewUsersController(usersService *service.UsersService) UsersController {
	return UsersController{service: usersService}
}

func (u *UsersController) CreateUser(ctx context.Context, g *gin.Context, pathParams map[string]string) error {

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
		if err == fmt.Errorf("error: Tenant '%s' not found", incomingUser.TenantId) {
			g.JSON(http.StatusNotFound, err)
			return nil
		}
		return err
	}
	return nil
}
