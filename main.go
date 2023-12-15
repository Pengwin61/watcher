package main

import (
	"log"
	"watcher/internal/configs"
	"watcher/internal/connections"

	"watcher/internal/logs"
	"watcher/internal/watch"
	"watcher/internal/webapp"
)

//

func main() {

	errCh := make(chan error)

	configs.InitConfigsViper()

	logfile := logs.InitLogs()

	err := connections.InitConnections()
	if err != nil {
		log.Fatalf("can`t create client: %s", err)
	}

	go watch.RunWatcher(errCh)

	go func() {
		for err := range errCh {
			log.Println(err)
		}
	}()

	// webClient := webapp.NewClient(params.Web.Port)
	// webClient.RunWeb(params.Web.User, params.Web.Pass,
	// 	params.Web.SslPub, params.Web.SslPriv)

	webapp.InitGin()

	defer logfile.CloseFile()
	defer connections.Conn.Database.CloseDB()
}
