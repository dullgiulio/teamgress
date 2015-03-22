package libteamgress

import (
	"time"
)

type Filter func(Event) bool

func GetByUser(user string) Filter {
	return func(e Event) bool {
		return e.User.UnixName == user
	}
}

func GetFromTime(time time.Time) Filter {
	unixTime := time.Unix()

	return func(e Event) bool {
		return e.Time.Unix() >= unixTime
	}
}
