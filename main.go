package main

import (
	"io"
	"log"
	"os"
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

	webClient := webapp.NewClient(params.Web.Port)
	webClient.RunWeb(params.Web.User, params.Web.Pass,
		params.Web.SslPub, params.Web.SslPriv)
  
}
