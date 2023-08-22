package configs

import "time"

type Params struct {
	Web
	Maintenance
	FreeIPA
	Paths
	UserQuota
	Servers
}

type Web struct {
	User    string
	Pass    string
	SslPub  string
	SslPriv string
	Port    string
}

type Maintenance struct {
	DaysRotation   string
	Schedule       time.Duration
	TimeExpiration time.Duration
}

type FreeIPA struct {
	Host  string
	User  string
	Pass  string
	Group string
}

type Paths struct {
	Home string
	Logs string
}

type UserQuota struct {
	Soft string
	Hard string
}

type Servers struct {
	User string
	Pass string
}
