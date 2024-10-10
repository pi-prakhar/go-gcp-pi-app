package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/services"
)

type UserHandler struct {
	Service *services.UserService
}

// func (h *UserHandler) CreateUser(c *gin.Context) {}

// func (h *UserHandler) GetUser(c *gin.Context) {}

// Handler to create a new Google user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var request models.GoogleUser

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CreateUser(c.Request.Context(), request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

// Handler to get a user by email
func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("mail")

	user, err := h.Service.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
