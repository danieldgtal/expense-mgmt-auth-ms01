package common

import (
	"auth-service/db"
	"auth-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StoreToken stores the generated token in the database
func StoreToken(userID uuid.UUID, tokenString string, createdTime time.Time, expirationTime time.Time) error {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Create a new AuthToken instance
	token := models.AuthToken{
		UserID:    userID,
		Token:     tokenString,
		CreatedAt: createdTime,
		ExpiresAt: expirationTime,
	}

	// Insert the token into the database
	if err := DB.Create(&token).Error; err != nil {
		return err // Return the error if there's an issue
	}

	return nil // Return nil if the operation was successful
}

// StoreResetToken updates the user's reset token and its expiration time in the database
func StoreResetToken(userID uuid.UUID, tokenString string, expirationTime time.Time) error {
	// Get the DB instance
	DB := db.GetDBInstance()

	// Find the user by userID
	var user models.User
	if err := DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return err // Return the error if user is not found
	}

	// Update the user's reset password token and expiration time
	user.ResetPasswordToken = tokenString
	user.ResetPasswordExpires = expirationTime

	// Save the updated user record
	if err := DB.Save(&user).Error; err != nil {
		return err // Return the error if there's an issue saving
	}

	return nil // Return nil if the operation was successful
}


// DeleteToken removes the specified token from the database
func DeleteToken(tokenString string) error {
	// Get the DB instance
	db := db.GetDBInstance()

	// Attempt to delete the token
	result := db.Where("token = ?", tokenString).Delete(&models.AuthToken{})
	if result.Error != nil {
		return result.Error // Return the error if there's an issue
	}
	
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return a not found error if no rows were affected
	}

	return nil // Return nil if the operation was successful
}

