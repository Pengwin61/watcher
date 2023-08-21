package main

import (
	"io"
	"log"
	"os"
	"watcher/configs"
	"watcher/watch"
	"watcher/webapp"
)

func main() {

	params := configs.InitConfigs()

	/*
	   Logging
	*/
	f, err := os.OpenFile(params.Paths.Logs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	/*






	 */

	go watch.RunWatcher(params)

	webClient := webapp.NewClient(params.Web.Port)
	webClient.RunWeb(params.Web.User, params.Web.Pass,
		params.Web.SslPub, params.Web.SslPriv)

	// app := new(webapp.Application)
	// app.Auth.Username = params.Web.User
	// app.Auth.Password = params.Web.Pass

	// mux := http.NewServeMux()
	// mux.HandleFunc("/status", app.BasicAuth(app.ProtectedHandler))
	// // mux.HandleFunc("/", app.UnprotectedHandler)

	// fs := http.FileServer(http.Dir("templates"))
	// mux.Handle("/", fs)

	//
	//
	//
	// srv := &http.Server{
	// 	Addr:         fmt.Sprint("0.0.0.0:", params.Web.Port),
	// 	Handler:      mux,
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// }
	// log.Printf("starting server on %s", srv.Addr)

	// err = srv.ListenAndServeTLS(params.SslPub, params.SslPriv)
	// if err != nil {
	// 	log.Printf("%s", err.Error())
	// }
}
