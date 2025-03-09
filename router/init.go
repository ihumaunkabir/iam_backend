package router

import (
	"iam_backend/handlers"
	controllers "iam_backend/jwork"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the routes for the application
func SetupRouter(userController *controllers.UserController) *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/register", handlers.RegisterHandler(userController))
		public.POST("/login", handlers.LoginHandler(userController))
	}

	// // Protected routes (would require authentication middleware)
	// protected := r.Group("/api/v1/protected")
	// {
	// 	protected.PUT("/user/roles", handlers.UpdateUserRolesHandler(userController))
	// }

	return r
}
