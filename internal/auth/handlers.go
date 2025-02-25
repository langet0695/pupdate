package auth

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

// Function to create JWT tokens with claims
func CreateToken(c *gin.Context) {
	username := "a-dummy-username"
	slog.Info("Generating Token for: ", "username", username)
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"iss": "pupdate",
		"aud": getRole(username),
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	c.IndentedJSON(http.StatusCreated, tokenString)
}

// Function to verify JWT tokens
func AuthenticateMiddleware(c *gin.Context) {
	bearerString := c.GetHeader("Authorization")
	tokenString := strings.SplitAfter(bearerString, " ")[1]

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})
		c.Abort()
		return
	}
	slog.Info("Token verified: ", "Claims", token)

	c.Next()
}
