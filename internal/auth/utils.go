package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func getRole(username string) string {
	if username == "admin" {
		return "admin"
	}
	return "other"
}

// Function to verify JWT tokens
func verifyToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
