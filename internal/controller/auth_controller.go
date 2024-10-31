package controller

import (
	"auth-service/db"
	"auth-service/internal/common"
	"auth-service/internal/models"
	"auth-service/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Signup handles user registration
func Signup(context *gin.Context) {
	// Parse the incoming JSON request to extract user data
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		// If parsing fails, send a standardized error response
		utils.SendResponse(context, http.StatusBadRequest, "Invalid request data", nil, gin.H{"error": err.Error()})
		return
	}
	
	// Call a model function to save the user
	if err := user.CreateUser(); err != nil {
		// If saving fails, send an error response
		utils.SendResponse(context, http.StatusInternalServerError, "User could not be created", nil, gin.H{"error": err.Error()})
		return
	}
	
	// Send a success response with a simplified user representation (like user ID) or an empty object
	utils.SendResponse(context, http.StatusCreated, "User registered successfully", gin.H{"userId": user.UserID}, nil)
}

// Login handles user login
func Login(context *gin.Context) {
	var loginReq models.LoginRequest
	if err := context.ShouldBindJSON(&loginReq); err != nil {
		utils.SendResponse(context, http.StatusBadRequest, "Invalid request data", nil, gin.H{"error": err.Error()})
		return
	}

	// Password cannot be empty
	if loginReq.Password == "" {
		utils.SendResponse(context, http.StatusBadRequest, "Password cannot be empty", nil, nil)
		return
	}
	
	// Get User Email
	var user models.User
	if err := user.GetUserByEmail(loginReq.Email); err != nil {
		utils.SendResponse(context, http.StatusUnauthorized, "Invalid credentials", nil, nil)
		return
	}

	// Check the password using the CheckPassword function
	if err := utils.CheckPassword(user.PasswordHash, user.Salt, loginReq.Password); err != nil {
		utils.SendResponse(context, http.StatusUnauthorized, "Invalid credentials", nil, nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.UserID) // Assuming user.UserID is a uuid.UUID
	if err != nil {
		utils.SendResponse(context, http.StatusInternalServerError, "Could not generate token", nil, nil)
		return
	}

	// Store Generated token in Auth_Token table
	createdTime := time.Now() // Get the current time for when the token is created
	expirationTime := time.Now().Add(time.Hour * 24) // Make a 24 hour expiration time
	
	err = common.StoreToken(user.UserID, token, createdTime, expirationTime)
	if err != nil {
		utils.SendResponse(context, http.StatusInternalServerError, "Could not store token", nil, nil)
		return
	}

	utils.SendResponse(context, http.StatusOK, "Login successful", gin.H{"userId": user.UserID, "token": token}, nil)
}

// ResetPassword handles the password reset request.
func ResetPassword(c *gin.Context) {
	var resetPassword models.ResetPassword

	if err := c.ShouldBindJSON(&resetPassword); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid request", nil, nil)
		return
	}

	// Check if the user exists
	var user models.User
	if err := user.GetUserByEmail(resetPassword.Email); err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Invalid credentials", nil, nil)
		return
	}

	// Generate a reset token
	token := utils.GenerateResetToken()

	// Store the token in the database using existing function
	expirationTime := time.Now().Add(time.Hour * 1) // Make a 1 hour expiration time
	
	resetToken := common.StoreResetToken(user.UserID, token, expirationTime)
	if resetToken != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not store token", nil, nil)
		return
	}

	// Send email with the reset token (implement email sending function)
	if err := utils.SendResetTokenEmail(resetPassword.Email, token); err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not send reset email", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Reset email sent", gin.H{"status": "success",
    "message": "Please check your email for instructions to reset your password.",}, nil)
}

// ConfirmResetPassword handles the password reset confirmation
func ConfirmResetPassword(c *gin.Context) {
	
	var confirmResetPassword models.ConfirmResetPassword
	// Parse the incoming JSON request to extract token and new password
	if err := c.ShouldBindJSON(&confirmResetPassword); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid request data", nil, gin.H{"error": err.Error()})
		return
	}

	// Get the DB instance
	DB := db.GetDBInstance()

	// Find the user by the reset token and ensure the token hasn't expired
	var user models.User
	if err := DB.Where("reset_password_token = ? AND reset_password_expires > ?", confirmResetPassword.Token, time.Now()).First(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid or expired reset token", nil, nil)
		return
	}

	// Use the retrieved salt to hash the new password
	hashedPassword, err := utils.HashPassword(confirmResetPassword.NewPassword, user.Salt) // Using the user's existing salt here
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not hash password", nil, nil)
		return
	}

	// Update the user's password hash and clear the reset token fields
	user.PasswordHash = hashedPassword
	user.ResetPasswordToken = ""
	user.ResetPasswordExpires = time.Time{} 

	// Save the updated user record
	if err := DB.Save(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not update password", nil, nil)
		return
	}

	// Send a success response
	utils.SendResponse(c, http.StatusOK, "Password successfully reset", nil, nil)
}

// logout handles user logout
func Logout(context *gin.Context) {
	// Retrieve the token and user ID from the context set by the middleware
	tokenString, exists := context.Get("tokenString")
	if !exists {
		utils.SendResponse(context, http.StatusBadRequest, "Token not found in context", nil, gin.H{"error": "Token not found"})
		return
	}

	// Ensure the token is a string
	token, ok := tokenString.(string)
	if !ok || token == "" {
		utils.SendResponse(context, http.StatusBadRequest, "Invalid token format", nil, gin.H{"error": "Invalid token"})
		return
	}

	// Call the DeleteToken function to remove the token
	if err := common.DeleteToken(token); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.SendResponse(context, http.StatusNotFound, "Token not found", nil, nil)
		} else {
			utils.SendResponse(context, http.StatusInternalServerError, "Failed to log out", nil, gin.H{"error": err.Error()})
		}
		return
	}

	utils.SendResponse(context, http.StatusOK, "User logged out successfully", nil, nil)
}

// GetUserInfo retrieves the user's complete information
func GetUserInfo(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Get the user ID from the JWT token in the middleware
	userID := c.MustGet("user_id").(string)

	// Find the user in the database
	var user models.User
	if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
			utils.SendResponse(c, http.StatusNotFound, "User not found", nil, nil)
			return
	}

	// Send response with user information
	utils.SendResponse(c, http.StatusOK, "User information retrieved successfully", gin.H{
			"user_id":           user.UserID,
			"first_name":        user.FirstName,
			"last_name":         user.LastName,
			"email":             user.Email,
			// Include other relevant fields
	}, nil)
}

// UpdatePassword handles the password change request
func UpdatePassword(c *gin.Context) {
	var updatePassword models.UpdatePassword
	// Get the DB instance
	DB := db.GetDBInstance()

	// Bind JSON to the request struct
	if err := c.ShouldBindJSON(&updatePassword); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid request data", nil, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the JWT token in the middleware (assumed to be set)
	userID := c.MustGet("user_id").(string)

	// Find the user in the database
	var user models.User
	if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil, nil)
		return
	}
   
	// Verify the current password
	if err := utils.CheckPassword(user.PasswordHash, user.Salt, updatePassword.CurrentPassword); err != nil {
		utils.SendResponse(c, http.StatusUnauthorized, "Current password is incorrect", nil, nil)
		return
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(updatePassword.NewPassword, user.Salt)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not hash password", nil, nil)
		return
	}

	// Update the user's password in the database
	user.PasswordHash = hashedPassword
	// user.Salt = salt
	if err := DB.Save(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, "Could not update password", nil, nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Password updated successfully", nil, nil)
}

// UpdateUserInfo handles the user info update request
func UpdateUserInfo(c *gin.Context) {
	// Get the DB instance
	DB := db.GetDBInstance()
	var updateUser models.UpdateUser
	// Bind JSON to the request struct
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "Invalid request data", nil, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the JWT token in the middleware (assumed to be set)
	userID := c.MustGet("user_id").(string)

	// Find the user in the database
	var user models.User
	if err := DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		utils.SendResponse(c, http.StatusNotFound, "User not found", nil, nil)
		return
	}

	// Update user fields
	user.FirstName = updateUser.FirstName
	user.LastName = updateUser.LastName
	user.Email = updateUser.Email

	// Save updated user information in the database
	if err := DB.Save(&user).Error; err != nil {
			utils.SendResponse(c, http.StatusInternalServerError, "Could not update user information", nil, nil)
			return
	}

	utils.SendResponse(c, http.StatusOK, "User information updated successfully", nil, nil)
}

