package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/langet/pupdate/internal"
	gomail "gopkg.in/mail.v2"
)

func downloadImages(link string, fileName string, targetFolder string) (dogFilePath string, err error) {
	tmp := internal.GetFilePath(targetFolder)
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
func getPuppy(fileName string) string {
	URL := "https://dog.ceo/api/breeds/image/random"
	TARGET_FOLDER := "/tmp"

	res, getErr := http.Get(URL)

	if getErr != nil {
		slog.Error(getErr.Error())
		panic(getErr.Error())
	}
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		slog.Error(readErr.Error())
		panic(readErr.Error())
	}
	dog_obj := dog{}
	jsonErr := json.Unmarshal(body, &dog_obj)
	if jsonErr != nil {
		slog.Error(jsonErr.Error())
		panic(jsonErr.Error())
	}
	filePath, err := downloadImages(dog_obj.Address, fileName, TARGET_FOLDER)
	if err != nil {
		slog.Error(err.Error())
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
