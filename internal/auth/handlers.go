package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

// Function to create JWT tokens with claims
func CreateToken(c *gin.Context) {
	username := "a-dummy-username"
	fmt.Println("Generating Token for: ", username)
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "pupdate",                        // Issuer
		"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		// return "", err
		fmt.Println(err)
	}

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	// return tokenString, nil
	c.IndentedJSON(http.StatusCreated, tokenString)
}

// Function to verify JWT tokens
func AuthenticateMiddleware(c *gin.Context) {
	// Retrieve the token from the cookie
	// tokenString, err := c.Cookie("token")
	bearerString := c.GetHeader("Authorization")
	tokenString := strings.SplitAfter(bearerString, " ")[1]
	fmt.Println("THIS IS RECIEVED TOKEN: ", tokenString)
	// fmt.Println("Result:", result[1])
	// if err != nil {
	// 	fmt.Println("Token missing in cookie")
	// 	c.Redirect(http.StatusSeeOther, "/login")
	// 	c.Abort()
	// 	return
	// }

	// Verify the token
	token, err := verifyToken(tokenString)
	if err != nil {
		fmt.Printf("Token verification failed: %v\\n", err)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	// Print information about the verified token
	fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

	// Continue with the next middleware or route handler
	c.Next()
}
