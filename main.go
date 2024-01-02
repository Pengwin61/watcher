package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"watcher/internal/configs"
	"watcher/internal/connections"

	"watcher/internal/logs"
	"watcher/internal/watch"
	"watcher/internal/webapp"

	"github.com/spf13/viper"
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

	// Read logs
	go func() {
		for err := range errCh {
			log.Println(err)
		}
	}()

	// Running Web
	go webapp.InitGin()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	log.Println("stopping application:", sign)

	log.Println("closing logfile:", viper.GetString("paths.logs"))
	defer logfile.CloseFile()

	log.Println("closing connections to database")
	defer connections.Conn.Database.CloseDB()
}
