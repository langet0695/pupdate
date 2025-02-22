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
	API_USER := os.Getenv("API_USER")
	fmt.Println("API_USER: ", API_USER)
	API_PASSWORD := os.Getenv("API_PASSWORD")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/createToken", gin.BasicAuth(gin.Accounts{API_USER: API_PASSWORD}), auth.CreateToken)
	router.GET("/subscriber/:email", subscription.GetSubscriberByEmail)
	router.POST("/subscriber", auth.AuthenticateMiddleware, subscription.CreateSubscriber)
	router.DELETE("/subscriber/:email", auth.AuthenticateMiddleware, subscription.DeleteSubscriber)
	router.GET("/subscribers", auth.AuthenticateMiddleware, subscription.GetActiveSubscribers)
	router.POST("/mail", auth.AuthenticateMiddleware, mail.SendMail)
	router.POST("/saveSubscribers", auth.AuthenticateMiddleware, mail.BackupSubscriptions)

	router.Run(":8080")
	return router
}
