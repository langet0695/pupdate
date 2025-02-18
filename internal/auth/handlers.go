package auth

import (
	"fmt"
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
	fmt.Println("Generating Token for: ", username)
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

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)

	c.IndentedJSON(http.StatusCreated, tokenString)
}

// Function to verify JWT tokens
func AuthenticateMiddleware(c *gin.Context) {
	// Retrieve the token from the cookie
	// tokenString, err := c.Cookie("token")
	bearerString := c.GetHeader("Authorization")
	tokenString := strings.SplitAfter(bearerString, " ")[1]

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}
	fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

	c.Next()
}
