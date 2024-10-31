package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// GenerateResetToken generates a 6-digit numerical reset token
func GenerateResetToken() string {
	// Create a new random source seeded with the current time
	randSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 100000 and 999999
	resetToken := randSource.Intn(900000) + 100000

	// Convert the number to a string and return it
	return strconv.Itoa(resetToken)
}

// SendResetTokenEmail sends the reset token to the user's email address
func SendResetTokenEmail(toEmail, resetToken string) error {
	// SMTP server configuration
	smtpHost := "sandbox.smtp.mailtrap.io" 
	smtpPort := "587"               
	smtpUser := "ec22aed451accb" 
	smtpPass := "58db41323b0501"     

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Compose the email message
	subject := "Password Reset Token Requested!"
	body := fmt.Sprintf("Your password reset token is: %s\nPlease use this token to reset your password.", resetToken)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{toEmail}, message)
	if err != nil {
		return err
	}

	return nil
}

// GenerateToken generates a JWT token for a user
func GenerateToken(userID uuid.UUID) (string, error) {
	// Get the secret key from viper
	secretKey := viper.GetString("JWT_SECRET")

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}


// VerifyToken verifies a JWT token and returns the user ID if the token is valid
func VerifyToken(tokenString string) (uuid.UUID, error) {
	// Parse the token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key
		return []byte(viper.GetString("JWT_SECRET")), nil
	})

	// Check if parsing the token failed
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse token: %v", err)
	}

	// Verify if the token is valid
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	// Extract the claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("could not parse token claims")
	}

	// Check if the "exp" claim is still valid
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return uuid.Nil, fmt.Errorf("token has expired")
		}
	} else {
		return uuid.Nil, fmt.Errorf("invalid expiration time")
	}

	// Extract the user ID
	userId, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID")
	}

	// Convert user ID to UUID
	return uuid.Parse(userId)
}
