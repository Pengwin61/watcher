package main

import (
	"log"
	"watcher/configs"
	"watcher/connections"

	"watcher/logs"
	"watcher/watch"
	"watcher/webapp"
)

//

func main() {
	errCh := make(chan error)

	params := configs.InitConfigs()

	logfile := logs.InitLogs(params.Paths.Logs)
	defer logfile.CloseFile()

	err := connections.InitConnections(params.FreeIPA.Host, params.FreeIPA.User, params.FreeIPA.Pass,
		params.Servers.User, params.Servers.Pass)
	if err != nil {
		log.Fatalf("can`t create client: %s", err)
	}

	defer connections.Conn.Database.CloseDB()

	go watch.RunWatcher(params, errCh)
	if err != nil {
		log.Println(<-errCh)
	}

	webClient := webapp.NewClient(params.Web.Port)
	webClient.RunWeb(params.Web.User, params.Web.Pass,
		params.Web.SslPub, params.Web.SslPriv)

}
