package core

import (
	"fmt"
	"log"
	"strings"
	"time"
	"watcher/connectors"
	"watcher/db"
)

func convertTime(t string) time.Time {

	layout := "2006-01-02T15:04:05"
	timeSession, err := time.Parse(layout, t)
	if err != nil {
		log.Println(err)
	}
	return timeSession
}

func viewTimeFormat(t string) string {

	time := convertTime(t)

	strDate := time.Format("02-01-2006")
	strTime := time.Format("15:04:05")

	res := fmt.Sprintln(strDate, "\n", strTime)

	return res
}

func viewStatusFormat(s *string) *string {

	switch *s {
	case "S":
		*s += "toped"

	case "R":
		*s += "unning"
	}
	return s
}

func viewHostname(s string) string {
	res := strings.Split(s, ".")

	return res[0]
}

func checkExpirationSession(t, state string,
	durationSession time.Duration) (bool, time.Duration) {

	var msk, _ = time.ParseDuration("3h")

	stopTimeSession := convertTime(t)
	delta := time.Since(stopTimeSession)
	delta = delta.Truncate(time.Second)

	delta += msk

	if delta >= durationSession && state != "R" {

		return true, delta
	}

	return false, delta
}

func containsIpaUser(array map[string]*connectors.User, value string) bool {
	for k := range array {
		if k == value {
			return true
		}
	}
	return false
}
func containsDbUser(array map[string]db.UserService, value string) bool {
	for k := range array {
		if k == value {
			return true
		}
	}
	return false
}

func difference(x2gosession map[string]*connectors.User, udssession map[string]db.UserService) (diff []string) {

	diffArray := []string{}

	for k := range x2gosession {
		if !containsDbUser(udssession, k) {
			diffArray = append(diffArray, k)
		}
	}

	for k := range udssession {
		if !containsIpaUser(x2gosession, k) {
			diffArray = append(diffArray, k)
		}
	}

	return diffArray
}

var isView = true

func printSesessionHeader() {

	if isView {

		fmt.Printf(" \t\t | %20s | %5s | %23s | %20s | %20s | %10s | %5s |\n",
			"--------------------", "-----", "-----------------------", "--------------------", "--------------------", "----------", "-----")

		fmt.Printf("X2GO RUN SESSION | %20s | %5s | %23s | %20s | %20s | %10s | %5s |\n",
			"Name", "State", "Hostname", "Start Time", "Stop Time", "Overtime", "DB ID")

		fmt.Printf(" \t\t | %20s | %5s | %23s | %20s | %20s | %10s | %5s |\n",
			"--------------------", "-----", "-----------------------", "--------------------", "--------------------", "----------", "-----")

	}
	isView = false
}

func printSessionHeaderEnd() {

	fmt.Printf(" \t\t | %20s | %5s | %23s | %20s | %20s | %10s | %5s |\n",
		"--------------------", "-----", "-----------------------", "--------------------", "--------------------", "----------", "-----")

}

func printSession(userSession, sessionState, hostname,
	startDateSession, stopDateSession string, overtime time.Duration, dbId int) {

	fmt.Printf(" \t\t | %20s | %5s | %23s | %20s | %20s | %10s | %5d |\n",
		userSession, sessionState, hostname, startDateSession, stopDateSession, overtime, dbId)
}
