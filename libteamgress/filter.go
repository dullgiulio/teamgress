package libteamgress

import (
	"time"
)

type Filter func(Event) bool

func GetByUser(key string) Filter {
	return func(e Event) bool {
		return e.User.Key == key
	}
}

func GetFromTime(time time.Time) Filter {
	unixTime := time.Unix()

	return func(e Event) bool {
		return e.Time.Unix() >= unixTime
	}
}
