package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/subscriber/:email", getSubscriberByEmail)
	router.POST("/subscriber", createSubscriber)
	router.DELETE("/subscriber/:email", deleteSubscriber)

	router.GET("/subscribers", getSubscribers)
	router.POST("/mail", sendMail)

	router.Run("localhost:8080")
}
