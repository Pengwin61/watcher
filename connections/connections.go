package connections

import (
	"log"
	"watcher/authenticators"
	"watcher/connectors"
	"watcher/db"
)

type Connections struct {
	IPA      *authenticators.Client
	Database *db.ClientPg
	SSH      *connectors.Client
}

var Conn *Connections

func InitConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass string) {

	Conn = getConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass)
}

func getConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass string) *Connections {

	conIpa, err := authenticators.NewClient(ipaHost, ipaUser, ipaPass)
	if err != nil {
		log.Fatalf("can not create freeIpa client; err: %s", err.Error())
	}

	conPg, err := db.NewClient()
	if err != nil {
		log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
	}

	conSSH, err := connectors.NewClient(srvUser, srvPass)
	if err != nil {
		log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
	}

	return &Connections{IPA: conIpa, Database: conPg, SSH: conSSH}
}
