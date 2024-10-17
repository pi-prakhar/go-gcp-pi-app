package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/auth/models"
	"github.com/pi-prakhar/go-gcp-pi-app/internal/user/services"
	errors "github.com/pi-prakhar/go-gcp-pi-app/pkg/error"
	"github.com/pi-prakhar/go-gcp-pi-app/pkg/utils"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request models.GoogleUser
	var res utils.Responder
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to parse data",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}

	if err = h.Service.CreateUser(c.Request.Context(), request); err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to create user",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}
	res = &utils.SuccessResponse[*models.GoogleUser]{
		Message:    "Successfully!! created user",
		StatusCode: http.StatusOK,
		Data:       &request,
	}

	res.Write(c.Writer)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	var email string = c.Param("mail")
	var res utils.Responder
	var user *models.GoogleUser
	var err error

	user, err = h.Service.GetUserByEmail(c.Request.Context(), email)
	if err == errors.ErrUserNotFound {
		res = &utils.SuccessResponse[any]{
			Message:    "User are not found",
			StatusCode: http.StatusOK,
		}
		res.Write(c.Writer)
		return

	} else if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to get user",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
	}

	res = &utils.SuccessResponse[*models.GoogleUser]{
		Message:    "Successfully!! fetched user",
		StatusCode: http.StatusOK,
		Data:       user,
	}
	res.Write(c.Writer)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	var res utils.Responder
	var users []*models.GoogleUser
	var err error

	users, err = h.Service.GetUsers(c.Request.Context())
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to get users",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}

	if users == nil {
		res = &utils.SuccessResponse[any]{
			Message:    "No users are found",
			StatusCode: http.StatusOK,
		}
		res.Write(c.Writer)
		return
	}

	res = &utils.SuccessResponse[[]*models.GoogleUser]{
		Message:    "Successfully!! fetched all users",
		StatusCode: http.StatusOK,
		Data:       users,
	}
	res.Write(c.Writer)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var request models.GoogleUser
	var updatedUser *models.GoogleUser
	var res utils.Responder
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to parse data",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}

	if updatedUser, err = h.Service.UpdateUser(c.Request.Context(), &request); err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to update user",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}
	res = &utils.SuccessResponse[*models.GoogleUser]{
		Message:    "Successfully!! updated user",
		StatusCode: http.StatusOK,
		Data:       updatedUser,
	}

	res.Write(c.Writer)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	var email string = c.Param("mail")
	var res utils.Responder
	var user *models.GoogleUser
	var err error = h.Service.DeleteUser(c.Request.Context(), email)

	if err == errors.ErrUserNotFound {
		res = &utils.ErrorResponse{
			Message:    "Failed to delete user",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	} else if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to delete user",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(c.Writer)
		return
	}

	res = &utils.SuccessResponse[*models.GoogleUser]{
		Message:    "Successfully!! deleted user",
		StatusCode: http.StatusOK,
		Data:       user,
	}
	res.Write(c.Writer)
}
