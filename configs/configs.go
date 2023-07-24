package configs

import (
	"flag"
	"log"
	"os"
	"time"
	"watcher/core"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"
)

type Params struct {
	Mode, Domain, PathHome, PathLogs, DaysRotation, HostIpa,
	UserIpa, UserPassIpa, GroupIpa, ActorsUser, ActorsPaswd,
	SoftQuota, HardQuota, WebIp, WebUser, WebPass,
	SslPub, SslPriv string
	Schedule time.Duration
}

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

	webIp := cfg.Section("web").Key("port").String()
	webUser := cfg.Section("web").Key("user").String()
	webPass := cfg.Section("web").Key("password").String()
	sslPub := cfg.Section("web").Key("ssl_public").String()
	sslPriv := cfg.Section("web").Key("ssl_private").String()

	pathHome := cfg.Section("paths").Key("home_dir").String()
	pathLogs := cfg.Section("paths").Key("logs").String()

	// pathTest := cfg.Section("paths").Key("test_dir").ValueWithShadows()
	// fmt.Println(pathTest)

	daysRotation := cfg.Section("maintenance").Key("home_dir_days_rotation").String()

	hostIpa := cfg.Section("FreeIpa").Key("host").String()
	userIpa := cfg.Section("FreeIpa").Key("username").String()
	userPassIpa := cfg.Section("FreeIpa").Key("password").String()
	groupIpa := cfg.Section("FreeIpa").Key("user_group").String()

	actorsUser := cfg.Section("servers").Key("username").String()
	actorsPaswd := cfg.Section("servers").Key("password").String()

	softQuota := cfg.Section("UserQuota").Key("softQuota").String()
	hardQuota := cfg.Section("UserQuota").Key("hardQuota").String()

	basePath := core.CreatePath(pathHome)

	/* Flags */
	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")
	flag.Parse()
	schedule, _ := time.ParseDuration(*scheduleFlag)

	var params = Params{Mode: mode, Domain: domain, PathHome: basePath,
		PathLogs: pathLogs, DaysRotation: daysRotation, HostIpa: hostIpa,
		UserIpa: userIpa, UserPassIpa: userPassIpa, GroupIpa: groupIpa,
		ActorsUser: actorsUser, ActorsPaswd: actorsPaswd, SoftQuota: softQuota,
		HardQuota: hardQuota, WebIp: webIp, WebUser: webUser, WebPass: webPass,
		SslPub: sslPub, SslPriv: sslPriv, Schedule: schedule}

	return params
}
