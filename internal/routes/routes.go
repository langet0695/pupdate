package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/langet/pupdate/internal/auth"
	"github.com/langet/pupdate/internal/mail"
	"github.com/langet/pupdate/internal/subscription"
)

func NewRouter() *gin.Engine {
	// Make mail loop over opt in users
	// Make get subscriver email show all instances of subscriber history
	router := gin.Default()
	router.POST("/createToken", gin.BasicAuth(gin.Accounts{"admin": "aPassword"}), auth.CreateToken)
	router.GET("/subscriber/:email", subscription.GetSubscriberByEmail)
	router.POST("/subscriber", auth.AuthenticateMiddleware, subscription.CreateSubscriber)
	router.DELETE("/subscriber/:email", auth.AuthenticateMiddleware, subscription.DeleteSubscriber)

	// router.GET("/subscribers", auth.AuthenticateMiddleware, subscription.GetActiveSubscribers)
	router.GET("/subscribers", subscription.GetActiveSubscribers)
	router.POST("/mail", mail.SendMail)

	// router.Run("localhost:8080")
	router.Run(":8080")
	return router
}
