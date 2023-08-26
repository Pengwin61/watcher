package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"watcher/core"
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

func (app *Application) TestH(w http.ResponseWriter, r *http.Request) {
	fmt.Println("123")
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
			// core.Tmp = remove(core.Tmp, k)
			initSession(v.SessionID, v.Hostname)
		}

	}

	http.Redirect(w, r, "/status", 301)

}

// func remove(slice []core.ViewSession, i int) []core.ViewSession {

// 	copy(slice[i:], slice[i+1:])
// 	return slice[:len(slice)-1]
// }

func initSession(sessionId, hostname string) {

	fmt.Println(sessionId, hostname)

}

// func (c *connectors.Client) TermSession() {
// 	c.TerminateSession(sessionId, hostname)
// }

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
