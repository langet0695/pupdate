package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rogpeppe/go-internal/lockedfile"

	gomail "gopkg.in/mail.v2"
)

func editSubscribers(path string, subscriber subscriber, action string) (outSubscriber subscriber, err error) {
	// Open the file in edit mode which gives us a write lock
	fmt.Println("OPENING IN EDIT MODE WRITE LOCK")
	jsonFile, err := lockedfile.Edit(path)
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var subscribers Subscribers
	json.Unmarshal(byteValue, &subscribers)

	// Execute logic to handle create and deletes
	if action == "delete" {
		for i := 0; i < len(subscribers); i++ {
			if subscriber.Email == subscribers[i].Email && subscribers[i].DateUnsubscribed.IsZero() {
				subscribers[i].DateUnsubscribed = time.Now()
				outSubscriber = subscribers[i]
			}
		}
	}
	if action == "create" {
		var exists bool = false
		for i := 0; i < len(subscribers); i++ {
			if subscriber.Email == subscribers[i].Email && subscribers[i].DateUnsubscribed.IsZero() {
				exists = true
				outSubscriber = subscribers[i]
			}
		}
		if !exists {
			subscribers = append(subscribers, subscriber)
			outSubscriber = subscriber
		}
	}
	// Construct the byte slice to and write it to destination. Make sure to truncate to remove excess
	out, err := json.Marshal(subscribers)
	if err != nil {
		return
	}
	bytes, err := jsonFile.WriteAt(out, 0)
	if err != nil {
		return
	}
	err = jsonFile.Truncate(int64(len(out)))
	if err != nil {
		return
	}
	fmt.Println("BYTES WRITTEN: ", bytes)
	return

}

func subscriptionNotificaitons(email string, action string) {

	MAIL_USER := os.Getenv("MAIL_USER")
	MAIL_PASSWORD := os.Getenv("MAIL_PASSWORD")

	message := gomail.NewMessage()
	message.SetHeader("From", MAIL_USER)
	message.SetHeader("To", email)

	if action == "new_subscription" {
		message.SetHeader("Subject", "Woof Woof!")
		body := fmt.Sprintf(`<p>Hello there,
							<br><br>
							Welcome to the paw-some world of daily dog images! &#127881; Did you know a recent study found that looking at pictures of dogs can boost your mood and reduce stress? (Science doesn’t lie!) So get ready for your inbox to be flooded with the cutest, silliest, and most photogenic pups you've ever seen. Whether you're having a ruff day or just need some tail-wagging joy, we've got you covered.
							<br><br>
							Our dogs have been training their whole lives to make your day better, one woof at a time. &#128054; So, sit back, relax, and enjoy the cuteness overload—no leash required!
							<br><br>
							Warning: Excessive smiling and spontaneous “aww-ing” may occur.
							<br><br>
							Paw-lease enjoy the daily dog magic! &#128062;
							<br><br>
							The Pupdate Team &#128021;&#128140;</p>
							<br>
							<i>To unsubscribe email: <b>pupdate+unsubscribe@gmail.com</b></i>`)
		message.SetBody("text/html", body)
	} else if action == "unsubscribing" {
		message.SetHeader("Subject", "I guess this is goodbye...")
		timeFormat := "2006-01-02 15:04:05.999999999 -0700 PST"
		timeStr := time.Now().Format(timeFormat)
		body := fmt.Sprintf(`<p>As of %s, you have unsubscribed from Pupdate. <br> You must like cats &#128572;</p>`, timeStr)
		message.SetBody("text/html", body)

	} else if action == "existing_subscription" {
		message.SetHeader("Heyyyyy, I know you!!")
		message.SetBody("text", "Silly goose, you've already subscribed!")

	}

	dialer := gomail.NewDialer("smtp.gmail.com", 587, MAIL_USER, MAIL_PASSWORD)

	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
