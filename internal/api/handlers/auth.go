package handlers

import (
	"net/http"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	UserType string `json:"user_type" binding:"required,oneof=admin enduser drone"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// Login generates a JWT token for the user
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.GenerateToken(req.Name, req.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}
