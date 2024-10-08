package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/services"
)

type UserHandler struct {
	services *services.UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {}

func (h *UserHandler) GetUser(c *gin.Context) {}
