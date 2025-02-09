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
	MAIL_PASSWORD := os.Getenv("MAIL_PASSWORD")
	// MAIL_USER := viperEnvVariable("MAIL_USER", ".env")
	// MAIL_PASSWORD := viperEnvVariable("MAIL_PASSWORD", ".env")

	message := buildMessage(MAIL_USER)
	// message.SetHeader("To", "tlange1124@gmail.com")
	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	subscribers := internal.FetchSubscribers(file_path, true)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAIL_USER, MAIL_PASSWORD)
	// Send the email
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
	// MAIL_USER := auth.viperEnvVariable("MAIL_USER", ".env")
	// MAIL_PASSWORD := viperEnvVariable("MAIL_PASSWORD", ".env")

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", MAIL_USER)
	message.SetHeader("To", MAIL_USER)

	timeFormat := "2006-01-02 15:04:05.999999999 -0700 PST"
	subject := fmt.Sprintf("Subscribers Save %s", time.Now().Format(timeFormat))
	message.SetHeader("Subject", subject)

	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := internal.GetFilePath(SUBSCRIBERS_PATH)
	// Set email body
	message.Attach(file_path)
	message.SetBody("text", "Subscriber data incase of recovery")

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAIL_USER, MAIL_PASSWORD)
	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
