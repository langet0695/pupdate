package subscription

import (
	"log/slog"
	"os"

	"net/http"

	"time"

	"github.com/gin-gonic/gin"

	"github.com/langet/pupdate/internal"
)

func GetActiveSubscribers(c *gin.Context) {
	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	subscribers := internal.FetchSubscribers(file_path, true)
	c.IndentedJSON(http.StatusOK, subscribers)
}

func GetSubscriberByEmail(c *gin.Context) {
	var subscriberHistory Subscribers
	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	subscribers := internal.FetchSubscribers(file_path, false)

	email := c.Param("email")

	for _, subscriber := range subscribers {
		if subscriber.Email == email {
			subscriberHistory = append(subscriberHistory, subscriber)
		}
	}
	if len(subscriberHistory) > 0 {
		c.IndentedJSON(http.StatusOK, subscriberHistory)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "subscriber not found"})
}

func CreateSubscriber(c *gin.Context) {
	var newSubscriber subscriber
	var email email
	if err := c.BindJSON(&email); err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err)
	}
	slog.Info("EMAIL FROM REQUEST: ", "Email", email)
	newSubscriber.Email = email.Email
	newSubscriber.DateSubscribed = time.Now()
	newSubscriber.DateUnsubscribed = time.Time{}

	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	outSubscriber, err := editSubscribers(file_path, newSubscriber, "create")
	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	if outSubscriber.DateSubscribed == outSubscriber.DateUnsubscribed {
		subscriptionNotificaitons(newSubscriber.Email, "new_subscription")
	} else {
		subscriptionNotificaitons(newSubscriber.Email, "existing_subscription")
	}

	c.IndentedJSON(http.StatusCreated, outSubscriber)
}

func DeleteSubscriber(c *gin.Context) {
	email := c.Param("email")
	var deleteSubscriber subscriber
	deleteSubscriber.Email = email

	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	outSubscriber, err := editSubscribers(file_path, deleteSubscriber, "delete")
	if err != nil {
		slog.Error(err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "subscriber not found"})
	}
	subscriptionNotificaitons(deleteSubscriber.Email, "unsubscribing")

	c.IndentedJSON(http.StatusOK, outSubscriber)

}
