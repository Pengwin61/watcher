package core

import (
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

func ManageSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService,
	conPg *db.ClientPg, conSsh *connectors.Client,
	expirationSession time.Duration) {

	cleanupSession(x2gosession, udssession, conPg, conSsh)
	personsSession := mergeSession(x2gosession, udssession)
	expirationOvertime(&personsSession, expirationSession, conPg, conSsh)

	ShowSession(&personsSession)

}

func cleanupSession(x2gosession map[string]*connectors.User, udssession map[string]db.UserService,
	conPg *db.ClientPg, conSsh *connectors.Client) error {

	if len(x2gosession) != len(udssession) {

		diff := difference(x2gosession, udssession)

		for _, k := range diff {

			if val, ok := udssession[k]; ok {

				err := conPg.UpdateTab(val.DbID)
				if err != nil {
					return err
				}
				log.Printf("session %s removed from database ID:%d, watcher didn't find session record in x2go", val.Username, val.DbID)
			} else if val, ok := x2gosession[k]; ok {
				conSsh.TerminateSession(val.SessionPid, val.Hostname, "sudo x2goterminate-session")
				log.Printf("session %s terminated, user %s logged in incorrectly.", val.SessionPid, val.UserSession)
			}

		}
	}
	return nil
}

func mergeSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService) []PersonSession {

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

func expirationOvertime(personsSession *[]PersonSession, expSesson time.Duration,
	conPg *db.ClientPg, conSsh *connectors.Client) error {

	for _, session := range *personsSession {
		expired, delta := checkExpirationSession(session.StopDateSession, session.SessionState, expSesson)

		if expired {

			conSsh.TerminateSession(session.SessionPid, session.Hostname, "sudo x2goterminate-session")
			err := conPg.UpdateTab(session.DbID)
			if err != nil {
				return err
			}

			log.Printf("session %s expired, overtime:%s update database ID:%d", session.UserSession, delta-expSesson, session.DbID)
		}
		if !expired && session.SessionState != "S" {
			log.Printf("X2GO RUN SESSION: | %20s | %s | %s | %s | %s | %5s | \n",
				session.UserSession, session.SessionState, session.Hostname,
				session.StartDateSession, session.StopDateSession, delta-expSesson)
		}

	}
	return nil
}
