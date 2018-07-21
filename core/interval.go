package core

import (
	"time"
)

type Limit struct {
	Interval time.Duration
	Limit    int64
}

// Converts string interval expression to Duration
func ToDuration(interval string) time.Duration {
	switch interval {
	case "SECOND", "s":
		return time.Duration(time.Second)
	case "MINUTE", "m":
		return time.Duration(time.Minute)
	case "HOUR", "h":
		return time.Duration(time.Hour)
	case "DAY", "d":
		return time.Duration(time.Hour * 24)
	case "WEEK", "w":
		return time.Duration(time.Hour * 24 * 7)
	}
	return -1
}

// Converts miliseconds to Duration
func ToDurationMillis(millis int64) time.Duration {
	return time.Duration(millis * int64(time.Millisecond))
}
