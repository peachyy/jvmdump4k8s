package util

import "time"

func FormartdateNow() string {
	return time.Now().Format("20060102150405")
}
