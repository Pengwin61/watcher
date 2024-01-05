package core

import (
	"log"
	"time"
	"watcher/internal/connections"
	"watcher/internal/utils"
)

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
	SessionID    string
	DbID         int
}

type viewServer struct {
	Hostname string
	State    string
	Ip       string
	Uptime   time.Time
}

var viewUsers = make([]ViewSession, 0)
var viewServers = make([]viewServer, 0)

func Show(person *[]PersonSession) {
	ShowSession(person)
	ShowServers()
}

func GetUsersView() []ViewSession {
	return viewUsers
}
func SetUserView(users []ViewSession) {
	viewUsers = users
}
func GetServerView() []viewServer {
	return viewServers
}

func GetPersonalView(user string) []ViewSession {
	tempusers := make([]ViewSession, 0)
	for _, v := range viewUsers {
		if v.Username == user {
			temp := &ViewSession{
				Username:     v.Username,
				Status:       v.Status,
				Hostname:     v.Hostname,
				StartSession: v.StartSession,
				StopSession:  v.StopSession,
				SessionID:    v.SessionID}

			tempusers = append(tempusers, *temp)
		}
	}
	return tempusers
}

func ShowSession(personsSession *[]PersonSession) {
	viewUsers = nil

	if personsSession != nil {

		for _, v := range *personsSession {

			vTmp := ViewSession{
				Username:     v.Username,
				Status:       *viewStatusFormat(&v.State),
				Hostname:     viewHostname(v.Hostname),
				StartSession: utils.ViewTimeFormat(v.StartDateSession),
				StopSession:  utils.ViewTimeFormat(v.StopDateSession),
				SessionID:    v.SessionID,
				DbID:         v.DbID}

			viewUsers = append(viewUsers, vTmp)
		}
	} else {
		viewUsers = nil
	}

}

func ShowServers() {
	viewServers = nil

	actorsList, err := connections.Conn.Database.GetEntity("uds_actortoken")
	if err != nil {
		log.Fatalf("can not get list actors: %s", err.Error())
	}

	if actorsList == nil {
		return
	}
	for k, v := range actorsList {

		vTmp := viewServer{
			Hostname: k,
			State:    "Up",
			Ip:       v,
			Uptime:   time.Now().Truncate(time.Minute)}
		viewServers = append(viewServers, vTmp)

	}
}
