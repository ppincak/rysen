package core

import (
	"time"
)

type IntervalType int

const (
	SECOND IntervalType = iota
	MINUTE
	HOUR
	DAY
	WEEK
)

type Limit struct {
	IntervalType IntervalType
	Limit        int64
}

func ToIntervalType(interval string) IntervalType {
	switch interval {
	case "SECOND", "s":
		return SECOND
	case "MINUTE", "m":
		return MINUTE
	case "HOUR", "h":
		return HOUR
	case "DAY", "d":
		return DAY
	case "WEEK", "w":
		return WEEK
	}
	return -1
}

// Converts miliseconds to Duration
func ToDuration(millis int64) time.Duration {
	return time.Duration(millis * int64(time.Millisecond))
}
