package utils

import (
	"os"
	"time"
)

func TimestampNow() string {
	return time.Now().Format(os.Getenv("TIME_FORMAT"))
}
