package utils

import "golang.org/x/crypto/bcrypt"

// Hash Password generates a bcrypt hash for the given password
func HashPassword(password string) (string, error) {
	// Generate a hashed password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify Password checks if the given password matches the hashed password

func VerifyPassword(password string, hashedPassword string) bool {
	// Compare the password with its hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
