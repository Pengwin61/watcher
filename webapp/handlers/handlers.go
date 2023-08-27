package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"watcher/connectors"
	"watcher/core"
	"watcher/db"
)

type Application struct {
	Auth struct {
		Username string
		Password string
	}
}

func (app *Application) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	data := core.Tmp
	tmpl, err := template.ParseFiles("templates/status.html")
	if err != nil {
		log.Printf("%s", err.Error())
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("can`t parse execute template:%s", err.Error())
	}
}

func (app *Application) UnprotectedHandler(w http.ResponseWriter, r *http.Request) {
	data := "print"
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("%s", err.Error())
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("can`t parse execute template:%s", err.Error())
	}

}

func (app *Application) BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.Auth.Username))
			expectedPasswordHash := sha256.Sum256([]byte(app.Auth.Password))

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

func (app *Application) TerminateSession(w http.ResponseWriter, r *http.Request) {

	sessionId := strings.TrimPrefix(r.RequestURI, "/status/terminate/")

	user := strings.SplitAfterN(sessionId, "-", 2)

	if len(user) == 1 {
		fmt.Println("string is nil")
	}

	u := strings.TrimRight(user[0], "-")

	for k, v := range core.Tmp {
		if u != v.Username {
			continue
		} else {

			fmt.Println(k)
			// termSession(v.SessionID, v.Hostname)
			terminationSession(v.SessionID, v.Hostname, v.DbID)
			core.Tmp = remove(core.Tmp, k)
		}

	}

	http.Redirect(w, r, "/status", 301)

}

func remove(slice []core.ViewSession, i int) []core.ViewSession {

	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func terminationSession(sessionId, hostname string, dbId int) {

	con, err := connectors.NewClientSSH("root", "fANu2d$E")
	if err != nil {
		fmt.Println("i can`t create connection to host", err)
	}

	conDb, err := db.NewClient()
	if err != nil {
		fmt.Println("i can`t create connection to database:", err)
	}

	cmdTerminated := "sudo x2goterminate-session " + sessionId

	con.ExecuteCmd(cmdTerminated, hostname)

	err = conDb.UpdateTab(dbId)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("session %s", sessionId)

}
