package subscription

import (
	"time"
)

type Subscribers []struct {
	Email            string    `json:"email"`
	DateSubscribed   time.Time `json:"dateSubscribed"`
	DateUnsubscribed time.Time `json:"dateUnsubscribed"`
}

type subscriber struct {
	Email            string    `json:"email"`
	DateSubscribed   time.Time `json:"dateSubscribed"`
	DateUnsubscribed time.Time `json:"dateUnsubscribed"`
}

type email struct {
	Email string `json:"email"`
}
