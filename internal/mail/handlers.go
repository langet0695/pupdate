package mail

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/langet/pupdate/internal"
	gomail "gopkg.in/mail.v2"
)

func SendMail(c *gin.Context) {

	MAIL_USER := os.Getenv("MAIL_USER")
	fmt.Println("MAIL_USER: ", MAIL_USER)
	MAIL_PASSWORD := os.Getenv("MAIL_PASSWORD")

	message := buildMessage(MAIL_USER)
	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	subscribers := internal.FetchSubscribers(file_path, true)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAIL_USER, MAIL_PASSWORD)

	for _, subscriber := range subscribers {
		message.SetHeader("To", subscriber.Email)
		if err := dialer.DialAndSend(message); err != nil {
			fmt.Println("Error:", err)
			panic(err)
		} else {
			fmt.Println("Email sent successfully!")
		}
	}
}

func BackupSubscriptions(c *gin.Context) {

	MAIL_USER := os.Getenv("MAIL_USER")
	MAIL_PASSWORD := os.Getenv("MAIL_PASSWORD")

	message := gomail.NewMessage()
	message.SetHeader("From", MAIL_USER)
	message.SetHeader("To", MAIL_USER)

	timeFormat := "2006-01-02 15:04:05.999999999 -0700 PST"
	subject := fmt.Sprintf("Subscribers Save %s", time.Now().Format(timeFormat))
	message.SetHeader("Subject", subject)

	SUBSCRIBERS_PATH := os.Getenv("SUBSCRIPTIONS_PATH")
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	message.Attach(file_path)

	message.SetBody("text", "Subscriber data incase of recovery")

	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAIL_USER, MAIL_PASSWORD)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
