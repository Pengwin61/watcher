package core

import (
	"fmt"
	"log"
	"strings"
	"time"
	"watcher/connectors"
	"watcher/db"
)

type PersonSession struct {
	UserSession      string
	SessionPid       string
	SessionState     string
	StartDateSession string
	StopDateSession  string
	Hostname         string
	DbID             int
	DbState          string
	DbUsername       string
}

func MergeSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService) []PersonSession {

	if len(x2gosession) != len(udssession) {

		diff := difference(x2gosession, udssession)

		for _, k := range diff {

			if val, ok := udssession[k]; ok {
				log.Printf("session removed from database id: %d", val.DbID)
			} else if val, ok := x2gosession[k]; ok {
				log.Printf("session %s terminated, user %s logged in incorrectly.", val.SessionPid, val.UserSession)
			}

		}
	}

	var PersonsSession = make([]PersonSession, 0)

	for xSession, xValue := range x2gosession {

		if val, ok := udssession[xSession]; ok {

			if strings.ContainsAny(xValue.Hostname, val.DepSvcName) {

				vTmp := &PersonSession{
					UserSession:      xValue.UserSession,
					SessionPid:       xValue.SessionPid,
					SessionState:     xValue.SessionState,
					StartDateSession: xValue.StartDateSession,
					StopDateSession:  xValue.StopDateSession,
					Hostname:         xValue.Hostname,
					DbID:             val.DbID,
					DbState:          val.State,
					DbUsername:       val.Username,
				}
				PersonsSession = append(PersonsSession, *vTmp)
			}
		}
	}
	return PersonsSession
}

func NewDiffer(personsSession []PersonSession, expSesson time.Duration) {
	for _, session := range personsSession {
		expired, delta := checkExpirationSession(session.StopDateSession, session.SessionState, expSesson)

		if expired {
			fmt.Printf("session terminate on host: %s, delta:%s", session.Hostname, delta)

			log.Printf("session %s expired, overtime:%s update database ID:%d", session.UserSession, delta-expSesson, session.DbID)
		}
		if !expired && session.SessionState != "S" {
			log.Printf("X2GO RUN SESSION: | %20s | %s | %s | %s | %s | %t\n",
				session.UserSession, session.SessionState, session.Hostname,
				session.StartDateSession, session.StopDateSession, expired)
		}

	}
}
