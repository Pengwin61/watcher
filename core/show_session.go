package core

import (
	"fmt"
	"watcher/connectors"
)

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
}

var Tmp = make([]ViewSession, 0)

func ShowSession(x2gosession map[string]*connectors.User) {
	Tmp = nil

	for k, v := range x2gosession {

		v.StartDateSession = viewTimeFormat(v.StartDateSession)
		v.StopDateSession = viewTimeFormat(v.StopDateSession)

		switch v.SessionState {
		case "S":
			v.SessionState += "toped"

		case "R":
			v.SessionState += "unning"
		}

		vTmp := ViewSession{
			Username:     k,
			Status:       v.SessionState,
			Hostname:     v.Hostname,
			StartSession: v.StartDateSession,
			StopSession:  v.StopDateSession}
		Tmp = append(Tmp, vTmp)
	}
}

func viewTimeFormat(t string) string {

	time := convertTime(t)

	strDate := time.Format("02-01-2006")
	strTime := time.Format("15:04:05")

	res := fmt.Sprintln(strDate, "\n", strTime)

	return res
}
