package main

import (
	"crypto/sha256"
	"crypto/subtle"
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

	// webPort := cfg.Section("web").Key("port").String()
	sslpub := cfg.Section("web").Key("ssl_public").String()
	sslpriv := cfg.Section("web").Key("ssl_private").String()

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

	// go runWatcher(params, schedule)
	app := new(application)
	app.auth.username = "admin"
	app.auth.password = "admin"

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/status", app.basicAuth(app.protectedHandler))
	mux.HandleFunc("/", app.unprotectedHandler)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("hello")
	// })

	srv := &http.Server{
		Addr:         ":8181",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
	// serve index (and anything else) as https

	err = srv.ListenAndServeTLS(sslpub, sslpriv)
	// err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", webPort), nil)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	runWatcher(params, schedule)
}
func redirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req,
		"https://"+req.Host+req.URL.String(),
		http.StatusMovedPermanently)
}

func index(w http.ResponseWriter, req *http.Request) {
	// all calls to unknown url paths should return 404
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	http.ServeFile(w, req, "index.html")
}

func (app *application) protectedHandler(w http.ResponseWriter, r *http.Request) {
	data := core.Tmp
	tmpl, err := template.ParseFiles("templates/status.html")
	if err != nil {
		log.Printf("%s", err.Error())
	}
	tmpl.Execute(w, data)
}

func (app *application) unprotectedHandler(w http.ResponseWriter, r *http.Request) {
	data := "print"
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("%s", err.Error())
	}
	tmpl.Execute(w, data)
}

func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

type application struct {
	auth struct {
		username string
		password string
	}
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

			core.ShowSession(x2gosession)

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
