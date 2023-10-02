package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"html/template"
	"log"
	"net/http"
	"strings"
	"watcher/connections"
	"watcher/core"
)

type Application struct {
	Auth struct {
		Username string
		Password string
	}
}

func (app *Application) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	data := core.GetUsersView()
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
		log.Println("string is nil")
	}

	u := strings.TrimRight(user[0], "-")

	for k, v := range core.GetUsersView() {
		if u != v.Username {
			continue
		} else {
			connections.Conn.SSH.TerminateSession(v.SessionID, v.Hostname)
			log.Printf("the session %s was terminated by the administrator", sessionId)

			err := connections.Conn.Database.UpdateTab(v.DbID)
			if err != nil {
				log.Println(err)
			}
			core.SetUserView(core.RemoveSlice(core.GetUsersView(), k))
		}

	}

	http.Redirect(w, r, "/status", 301)

}
