package main

import "time"

var subscribers = []subscriber{
	{Email: "tlange1124@gmail.com", DateSubscribed: time.Now(), DateUnsubscribed: time.Time{}},
	{Email: "domino.birze@gmail.com", DateSubscribed: time.Now(), DateUnsubscribed: time.Time{}},
}

var subscriber_map = map[string]subscriber{
	"tlange1124@gmail.com":   {DateSubscribed: time.Now(), DateUnsubscribed: time.Time{}},
	"domino.birze@gmail.com": {DateSubscribed: time.Now(), DateUnsubscribed: time.Time{}},
}
