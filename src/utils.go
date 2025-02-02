package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"time"

	"github.com/rogpeppe/go-internal/lockedfile"
)

func downloadImages(link string, fileName string, targetFolder string) (dogFilePath string, err error) {
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

func fetchSubscribers(path string, active bool) Subscribers {

	jsonFile, err := lockedfile.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	fmt.Println("READ LOCK")
	var subscribers Subscribers
	var outSubscribers Subscribers

	json.Unmarshal(byteValue, &subscribers)
	if active {
		for _, subscriber := range subscribers {
			if subscriber.DateUnsubscribed.IsZero() {
				outSubscribers = append(outSubscribers, subscriber)
			}
		}
		return outSubscribers
	}
	return subscribers
}

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
