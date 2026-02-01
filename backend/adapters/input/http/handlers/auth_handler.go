package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"torque-dms/adapters/input/http/dto/request"
	"torque-dms/adapters/input/http/dto/response"
	"torque-dms/core/identity/ports/input"
)

type AuthHandler struct {
	authService input.AuthService
}

func NewAuthHandler(authService input.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, user, err := h.authService.Register(input.RegisterInput{
		Type:         req.Type,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BusinessName: req.BusinessName,
		Email:        req.Email,
		Phone:        req.Phone,
		Username:     req.Username,
		Password:     req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response.RegisterResponse{
		Entity: response.EntityResponse{
			ID:           entity.ID,
			Type:         string(entity.Type),
			FirstName:    entity.FirstName,
			LastName:     entity.LastName,
			BusinessName: entity.BusinessName,
			Email:        entity.Email,
			Status:       string(entity.Status),
			CreatedAt:    entity.CreatedAt,
		},
		User: response.UserResponse{
			ID:        user.ID,
			EntityID:  user.EntityID,
			Username:  user.Username,
			Status:    string(user.Status),
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.Login(input.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		User: response.UserResponse{
			ID:        result.User.ID,
			EntityID:  result.User.EntityID,
			Username:  result.User.Username,
			LastLogin: result.User.LastLogin,
			Status:    string(result.User.Status),
			CreatedAt: result.User.CreatedAt,
		},
		Token: result.Token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	err := h.authService.Logout(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := h.authService.ChangePassword(input.ChangePasswordInput{
		UserID:      userID.(uint),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}