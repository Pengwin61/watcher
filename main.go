package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
	"watcher/authenticators"
	"watcher/connectors"
	"watcher/core"
	"watcher/db"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

}

type Params struct {
	mode, domain, basePath, daysRotation, hostIpa, userIpa,
	userPassIpa, groupIpa, actorsUser, actorsPaswd,
	softQuota, hardQuota string
}

type ViewSession struct {
	Username     string
	Status       string
	Hostname     string
	StartSession string
	StopSession  string
}

var Tmp = make([]ViewSession, 0)

func main() {

	cfg, err := ini.Load("settings.cfg")
	if err != nil {
		fmt.Printf("fail to read file: %v", err)
		os.Exit(1)
	}
	/* Flags */
	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")

	flag.Parse()

	/* settings.cfg */

	mode := cfg.Section("").Key("app_mode").String()
	domain := cfg.Section("").Key("domain").String()

	pathFlag := cfg.Section("paths").Key("home_dir").String()
	daysRotation := cfg.Section("paths").Key("home_dir_days_rotation").String()
	logsPath := cfg.Section("paths").Key("logs").String()

	hostIpa := cfg.Section("FreeIpa").Key("host").String()
	userIpa := cfg.Section("FreeIpa").Key("username").String()
	userPassIpa := cfg.Section("FreeIpa").Key("password").String()
	groupIpa := cfg.Section("FreeIpa").Key("user_group").String()

	actorsUser := cfg.Section("servers").Key("username").String()
	actorsPaswd := cfg.Section("servers").Key("password").String()

	softQuota := cfg.Section("UserQuota").Key("softQuota").String()
	hardQuota := cfg.Section("UserQuota").Key("hardQuota").String()

	schedule, _ := time.ParseDuration(*scheduleFlag)
	basePath := core.CreatePath(pathFlag)

	/*

	 */
	var params = Params{mode: mode, domain: domain, basePath: basePath,
		daysRotation: daysRotation, hostIpa: hostIpa, userIpa: userIpa,
		userPassIpa: userPassIpa, groupIpa: groupIpa, actorsUser: actorsUser,
		actorsPaswd: actorsPaswd, softQuota: softQuota, hardQuota: hardQuota}
	//

	/*
	   Logging
	*/
	f, err := os.OpenFile(logsPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	/*






	 */

	go runWatcher(params, schedule)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := Tmp
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, data)
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

// Start Program
func runWatcher(params Params, schedule time.Duration) {

	c, err := authenticators.NewClient(params.hostIpa, params.userIpa, params.userPassIpa)
	if err != nil {
		log.Fatalf("can not create freeIpa client; err: %s", err.Error())
	}

	conPg, err := db.NewClient()
	if err != nil {
		log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
	}
	defer conPg.CloseDB()

	conSSH, err := connectors.NewClient(params.actorsUser, params.actorsPaswd)
	if err != nil {
		log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
	}

	for {

		if params.mode == "production" {

			actorsList, err := conPg.GetEntity("uds_actortoken")
			if err != nil {
				log.Fatalf("can not get list actors: %s", err.Error())
			}

			usersList, err := c.GetUser(params.groupIpa)
			if err != nil {
				log.Printf("can not get user list in FreeIPA; err: %s", err.Error())
			}

			userListID, err := c.GetUserID(usersList)
			if err != nil {
				log.Printf("can not get user list ID; err: %s", err.Error())
			}

			/* Удаление папки */
			err = core.DirExpired(params.basePath, params.daysRotation, usersList)
			if err != nil {
				log.Printf("can not delete directory; err: %s", err.Error())
			}

			err = core.CreateDirectory(params.basePath, usersList, userListID)
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
			Tmp = nil
			for k, v := range x2gosession {

				vTmp := ViewSession{
					Username: k, Status: v.SessionState, Hostname: v.Hostname, StartSession: v.StartDateSession, StopSession: v.StopDateSession}
				Tmp = append(Tmp, vTmp)
			}

			udssession, err := conPg.GetNewRequest()
			if err != nil {
				log.Fatalf("can not; err: %s", err.Error())
			}

			err = core.DiffSession(x2gosession, udssession, conPg, conSSH, actorsList, params.domain)
			if err != nil {
				log.Fatal("can not:", err.Error())
			}

			// err = core.SetQuota(params.softQuota, params.hardQuota, usersList)
			// if err != nil {
			// 	log.Printf("can not set quota: %s", err.Error())
			// }
		} else {

			log.Println("APP MODE:", params.mode)

		}

		time.Sleep(schedule)
	}
}
