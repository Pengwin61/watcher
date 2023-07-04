package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {

	cfg, err := ini.Load("settings.cfg")
	if err != nil {
		fmt.Printf("fail to read file: %v", err)
		os.Exit(1)
	}
	/* Flags */
	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")
	softQuotaFlag := flag.String("soft", "1G", "Soft quota")
	hardQuotaFlag := flag.String("hard", "1G", "Hard quota")

	flag.Parse()

	/* settings.cfg */

	mode := cfg.Section("").Key("app_mode").String()
	domain := cfg.Section("").Key("domain").String()

	pathFlag := cfg.Section("paths").Key("home_dir").String()
	daysRotation := cfg.Section("paths").Key("home_dir_days_rotation").String()
	logsPath := cfg.Section("paths").Key("logs").String()

	hostIpa := cfg.Section("FreeIpa").Key("host").String()
	userIpa := cfg.Section("FreeIpa").Key("username").String()
	userpassIpa := cfg.Section("FreeIpa").Key("password").String()
	groudIpa := cfg.Section("FreeIpa").Key("user_group").String()

	actorsUser := cfg.Section("servers").Key("username").String()
	actorsPaswd := cfg.Section("servers").Key("password").String()

	schedule, _ := time.ParseDuration(*scheduleFlag)
	basePath := core.CreatePath(pathFlag)

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
	runWatcher(mode, domain, basePath, daysRotation, *softQuotaFlag, *hardQuotaFlag,
		hostIpa, userIpa, userpassIpa, groudIpa, actorsUser,
		actorsPaswd, schedule)
}

// Start Programm
func runWatcher(appMode, domain, basePath, daysRotation, softQuotaFlag, hardQuotaFlag,
	hostIpa, userIpa, userpassIpa, groudIpa,
	actorsUser, actorsPaswd string, schedule time.Duration) []string {

	c, err := authenticators.NewClient(hostIpa, userIpa, userpassIpa)
	if err != nil {
		log.Fatalf("can not create freeIpa client; err: %s", err.Error())
	}

	conPg, err := db.NewClient()
	if err != nil {
		log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
	}
	defer conPg.CloseDB()

	conSSH, err := connectors.NewClient(actorsUser, actorsPaswd)
	if err != nil {
		log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
	}

	for {

		if appMode == "production" {
			// log.Println("APP MODE:", appMode)

			actorsList, err := conPg.GetEntity("uds_actortoken")
			if err != nil {
				log.Fatalf("can not get list actors: %s", err.Error())
			}

			//
			//

			usersList, err := c.GetUser(groudIpa)
			if err != nil {
				log.Printf("can not get user list in FreeIPA; err: %s", err.Error())
			}

			userListID, err := c.GetUserID(usersList)
			if err != nil {
				log.Printf("can not get user list ID; err: %s", err.Error())
			}

			/* Удаление папки */
			err = core.DirExpired(basePath, daysRotation, usersList)
			if err != nil {
				log.Printf("can not delete directory; err: %s", err.Error())
			}

			err = core.CreateDirectory(basePath, usersList, userListID)
			if err != nil {
				log.Printf("can not create directory; err: %s", err.Error())
			}

			sshstdout := conSSH.ConnectHost("x2golistsessions_root", actorsList)

			x2gosession, err := connectors.GetSessionX2go(sshstdout)
			if err != nil {
				// log.Fatalf("list session x2go is empty: %s", err.Error())
				log.Printf("list session x2go is empty: %s", err.Error())
			}

			udssession, err := conPg.GetNewRequest()
			if err != nil {
				log.Fatalf("can not; err: %s", err.Error())
			}

			err = core.DiffSession(x2gosession, udssession, conPg, conSSH, actorsList, domain)
			if err != nil {
				log.Fatal("can not:", err.Error())
			}

			err = core.SetQuota(softQuotaFlag, hardQuotaFlag, usersList)
			if err != nil {
				log.Println("can not: set quota: %s", err.Error())
			}

		} else {

			log.Println("APP MODE:", appMode)

		}

		time.Sleep(schedule)
	}
}
