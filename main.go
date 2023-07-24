package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"watcher/authenticators"
	"watcher/configs"
	"watcher/connectors"
	"watcher/core"
	"watcher/db"
	"watcher/webapp"
)

func main() {

	params := configs.InitConfigs()

	/* Flags */
	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")
	flag.Parse()

	schedule, _ := time.ParseDuration(*scheduleFlag)

	/*
	   Logging
	*/
	f, err := os.OpenFile(params.PathLogs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	/*






	 */

	go runWatcher(params, schedule)

	app := new(webapp.Application)
	app.Auth.Username = params.WebUser
	app.Auth.Password = params.WebPass

	mux := http.NewServeMux()
	mux.HandleFunc("/status", app.BasicAuth(app.ProtectedHandler))
	// mux.HandleFunc("/", app.UnprotectedHandler)

	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/", fs)

	//
	//
	//
	srv := &http.Server{
		Addr:         params.WebIp,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("starting server on %s", srv.Addr)

	err = srv.ListenAndServeTLS(params.SslPub, params.SslPriv)
	if err != nil {
		log.Printf("%s", err.Error())
	}
}

// Start Program
func runWatcher(params configs.Params, schedule time.Duration) {

	c, err := authenticators.NewClient(params.HostIpa, params.UserIpa, params.UserPassIpa)
	if err != nil {
		log.Fatalf("can not create freeIpa client; err: %s", err.Error())
	}

	conPg, err := db.NewClient()
	if err != nil {
		log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
	}
	defer conPg.CloseDB()

	conSSH, err := connectors.NewClient(params.ActorsUser, params.ActorsPaswd)
	if err != nil {
		log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
	}

	for {

		if params.Mode == "production" {

			actorsList, err := conPg.GetEntity("uds_actortoken")
			if err != nil {
				log.Fatalf("can not get list actors: %s", err.Error())
			}

			usersList, err := c.GetUser(params.GroupIpa)
			if err != nil {
				log.Printf("can not get user list in FreeIPA; err: %s", err.Error())
			}

			userListID, err := c.GetUserID(usersList)
			if err != nil {
				log.Printf("can not get user list ID; err: %s", err.Error())
			}

			/* Удаление папки */
			err = core.DirExpired(params.PathHome, params.DaysRotation, usersList)
			if err != nil {
				log.Printf("can not delete directory; err: %s", err.Error())
			}

			err = core.CreateDirectory(params.PathHome, usersList, userListID)
			if err != nil {
				log.Printf("can not create directory; err: %s", err.Error())
			}

			sshstdout := conSSH.ConnectHost("x2golistsessions_root", actorsList)

			x2gosession, err := connectors.GetSessionX2go(sshstdout)
			if err != nil {
				log.Printf("list session x2go is empty: %s", err.Error())
			}

			//
			//
			//
			//
			//
			//
			//

			core.ShowSession(x2gosession)

			udssession, err := conPg.GetNewRequest()
			if err != nil {
				log.Fatalf("can not; err: %s", err.Error())
			}

			err = core.DiffSession(x2gosession, udssession, conPg, conSSH, actorsList, params.Domain)
			if err != nil {
				log.Fatal("can not:", err.Error())
			}

			// err = core.SetQuota(params.softQuota, params.hardQuota, usersList)
			// if err != nil {
			// 	log.Printf("can not set quota: %s", err.Error())
			// }
		} else {

			log.Printf("APP MODE:%s", params.Mode)

		}

		time.Sleep(schedule)
	}
}
