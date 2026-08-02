package handler

import (
	"strconv"

	"github.com/davidcm146/bus-booking-be/internal/dto/request"
	"github.com/davidcm146/bus-booking-be/internal/service"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/davidcm146/bus-booking-be/internal/shared/response"
	"github.com/davidcm146/bus-booking-be/internal/shared/validator"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler with the given service.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create godoc
// @Summary      Create a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body  request.CreateUserRequest  true  "Create user payload"
// @Success      201
// @Failure      400
// @Failure      409
// @Router       /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validator.FormatError(c.Request.Context(), err))
		return
	}

	result, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, result)
}

// GetByID godoc
// @Summary      Get user by ID
// @Tags         users
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200
// @Failure      404
// @Router       /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ValidationError(c, i18n.T(c.Request.Context(), constant.MsgKeyInvalidIDFormat))
		return
	}

	result, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// GetAll godoc
// @Summary      Get all users
// @Tags         users
// @Produce      json
// @Success      200
// @Router       /api/v1/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	result, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// Update godoc
// @Summary      Update a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path  string                     true  "User ID"
// @Param        body  body  request.UpdateUserRequest   true  "Update user payload"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      409
// @Router       /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ValidationError(c, i18n.T(c.Request.Context(), constant.MsgKeyInvalidIDFormat))
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, validator.FormatError(c.Request.Context(), err))
		return
	}

	result, err := h.userService.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// Delete godoc
// @Summary      Delete a user
// @Tags         users
// @Param        id  path  string  true  "User ID"
// @Success      204
// @Failure      404
// @Router       /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ValidationError(c, i18n.T(c.Request.Context(), constant.MsgKeyInvalidIDFormat))
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}
