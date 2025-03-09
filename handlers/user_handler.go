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

		// Here you would normally generate a JWT token
		// For simplicity, just returning user information
		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user": gin.H{
				"id":       user.ID.Hex(),
				"username": user.Username,
				"email":    user.Email,
				"roles":    user.Roles,
			},
		})
	}
}

// GetUserHandler retrieves user information
func GetUserHandler(userController *controllers.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			return
		}

		user, err := userController.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":        user.ID.Hex(),
				"username":  user.Username,
				"email":     user.Email,
				"roles":     user.Roles,
				"active":    user.Active,
				"lastLogin": user.LastLogin,
				"createdAt": user.CreatedAt,
				"updatedAt": user.UpdatedAt,
			},
		})
	}
}

// UpdateUserRolesHandler updates roles for a user
func UpdateUserRolesHandler(userController *controllers.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var roleRequest struct {
			UserID string   `json:"user_id" binding:"required"`
			Roles  []string `json:"roles" binding:"required"`
		}

		if err := c.ShouldBindJSON(&roleRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userController.UpdateUserRoles(
			c.Request.Context(),
			roleRequest.UserID,
			roleRequest.Roles,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User roles updated successfully",
		})
	}
}
