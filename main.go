package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"watcher/configs"
	"watcher/watch"
	"watcher/webapp"
)

func main() {

	params := configs.InitConfigs()

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

	go watch.RunWatcher(params)

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
		Addr:         fmt.Sprint("0.0.0.0:", params.WebPort),
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
