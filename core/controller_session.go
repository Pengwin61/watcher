package core

import (
	"log"
	"strings"
	"time"
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
	udssession map[string]db.UserService,
	conPg *db.ClientPg, conSsh *connectors.Client,
	timeExpiration time.Duration) {

	cleanupSession(x2gosession, udssession, conPg, conSsh)
	personsSession := mergeSession(x2gosession, udssession)
	expirationOvertime(&personsSession, timeExpiration, conPg, conSsh)

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
				conSsh.TerminateSession(val.SessionID, val.Hostname)
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

func expirationOvertime(personsSession *[]PersonSession, timeExpiration time.Duration,
	conPg *db.ClientPg, conSsh *connectors.Client) error {

	for _, session := range *personsSession {
		expired, delta := checkExpirationSession(session.StopDateSession, session.State, timeExpiration)

		if expired {

			conSsh.TerminateSession(session.SessionID, session.Hostname)
			err := conPg.UpdateTab(session.DbID)
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
