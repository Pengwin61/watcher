package configs

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// func InitConfigs() Params {
// 	// loads values from .env into the system
// 	// if err := godotenv.Load("config/.env"); err != nil {
// 	// 	log.Print("No .env file found")
// 	// 	os.Exit(1)
// 	// }

// 	// loads values from settings.cfg
// 	cfg, err := ini.Load("config/settings.cfg")
// 	if err != nil {
// 		log.Printf("fail to read file: %v", err)
// 		os.Exit(1)
// 	}

// 	webPort := cfg.Section("web").Key("port").String()
// 	webUser := cfg.Section("web").Key("user").String()
// 	webPass := cfg.Section("web").Key("password").String()
// 	sslPub := cfg.Section("web").Key("ssl_public").String()
// 	sslPriv := cfg.Section("web").Key("ssl_private").String()

// 	pathHome := cfg.Section("paths").Key("home_dir").String()
// 	pathLogs := cfg.Section("paths").Key("logs").String()

// 	daysRotation := cfg.Section("maintenance").Key("home_dir_days_rotation").String()
// 	timeExpiration, err := cfg.Section("maintenance").Key("time_expiration_session").Duration()
// 	if err != nil {
// 		fmt.Println("can parse time_expiration_session to ini file")
// 	}

// 	hostIpa := cfg.Section("FreeIpa").Key("host").String()
// 	userIpa := cfg.Section("FreeIpa").Key("username").String()
// 	userPassIpa := cfg.Section("FreeIpa").Key("password").String()

// 	userPassIpa, err = decodingPassword(userPassIpa)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	groupIpa := cfg.Section("FreeIpa").Key("master_group").String()

// 	actorsUser := cfg.Section("servers").Key("username").String()
// 	actorsPaswd := cfg.Section("servers").Key("password").String()

// 	actorsPaswd, err = decodingPassword(actorsPaswd)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	softQuota := cfg.Section("UserQuota").Key("softQuota").String()
// 	hardQuota := cfg.Section("UserQuota").Key("hardQuota").String()

// 	basePath := core.CreatePath(pathHome)

// 	/* Flags */
// 	scheduleFlag := flag.String("schedule", "10m", "Delault time for updates")
// 	flag.Parse()
// 	schedule, _ := time.ParseDuration(*scheduleFlag)

// 	var params = Params{
// 		Maintenance: Maintenance{DaysRotation: daysRotation, Schedule: schedule, TimeExpiration: timeExpiration},
// 		FreeIPA:     FreeIPA{Host: hostIpa, User: userIpa, Pass: userPassIpa, Group: groupIpa},
// 		Paths:       Paths{Home: basePath, Logs: pathLogs},
// 		Servers:     Servers{User: actorsUser, Pass: actorsPaswd},
// 		UserQuota:   UserQuota{Soft: softQuota, Hard: hardQuota},
// 		Web:         Web{User: webUser, Pass: webPass, Port: webPort, SslPub: sslPub, SslPriv: sslPriv}}

// 	return params
// }

// func decodingPassword(encodePassword string) (string, error) {

// 	decodePass, err := base64.StdEncoding.DecodeString(encodePassword)
// 	if err != nil {
// 		return "", err
// 	}
// 	password := bytes.NewBuffer(decodePass).String()

// 	return password, err
// }

func InitConfigsViper() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// set default
	viper.SetDefault("web.port", "80")
	viper.SetDefault("schedule.interval", "2m")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}
}
