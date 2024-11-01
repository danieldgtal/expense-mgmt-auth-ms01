package models

import (
	"auth-service/db"
	"auth-service/utils"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id"`
	FirstName         string    `gorm:"not null" json:"first_name"`
	LastName          string    `gorm:"not null" json:"last_name"`
	Email             string    `gorm:"unique;not null" json:"email"`
	PasswordHash      string    `gorm:"not null" json:"password"` // Change json tag to "password"
	Salt              string    `gorm:"not null" json:"-"`
	IsEmailVerified   bool      `gorm:"default:false" json:"is_email_verified"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	ResetPasswordToken string    `gorm:"size:255" json:"-"`
	ResetPasswordExpires time.Time `json:"reset_password_expires"`
	Currency          string    `gorm:"type:char(3);default:CAD;check:currency in ('CAD', 'USD')" json:"currency"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

type ConfirmResetPassword struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdatePassword struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type UpdateUser struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}


// CreateUser saves a new user to the database after validating and hashing the password
func (u *User) CreateUser() error {
	
	// Use GetDBInstance to get the DB instance
	DB:= db.GetDBInstance()
	// Validate the user object
	if err := u.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Generate a new UUID for the user
	u.UserID = uuid.New()

	// Generate a unique salt for the password
	u.Salt = utils.GenSalt()

	// Hash the password with the generated salt
	hashedPassword, err := utils.HashPassword(u.PasswordHash, u.Salt) // Hash the password with the salt
	
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.PasswordHash = hashedPassword
	u.IsEmailVerified = true
	u.CreatedAt = time.Now()
	u.ResetPasswordToken = ""
	u.ResetPasswordExpires = time.Time{}
	u.Currency = "CAD"

	// Save the user to the database using GORM
	if err := DB.Create(u).Error; err != nil {
		return err // Return error if saving to the database fails
	}

	return nil // Return nil if the operation was successful
}

// Validate checks the required fields and constraints
func (u *User) Validate() error {
	// Trim leading and trailing spaces from user input
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	// Check required fields
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("first name and last name cannot be empty")
	}
	if u.Email == "" || !utils.IsValidEmail(u.Email) {
		return errors.New("invalid email address")
	}
	if len(u.PasswordHash) == 0 {
		return errors.New("password cannot be empty")
	}

	// Check for unique email
	var count int64
	if err := db.DB.Model(&User{}).Where("email = ?", u.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already exists")
	}

	return nil
}

// GetUserByEmail retrieves a user by email from the database
func (u *User) GetUserByEmail(email string) error {
	if !utils.IsValidEmail(email) {
		return errors.New("invalid email format")
	}

	DB := db.GetDBInstance()
	return DB.Where("email = ?", email).First(u).Error
}
