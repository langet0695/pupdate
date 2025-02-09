package mail

import (
	"fmt"
	"os"

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
