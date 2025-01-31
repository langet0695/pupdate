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

	// Loop through the list of subscribers, looking for
	// an aemaillbum whose ID value matches the parameter.
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

	// Write logic to confirm subscriber doesn't currently exist if so send error.

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
	// Create a new message
	getPuppy()
	daily_quote := getQuote()
	fmt.Println("A Quote for Mail", daily_quote)
	MAIL_USER := viperEnvVariable("MAIL_USER")
	MAIL_PASSWORD := viperEnvVariable("MAIL_PASSWORD")

	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", MAIL_USER)
	message.SetHeader("To", "tlange1124@gmail.com")
	subject := fetchSubject()
	message.SetHeader("Subject", subject)

	// Set email body
	body := fmt.Sprintf(`<p style="font-size: 16px;"><em>%s<br/>   - %s</em></p><img src="cid:daily_dog.jpg" alt="A good looking dog"/>`, daily_quote.Quote, daily_quote.Author)
	message.Embed("/Users/tlange/pupdate/tmp/daily_dog.jpg")
	message.SetBody("text/html", body)

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

func getPuppy() {
	// Execute fetch call to https://dog.ceo/api/breeds/image/random Fetch!
	fmt.Printf("FETCHING A DOGGIE!")
	url := "https://dog.ceo/api/breeds/image/random"
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println(string(body))
	dog_obj := dog{}
	jsonErr := json.Unmarshal(body, &dog_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Println("Dog Address", dog_obj.Address, "Response Status", dog_obj.Status)
	fmt.Println("DOWNLOADING")
	downloadImages(dog_obj.Address)
}

func getQuote() quote {
	fmt.Printf("Feting a Meaningless Platitude")
	url := "https://quoteslate.vercel.app/api/quotes/random"
	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println(string(body))

	quote_obj := quote{}
	jsonErr := json.Unmarshal(body, &quote_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Println("Quote", quote_obj.Quote)

	return quote_obj
}

func downloadImages(link string) error {

	tmp := getFilePath("/tmp")

	fmt.Println(tmp)
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.Create(filepath.Join(tmp, "daily_dog.jpg"))
	if err != nil {
		return err
	}
	defer file.Close()
	file.ReadFrom(res.Body)
	return nil
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
