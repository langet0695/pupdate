package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/langet/pupdate/internal/auth"
	"github.com/langet/pupdate/internal/mail"
	"github.com/langet/pupdate/internal/subscription"
)

func NewRouter() *gin.Engine {
	// Make mail loop over opt in users
	// Make get subscriver email show all instances of subscriber history
	APP_USER := os.Getenv("APP_USER")
	APP_PASSWORD := os.Getenv("APP_PASSWORD")

	router := gin.Default()
	router.POST("/createToken", gin.BasicAuth(gin.Accounts{APP_USER: APP_PASSWORD}), auth.CreateToken)
	router.GET("/subscriber/:email", subscription.GetSubscriberByEmail)
	router.POST("/subscriber", auth.AuthenticateMiddleware, subscription.CreateSubscriber)
	router.DELETE("/subscriber/:email", auth.AuthenticateMiddleware, subscription.DeleteSubscriber)

	router.GET("/subscribers", auth.AuthenticateMiddleware, subscription.GetActiveSubscribers)
	// router.GET("/subscribers", subscription.GetActiveSubscribers)
	router.POST("/mail", auth.AuthenticateMiddleware, mail.SendMail)
	router.POST("saveSubscribers", auth.AuthenticateMiddleware, mail.BackupSubscriptions)

	// router.Run("localhost:8080")
	router.Run(":8080")
	return router
}
