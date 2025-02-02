package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	// Make mail loop over opt in users
	// Make get subscriver email show all instances of subscriber history
	router := gin.Default()
	router.GET("/subscriber/:email", getSubscriberByEmail)
	router.POST("/subscriber", createSubscriber)
	router.DELETE("/subscriber/:email", deleteSubscriber)

	router.GET("/subscribers", getActiveSubscribers)
	router.POST("/mail", sendMail)

	router.Run("localhost:8080")
	// router.Run(":8080")
}
