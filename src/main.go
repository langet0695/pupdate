package main

import (
	"github.com/gin-gonic/gin"
)

// type Todo struct {
// 	Text string
// 	Done bool
// }

// var todos []Todo
// var loggedInUser string
// var secretKey = []byte("your-secret-key")

// // Function to create JWT tokens with claims
// func createToken(c *gin.Context) {
// 	username := "a-dummy-username"
// 	fmt.Println("Generating Token for: ", username)
// 	// Create a new JWT token with claims
// 	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": username,                         // Subject (user identifier)
// 		"iss": "pupdate",                        // Issuer
// 		"aud": getRole(username),                // Audience (user role)
// 		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
// 		"iat": time.Now().Unix(),                // Issued at
// 	})

// 	tokenString, err := claims.SignedString(secretKey)
// 	if err != nil {
// 		// return "", err
// 		fmt.Println(err)
// 	}

// 	// Print information about the created token
// 	fmt.Printf("Token claims added: %+v\n", claims)
// 	// return tokenString, nil
// 	c.IndentedJSON(http.StatusCreated, tokenString)
// }

// func getRole(username string) string {
// 	if username == "senior" {
// 		return "senior"
// 	}
// 	return "employee"
// }

// // Function to verify JWT tokens
// func authenticateMiddleware(c *gin.Context) {
// 	// Retrieve the token from the cookie
// 	// tokenString, err := c.Cookie("token")
// 	bearerString := c.GetHeader("Authorization")
// 	tokenString := strings.SplitAfter(bearerString, " ")[1]
// 	fmt.Println("THIS IS RECIEVED TOKEN: ", tokenString)
// 	// fmt.Println("Result:", result[1])
// 	// if err != nil {
// 	// 	fmt.Println("Token missing in cookie")
// 	// 	c.Redirect(http.StatusSeeOther, "/login")
// 	// 	c.Abort()
// 	// 	return
// 	// }

// 	// Verify the token
// 	token, err := verifyToken(tokenString)
// 	if err != nil {
// 		fmt.Printf("Token verification failed: %v\\n", err)
// 		c.Redirect(http.StatusSeeOther, "/login")
// 		c.Abort()
// 		return
// 	}

// 	// Print information about the verified token
// 	fmt.Printf("Token verified successfully. Claims: %+v\\n", token.Claims)

// 	// Continue with the next middleware or route handler
// 	c.Next()
// }

// // Function to verify JWT tokens
// func verifyToken(tokenString string) (*jwt.Token, error) {
// 	// Parse the token with the secret key
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	// Check for verification errors
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the token is valid
// 	if !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	// Return the verified token
// 	return token, nil
// }

func main() {
	// Make mail loop over opt in users
	// Make get subscriver email show all instances of subscriber history
	router := gin.Default()
	// gin.BasicAuth(gin.Accounts{"admin": "aPassword"})
	router.POST("/createToken", gin.BasicAuth(gin.Accounts{"admin": "aPassword"}), createToken)
	router.GET("/subscriber/:email", getSubscriberByEmail)
	router.POST("/subscriber", authenticateMiddleware, createSubscriber)
	router.DELETE("/subscriber/:email", authenticateMiddleware, deleteSubscriber)

	router.GET("/subscribers", authenticateMiddleware, getActiveSubscribers)
	router.POST("/mail", sendMail)

	// router.Run("localhost:8080")
	router.Run(":8080")
}
