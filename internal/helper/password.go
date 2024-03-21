package helper

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Get the salt cost factor from the environment variable
	costFactor := os.Getenv("BCRYPT_SALT")
	if costFactor == "" {
		costFactor = "10" // Default to cost factor of 10 if not provided
	}

	// Convert the cost factor to int
	cost, err := strconv.Atoi(costFactor)
	if err != nil {
		return "", err
	}

	// Generate a salted hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	// Compare the hashed password with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}

func CompareHashedPasswordAndPassword(hashedPassword, password string) error {
	// Compare the hashed password with the provided password
	hashPassword, err := HashPassword(hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err
}
