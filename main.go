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

}
