package connections

import (
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

func InitConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass string) error {
	var err error

	Conn, err = getConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass)
	if err != nil {
		return err
	}
	return err
}

func getConnections(ipaHost, ipaUser, ipaPass, srvUser, srvPass string) (*Connections, error) {

	conIpa, err := authenticators.NewClient(ipaHost, ipaUser, ipaPass)
	if err != nil {
		// log.Fatalf("can not create freeIpa client; err: %s", err.Error())
		return nil, err
	}

	conPg, err := db.NewClient()
	if err != nil {
		// log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
		return nil, err
	}

	conSSH, err := connectors.NewClient(srvUser, srvPass)
	if err != nil {
		// log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
		return nil, err
	}

	return &Connections{IPA: conIpa, Database: conPg, SSH: conSSH}, err
}
