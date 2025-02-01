package main

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

type dog struct {
	Address string `json:"message"`
	Status  string `json:"status"`
}

type quote struct {
	Id     int      `json:"id"`
	Quote  string   `json:"quote"`
	Author string   `json:"author"`
	Length int      `json:"length"`
	Tags   []string `json:"tags"`
}

type email struct {
	Email string `json:email`
}
