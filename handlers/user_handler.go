package handlers

import (
	"net/http"

	controllers "iam_backend/jwork"

	"github.com/gin-gonic/gin"
)

// RegisterHandler handles user registration
func RegisterHandler(userController *controllers.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var registrationRequest struct {
			Username string `json:"username" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}

		if err := c.ShouldBindJSON(&registrationRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userController.RegisterUser(
			c.Request.Context(),
			registrationRequest.Username,
			registrationRequest.Email,
			registrationRequest.Password,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User registered successfully",
			"user_id": user.ID.Hex(),
		})
	}
}

// LoginHandler handles user authentication
func LoginHandler(userController *controllers.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userController.AuthenticateUser(
			c.Request.Context(),
			loginRequest.Username,
			loginRequest.Password,
		)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user_id": user.ID,
		})
	}
}
