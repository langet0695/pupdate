package auth

import (
	"time"
)

type Subscribers []struct {
	Email            string    `json:"email"`
	DateSubscribed   time.Time `json:"dateSubscribed"`
	DateUnsubscribed time.Time `json:"dateUnsubscribed"`
}
