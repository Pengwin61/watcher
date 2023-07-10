package core

import (
	"fmt"
	"log"
	"strings"
	"time"
	"watcher/connectors"
	"watcher/db"
)

var dur, _ = time.ParseDuration("4h")
var isExpired bool

func DiffSession(x2gosession map[string]*connectors.User,
	udssession map[string]db.UserService,
	conPg *db.ClientPg, conSsh *connectors.Client,
	actorsList map[string]string, domain string) error {

	var err error

	for session, v := range x2gosession {

		if v.SessionState == "S" {

			if checkExpiration(v.StopDateSession) {

				if val, ok := udssession[session]; ok {

					if checkHostMatches(v.Hostname, val.DepSvcName, domain) {
						hostname := strings.TrimRight(v.Hostname, domain)

						if host, ok := actorsList[hostname]; ok {
							conSsh.TerminateSession(v.SessionPid, host, "x2goterminate-session", conSsh)

							err := conPg.UpdateTab(val.User_service_id)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		} else {
			log.Printf("X2GO RUN SESSION: | %s | %s | %s | %s | %s |\n",
				v.UserSession, v.SessionState, v.Hostname, v.StartDateSession, v.StopDateSession)

		}

		/* check diff sessions */
		diff := difference(x2gosession, udssession)
		log.Println("session removed from database:", diff)

		for _, k := range diff {
			if val, ok := udssession[k]; ok {
				err := conPg.UpdateTab(val.User_service_id)

				if err != nil {
					return err
				}
			}
			if val, ok := x2gosession[k]; ok {

				hostname := strings.TrimRight(val.Hostname, domain)
				hostname = strings.TrimRight(hostname, fmt.Sprint(".", domain))

				if host, ok := actorsList[hostname]; ok {
					conSsh.TerminateSession(val.SessionPid, host, "x2goterminate-session", conSsh)
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

func checkExpiration(t string) bool {

	stopTimeSession := convertTime(t)
	delta := time.Since(stopTimeSession)
	delta = delta.Truncate(time.Second)

	if delta >= dur {
		if delta <= 0 {
			log.Fatal("session sub zero =)")
		}
		isExpired = true
		return isExpired
	}
	return false
}

func checkHostMatches(hostname, depSvcName, domain string) bool {

	hostname = strings.TrimRight(hostname, fmt.Sprint(".", domain))
	depSvcName = strings.TrimLeft(depSvcName, "s-")

	if strings.EqualFold(depSvcName, hostname) {
		return true
	} else {
		fmt.Println("not found:", depSvcName, hostname, "не равны")
	}
	return false
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
