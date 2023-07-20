package webapp

import (
	"crypto/sha256"
	"crypto/subtle"
	"html/template"
	"log"
	"net/http"
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
	tmpl.Execute(w, data)
}

func (app *Application) UnprotectedHandler(w http.ResponseWriter, r *http.Request) {
	data := "print"
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("%s", err.Error())
	}
	tmpl.Execute(w, data)
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
