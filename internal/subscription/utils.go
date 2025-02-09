package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/rogpeppe/go-internal/lockedfile"
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
