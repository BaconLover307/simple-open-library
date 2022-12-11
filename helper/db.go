package helper

import (
	"fmt"
	"time"
)

func ParseDBTime(dateTime string) string {
	schedule, err := time.Parse(time.RFC3339, dateTime)
	PanicIfError(err)
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
        schedule.Year(), schedule.Month(), schedule.Day(),
        schedule.Hour(), schedule.Minute(), schedule.Second())
}