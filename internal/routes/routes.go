package routes

import (
	"auth-service/internal/controller"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(server *gin.Engine) {
	server.POST("/api/v1/signup", controller.Signup) 										// Endpoint for user registration
	server.POST("/api/v1/login", controller.Login) 											// Endpoint for user Login
	server.POST("/api/v1/logout", controller.Logout) 										// Endpoint for user logout
	server.POST("/api/v1/token/refresh", controller.RefreshToken) 			// Endpoint for refresh JWT token
	server.POST("/api/v1/forgot-password", controller.ForgotPassword) 	// Request password reset
	server.POST("/api/v1/reset-password", controller.ResetPassword) 		// Reset Password using token
	server.PUT("/api/v1/change-password", controller.ChangePassword) 		// Change password while logged in
	server.GET("/api/v1/profile", controller.GetProfile) 								// Retrieve the logged in user profile
	server.PUT("/api/v1/profile", controller.UpdateProfile) 						// Update the logged-in user's profil
	server.GET("/api/v1/verify-account", controller.VerifyAccount) 			// Verify user email account
	server.DELETE("/api/v1/delete-account", controller.DeleteAccount) 	// Delete user account
}