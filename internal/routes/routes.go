package routes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/langet/pupdate/internal/auth"
	"github.com/langet/pupdate/internal/mail"
	"github.com/langet/pupdate/internal/subscription"
)

func NewRouter() *gin.Engine {
	APP_USER := os.Getenv("APP_USER")
	fmt.Println("APP_USER: ", APP_USER)
	APP_PASSWORD := os.Getenv("APP_PASSWORD")
	// APP_USER := "admin"
	// APP_PASSWORD := "pass"

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/createToken", gin.BasicAuth(gin.Accounts{APP_USER: APP_PASSWORD}), auth.CreateToken)
	router.GET("/subscriber/:email", subscription.GetSubscriberByEmail)
	router.POST("/subscriber", auth.AuthenticateMiddleware, subscription.CreateSubscriber)
	router.DELETE("/subscriber/:email", auth.AuthenticateMiddleware, subscription.DeleteSubscriber)
	// router.POST("/subscriber", subscription.CreateSubscriber)
	// router.DELETE("/subscriber/:email", subscription.DeleteSubscriber)
	router.GET("/subscribers", auth.AuthenticateMiddleware, subscription.GetActiveSubscribers)
	router.POST("/mail", auth.AuthenticateMiddleware, mail.SendMail)
	router.POST("/saveSubscribers", auth.AuthenticateMiddleware, mail.BackupSubscriptions)
	// router.POST("/saveSubscribers", mail.BackupSubscriptions)

	router.Run(":8080")
	return router
}
