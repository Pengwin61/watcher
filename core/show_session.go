package core

import "watcher/connectors"

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

		vTmp := ViewSession{
			Username:     k,
			Status:       v.SessionState,
			Hostname:     v.Hostname,
			StartSession: v.StartDateSession,
			StopSession:  v.StopDateSession}
		Tmp = append(Tmp, vTmp)
	}
}
