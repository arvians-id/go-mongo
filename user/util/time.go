package util

import (
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

func CurrentTime() time.Time {
	timeNow, err := time.Parse(time.RFC3339Nano, time.Now().Format(time.RFC3339Nano))
	if err != nil {
		return time.Time{}
	}

	return timeNow
}
