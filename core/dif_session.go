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

		expired, delta := checkExpirationSession(v.StopDateSession, v.SessionState, v.UserSession, expirationSession)

		if expired {

			if val, ok := udssession[session]; ok {
				hostEqual, hostname := checkHostMatches(v.Hostname, val.DepSvcName, domain)

				if hostEqual {
					if host, ok := actorsList[hostname]; ok {
						conSsh.TerminateSession(v.SessionPid, host, "sudo x2goterminate-session")

						err := conPg.UpdateTab(val.UserServiceId)
						if err != nil {
							return err
						}
						log.Printf("session %s expired, overtime:%s update database ID:%d", v.UserSession, delta-expirationSession, val.UserServiceId)
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
				err := conPg.UpdateTab(val.UserServiceId)
				log.Printf("session %s removed from database ID:%d, watcher didn't find session record in x2go, state:%s", val.Username, val.UserServiceId, val.State)

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

func convertTime(t string) time.Time {

	layout := "2006-01-02T15:04:05"
	timeSession, err := time.Parse(layout, t)
	if err != nil {
		log.Println(err)
	}
	return timeSession
}

func checkExpirationSession(t, state, user string, durationSession time.Duration) (bool, time.Duration) {

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
