package core

import (
	"log"
	"strings"
	"time"
	"watcher/connections"
	"watcher/connectors"
	"watcher/db"
)

type PersonSession struct {
	Username         string
	SessionID        string
	State            string
	StartDateSession string
	StopDateSession  string
	Hostname         string
	DbID             int
	DbState          string
	DbUsername       string
}

func ManageSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService, timeExpiration time.Duration) error {

	err := cleanupSession(x2gosession, udssession)
	if err != nil {
		return err
	}
	personsSession := mergeSession(x2gosession, udssession)
	err = expirationOvertime(&personsSession, timeExpiration)
	if err != nil {
		return err
	}

	ShowSession(&personsSession)

	return err
}

func cleanupSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService) error {

	if len(x2gosession) != len(udssession) {

		diff := difference(x2gosession, udssession)

		for _, k := range diff {

			if val, ok := udssession[k]; ok {

				err := connections.Conn.Database.UpdateTab(val.DbID)
				if err != nil {
					return err
				}
				log.Printf("session %s removed from database ID:%d, watcher didn't find session record in x2go", val.Username, val.DbID)
			} else if val, ok := x2gosession[k]; ok {

				connections.Conn.SSH.TerminateSession(val.SessionID, val.Hostname)
				log.Printf("session %s terminated, user %s logged in incorrectly.", val.SessionID, val.Username)
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
					Username:         xValue.Username,
					SessionID:        xValue.SessionID,
					State:            xValue.State,
					StartDateSession: xValue.InitTime,
					StopDateSession:  xValue.LastTime,
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

func expirationOvertime(personsSession *[]PersonSession,
	timeExpiration time.Duration) error {

	for _, session := range *personsSession {
		expired, delta := checkExpirationSession(session.StopDateSession, session.State, timeExpiration)

		if expired {

			connections.Conn.SSH.TerminateSession(session.SessionID, session.Hostname)
			err := connections.Conn.Database.UpdateTab(session.DbID)
			if err != nil {
				return err
			}

			log.Printf("session %s expired, overtime:%s update database ID:%d", session.Username, delta-timeExpiration, session.DbID)
		}
		if !expired && session.State != "S" {
			untilEnd := timeExpiration - delta

			printSesessionHeader()

			printSession(session.Username, session.State, session.Hostname,
				session.StartDateSession, session.StopDateSession, untilEnd.Truncate(time.Minute), session.DbID)

		}
		printSessionHeaderEnd()
	}
	return nil
}
