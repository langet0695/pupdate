package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand/v2"

	"net/http"

	"time"

	"github.com/gin-gonic/gin"

	gomail "gopkg.in/mail.v2"
)

func getActiveSubscribers(c *gin.Context) {
	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := getFilePath(SUBSCRIBERS_PATH)
	subscribers := fetchSubscribers(file_path, true)
	c.IndentedJSON(http.StatusOK, subscribers)
}

func getSubscriberByEmail(c *gin.Context) {
	var subscriberHistory Subscribers
	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := getFilePath(SUBSCRIBERS_PATH)
	subscribers := fetchSubscribers(file_path, false)

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

func createSubscriber(c *gin.Context) {
	var newSubscriber subscriber
	var email email
	if err := c.BindJSON(&email); err != nil {
		return
	}
	fmt.Println("EMAIL FROM REQUEST: ", email.Email)
	newSubscriber.Email = email.Email
	newSubscriber.DateSubscribed = time.Now()
	newSubscriber.DateUnsubscribed = time.Time{}

	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := getFilePath(SUBSCRIBERS_PATH)
	outSubscriber, err := editSubscribers(file_path, newSubscriber, "create")
	if err != nil {
		fmt.Println(err)
	}

	c.IndentedJSON(http.StatusCreated, outSubscriber)
}

func deleteSubscriber(c *gin.Context) {
	email := c.Param("email")
	var deleteSubscriber subscriber
	deleteSubscriber.Email = email

	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := getFilePath(SUBSCRIBERS_PATH)
	outSubscriber, err := editSubscribers(file_path, deleteSubscriber, "delete")
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "subscriber not found"})
	}

	c.IndentedJSON(http.StatusOK, outSubscriber)

}

func sendMail(c *gin.Context) {

	MAIL_USER := viperEnvVariable("MAIL_USER", ".env")
	MAIL_PASSWORD := viperEnvVariable("MAIL_PASSWORD", ".env")

	message := buildMessage(MAIL_USER)
	// message.SetHeader("To", "tlange1124@gmail.com")
	SUBSCRIBERS_PATH := "/tmp/subscriptions.json"
	file_path := getFilePath(SUBSCRIBERS_PATH)
	subscribers := fetchSubscribers(file_path, true)

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
	// if err := dialer.DialAndSend(message); err != nil {
	// 	fmt.Println("Error:", err)
	// 	panic(err)
	// } else {
	// 	fmt.Println("Email sent successfully!")
	// }
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

	filePath, err := downloadImages(dog_obj.Address, fileName, TARGET_FOLDER)
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

// func downloadImages(link string, fileName string, targetFolder string) (dogFilePath string, err error) {
// 	tmp := getFilePath(targetFolder)
// 	dogFilePath = filepath.Join(tmp, fileName)

// 	res, err := http.Get(link)
// 	if err != nil {
// 		return
// 	}
// 	defer res.Body.Close()

// 	file, err := os.Create(dogFilePath)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()
// 	file.ReadFrom(res.Body)

// 	return
// }

// func getFilePath(tmpPath string) string {
// 	curDir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	path := fmt.Sprintf("%s%s", curDir, tmpPath)
// 	return path
// }

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

// func fetchSubscribers(path string, active bool) Subscribers {

// 	jsonFile, err := lockedfile.Open(path)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	defer jsonFile.Close()

// 	byteValue, _ := io.ReadAll(jsonFile)
// 	fmt.Println("READ LOCK")
// 	var subscribers Subscribers
// 	var outSubscribers Subscribers

// 	json.Unmarshal(byteValue, &subscribers)
// 	if active {
// 		for _, subscriber := range subscribers {
// 			if subscriber.DateUnsubscribed.IsZero() {
// 				outSubscribers = append(outSubscribers, subscriber)
// 			}
// 		}
// 		return outSubscribers
// 	}
// 	return subscribers
// }

// func editSubscribers(path string, subscriber subscriber, action string) (outSubscriber subscriber, err error) {
// 	// Open the file in edit mode which gives us a write lock
// 	fmt.Println("OPENING IN EDIT MODE WRITE LOCK")
// 	jsonFile, err := lockedfile.Edit(path)
// 	if err != nil {
// 		return
// 	}
// 	defer jsonFile.Close()
// 	byteValue, _ := io.ReadAll(jsonFile)

// 	var subscribers Subscribers
// 	json.Unmarshal(byteValue, &subscribers)

// 	// Execute logic to handle create and deletes
// 	if action == "delete" {
// 		for i := 0; i < len(subscribers); i++ {
// 			if subscriber.Email == subscribers[i].Email && subscribers[i].DateUnsubscribed.IsZero() {
// 				subscribers[i].DateUnsubscribed = time.Now()
// 				outSubscriber = subscribers[i]
// 			}
// 		}
// 	}
// 	if action == "create" {
// 		var exists bool = false
// 		for i := 0; i < len(subscribers); i++ {
// 			if subscriber.Email == subscribers[i].Email && subscribers[i].DateUnsubscribed.IsZero() {
// 				exists = true
// 				outSubscriber = subscribers[i]
// 			}
// 		}
// 		if !exists {
// 			subscribers = append(subscribers, subscriber)
// 			outSubscriber = subscriber
// 		}
// 	}
// 	// Construct the byte slice to and write it to destination. Make sure to truncate to remove excess
// 	out, err := json.Marshal(subscribers)
// 	if err != nil {
// 		return
// 	}
// 	bytes, err := jsonFile.WriteAt(out, 0)
// 	if err != nil {
// 		return
// 	}
// 	err = jsonFile.Truncate(int64(len(out)))
// 	if err != nil {
// 		return
// 	}
// 	fmt.Println("BYTES WRITTEN: ", bytes)
// 	return

// }
