package controller

import (
	"context"
	"net/http"

	"github.com/ansonallard/users-service/internal/service"
	"github.com/gin-gonic/gin"
)

type TenantsController struct {
	service service.TenantsService
}

func NewTenantsContorller(service service.TenantsService) TenantsController {
	return TenantsController{service: service}
}

func (t *TenantsController) CreateTenant(ctx context.Context, g *gin.Context) error {
	response, err := t.service.Create(ctx)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err)
	}
	g.JSON(http.StatusCreated, response)
	return nil
}
