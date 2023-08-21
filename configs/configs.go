package configs

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"watcher/core"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

func InitConfigs() Params {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}

	// loads values from settings.cfg
	cfg, err := ini.Load("settings.cfg")
	if err != nil {
		log.Printf("fail to read file: %v", err)
		os.Exit(1)
	}

	mode := cfg.Section("").Key("app_mode").String()
	domain := cfg.Section("").Key("domain").String()

	webPort, err := cfg.Section("web").Key("port").Int()
	if err != nil {
		fmt.Println("can parse port to ini file")
	}
	webUser := cfg.Section("web").Key("user").String()
	webPass := cfg.Section("web").Key("password").String()
	sslPub := cfg.Section("web").Key("ssl_public").String()
	sslPriv := cfg.Section("web").Key("ssl_private").String()

	pathHome := cfg.Section("paths").Key("home_dir").String()
	pathLogs := cfg.Section("paths").Key("logs").String()

	daysRotation := cfg.Section("maintenance").Key("home_dir_days_rotation").String()
	timeExpiration, err := cfg.Section("maintenance").Key("time_expiration_session").Duration()
	if err != nil {
		fmt.Println("can parse time_expiration_session to ini file")
	}

	hostIpa := cfg.Section("FreeIpa").Key("host").String()
	userIpa := cfg.Section("FreeIpa").Key("username").String()
	userPassIpa := cfg.Section("FreeIpa").Key("password").String()

	userPassIpa, err = decodingPassword(userPassIpa)
	if err != nil {
		fmt.Println(err.Error())
	}

	groupIpa := cfg.Section("FreeIpa").Key("master_group").String()

	actorsUser := cfg.Section("servers").Key("username").String()
	actorsPaswd := cfg.Section("servers").Key("password").String()

	actorsPaswd, err = decodingPassword(actorsPaswd)
	if err != nil {
		fmt.Println(err.Error())
	}

	softQuota := cfg.Section("UserQuota").Key("softQuota").String()
	hardQuota := cfg.Section("UserQuota").Key("hardQuota").String()

	basePath := core.CreatePath(pathHome)

	/* Flags */
	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")
	flag.Parse()
	schedule, _ := time.ParseDuration(*scheduleFlag)

	var params = Params{
		Maintenance: Maintenance{DaysRotation: daysRotation, Mode: mode, Domain: domain,
			Schedule: schedule, TimeExpiration: timeExpiration},
		FreeIPA:   FreeIPA{Host: hostIpa, User: userIpa, Pass: userPassIpa, Group: groupIpa},
		Paths:     Paths{Home: basePath, Logs: pathLogs},
		Servers:   Servers{User: actorsUser, Pass: actorsPaswd},
		UserQuota: UserQuota{Soft: softQuota, Hard: hardQuota},
		Web:       Web{User: webUser, Pass: webPass, Port: webPort, SslPub: sslPub, SslPriv: sslPriv}}

	return params
}

func decodingPassword(encodePassword string) (string, error) {

	decodePass, err := base64.StdEncoding.DecodeString(encodePassword)
	if err != nil {
		return "", err
	}
	password := bytes.NewBuffer(decodePass).String()

	return password, err
}
