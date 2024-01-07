package utils

import (
	"fmt"
	"log"
	"time"
)

func ConvertTime(t string) time.Time {

	layout := "2006-01-02T15:04:05"
	timeSession, err := time.Parse(layout, t)
	if err != nil {
		log.Println(err)
	}
	return timeSession
}

func ViewTimeFormat(t string) string {

	time := ConvertTime(t)

	strDate := time.Format("02-01-2006")
	strTime := time.Format("15:04:05")

	res := fmt.Sprintln(strDate, "\n", strTime)

	return res
}
