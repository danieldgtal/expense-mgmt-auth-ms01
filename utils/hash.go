package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// GenSalt generates a random salt for password hashing
func GenSalt() string {
	salt := make([]byte, 16) // 16 bytes for salt
	if _, err := rand.Read(salt); err != nil {
		log.Fatalf("failed to generate salt: %v", err)
	}
	return fmt.Sprintf("%x", salt) // Convert bytes to hexadecimal string
}

// HashPassword hashes the password using the provided salt
func HashPassword(password, salt string) (string, error) {
	// Combine the salt and password
	saltedPassword := salt + password

	// Hash the salted password
	hash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares the stored hash with the hashed combination of the given password and salt.
func CheckPassword(storedHash, salt, givenPassword string) error {
	// Combine the salt and given password
	saltedGivenPassword := salt + givenPassword

	// Hash the given password with the salt
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(saltedGivenPassword))
}


// IsValidEmail checks if the provided email is valid
func IsValidEmail(email string) bool {
	// Use a basic regex for email validation
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}