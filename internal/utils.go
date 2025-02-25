package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/rogpeppe/go-internal/lockedfile"
)

func GetFilePath(tmpPath string) string {
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path := fmt.Sprintf("%s%s", curDir, tmpPath)
	return path
}

func FetchSubscribers(path string, active bool) Subscribers {

	jsonFile, err := lockedfile.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	slog.Info("Acquired READ Lock")
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
