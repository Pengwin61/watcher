package core

import (
	"fmt"
	"log"
	"strings"
	"time"
	"watcher/connectors"
	"watcher/db"
)

func DiffSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService,
	conPg *db.ClientPg, conSsh *connectors.Client,
	actorsList map[string]string, domain string, expirationSession time.Duration) error {

	var err error

	for session, v := range x2gosession {

		expired, delta := checkExpirationSession(v.StopDateSession, v.SessionState, expirationSession)

		if expired {

			if val, ok := udssession[session]; ok {
				hostEqual, hostname := checkHostMatches(v.Hostname, val.DepSvcName, domain)

				if hostEqual {
					if host, ok := actorsList[hostname]; ok {
						conSsh.TerminateSession(v.SessionPid, host, "sudo x2goterminate-session")

						err := conPg.UpdateTab(val.DbID)
						if err != nil {
							return err
						}
						log.Printf("session %s expired, overtime:%s update database ID:%d", v.UserSession, delta-expirationSession, val.DbID)
					}
				}
			}
		}
		if !expired && v.SessionState != "S" {
			log.Printf("X2GO RUN SESSION: | %20s | %s | %s | %s | %s | %t\n",
				v.UserSession, v.SessionState, v.Hostname, v.StartDateSession, v.StopDateSession, expired)
		}

		/* check diff sessions */
		diff := difference(x2gosession, udssession)

		// deletes the session when the user presses the logoff button
		for _, k := range diff {
			if val, ok := udssession[k]; ok {
				err := conPg.UpdateTab(val.DbID)
				log.Printf("session %s removed from database ID:%d, watcher didn't find session record in x2go", val.Username, val.DbID)

				if err != nil {
					return err
				}
			}
			// deletes session when user connected bypassing openuds
			if val, ok := x2gosession[k]; ok {

				hostname := strings.TrimRight(val.Hostname, domain)
				hostname = strings.TrimRight(hostname, fmt.Sprint(".", domain))

				if host, ok := actorsList[hostname]; ok {
					conSsh.TerminateSession(val.SessionPid, host, "sudo x2goterminate-session")
					log.Printf("session %s terminated, user %s logged in incorrectly.", val.SessionPid, val.UserSession)
				}
			}
		}
	}
	return err
}

func checkHostMatches(hostname, depSvcName, domain string) (bool, string) {

	hostname = strings.TrimRight(hostname, fmt.Sprint(".", domain))
	depSvcName = strings.TrimLeft(depSvcName, "s-")

	if strings.EqualFold(depSvcName, hostname) {
		return true, hostname
	} else {
		fmt.Println("not found:", depSvcName, hostname, "are not equal")
	}
	return false, hostname
}
