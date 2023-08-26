package webapp

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"watcher/webapp/handlers"
)

type Client struct {
	con *http.Server
	mux *http.ServeMux
}

func NewClient(webPort string) *Client {
	mux, srv := initWeb(webPort)

	return &Client{con: srv, mux: mux}
}

func initWeb(webPort string) (*http.ServeMux, *http.Server) {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:         fmt.Sprint("0.0.0.0:", webPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("starting server on %s", srv.Addr)

	return mux, srv
}

func (c *Client) RunWeb(user, pass, sslpub, sslpriv string) {

	c.runHandlers(user, pass)

	err := c.con.ListenAndServeTLS(sslpub, sslpriv)
	if err != nil {
		log.Printf("%s", err.Error())
	}
}

func (c *Client) runHandlers(webUser, webPass string) {

	app := new(handlers.Application)
	app.Auth.Username = webUser
	app.Auth.Password = webPass

	fs := http.FileServer(http.Dir("templates"))

	c.mux.Handle("/", fs)
	c.mux.HandleFunc("/status", app.BasicAuth(app.ProtectedHandler))
	c.mux.HandleFunc("/status/terminate/", app.TerminateSession)

	c.mux.HandleFunc("/test", app.TestH)

}
