package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"

	"net/http"

	"time"

	"github.com/gin-gonic/gin"

	gomail "gopkg.in/mail.v2"
)

func getSubscribers(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, subscribers)
}

func getSubscriberByEmail(c *gin.Context) {
	email := c.Param("email")

	for _, sub := range subscribers {
		if sub.Email == email {
			c.IndentedJSON(http.StatusOK, sub)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "subscriber not found"})
}

func createSubscriber(c *gin.Context) {
	var newSubscriber subscriber

	if err := c.BindJSON(&newSubscriber); err != nil {
		return
	}

	// TODO Write logic to confirm subscriber doesn't currently exist if so send error.

	subscribers = append(subscribers, newSubscriber)
	c.IndentedJSON(http.StatusCreated, newSubscriber)
}

func deleteSubscriber(c *gin.Context) {
	email := c.Param("email")

	for i, sub := range subscribers {
		if sub.Email == email {
			subscribers[i].DateUnsubscribed = time.Now()
			c.IndentedJSON(http.StatusOK, sub)
			c.IndentedJSON(http.StatusOK, subscribers)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "subscriber not found"})

}

func sendMail(c *gin.Context) {

	MAIL_USER := viperEnvVariable("MAIL_USER")
	MAIL_PASSWORD := viperEnvVariable("MAIL_PASSWORD")

	message := buildMessage(MAIL_USER)
	message.SetHeader("To", "tlange1124@gmail.com")

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

func getPuppy(fileName string) string {
	URL := "https://dog.ceo/api/breeds/image/random"
	TARGET_FOLDER := "/tmp"

	res, getErr := http.Get(URL)

	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := io.ReadAll(res.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}
	dog_obj := dog{}
	jsonErr := json.Unmarshal(body, &dog_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	err, filePath := downloadImages(dog_obj.Address, fileName, TARGET_FOLDER)
	if err != nil {
		panic(err.Error())
	}
	return filePath

}

func getQuote() quote {
	URL := "https://quoteslate.vercel.app/api/quotes/random"

	res, getErr := http.Get(URL)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	quote_obj := quote{}
	jsonErr := json.Unmarshal(body, &quote_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return quote_obj
}

func downloadImages(link string, fileName string, targetFolder string) (err error, dogFilePath string) {
	tmp := getFilePath(targetFolder)
	dogFilePath = filepath.Join(tmp, fileName)

	res, err := http.Get(link)
	if err != nil {
		return
	}
	defer res.Body.Close()

	file, err := os.Create(dogFilePath)
	if err != nil {
		return
	}
	defer file.Close()
	file.ReadFrom(res.Body)

	return
}

func getFilePath(tmpPath string) string {
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := fmt.Sprintf("%s%s", curDir, tmpPath)
	return path
}

func fetchSubject() string {
	var subjectBase string

	if rand.IntN(100)%5 < 4 {
		subjectBase = subjectStandard[rand.IntN(len(subjectStandard))]
	} else {
		subjectBase = subjectPuns[rand.IntN(len(subjectPuns))]
	}

	currentTime := time.Now().Local().Format("2006-01-02")

	subject := fmt.Sprintf(`%s - %s`, subjectBase, currentTime)
	return subject
}

func buildMessage(fromAddress string) *gomail.Message {
	FILE_NAME := "daily_dog.jpg"

	daily_quote := getQuote()
	filePath := getPuppy(FILE_NAME)

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", fromAddress)
	subject := fetchSubject()
	message.SetHeader("Subject", subject)

	// Set email body
	body := fmt.Sprintf(`<p style="font-size: 16px;"><em>%s<br/>   - %s</em></p><img src="cid:daily_dog.jpg" alt="A good looking dog"/>`, daily_quote.Quote, daily_quote.Author)
	message.Embed(filePath)
	message.SetBody("text/html", body)

	return message
}
