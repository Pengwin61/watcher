package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"watcher/configs"
	"watcher/logs"
	"watcher/watch"
	"watcher/webapp"
)

func main() {

	params := configs.InitConfigs()

	logfile := logs.InitLogs(params.Paths.Logs)
	defer logfile.CloseFile()

	/*


	 */

	go watch.RunWatcher(params)

	app := new(webapp.Application)
	app.Auth.Username = params.Web.User
	app.Auth.Password = params.Web.Pass

	mux := http.NewServeMux()
	mux.HandleFunc("/status", app.BasicAuth(app.ProtectedHandler))
	// mux.HandleFunc("/", app.UnprotectedHandler)

	fs := http.FileServer(http.Dir("templates"))
	mux.Handle("/", fs)

	//
	//
	//
	srv := &http.Server{
		Addr:         fmt.Sprint("0.0.0.0:", params.Web.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("starting server on %s", srv.Addr)

	err := srv.ListenAndServeTLS(params.SslPub, params.SslPriv)
	if err != nil {
		log.Printf("%s", err.Error())
	}
}
