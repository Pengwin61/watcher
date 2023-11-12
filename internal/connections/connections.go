package connections

import (
	"watcher/internal/authenticators"
	"watcher/internal/connectors"
	"watcher/internal/db"
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
		return nil, err
	}

	conPg, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	conSSH, err := connectors.NewClient(srvUser, srvPass)
	if err != nil {
		return nil, err
	}

	return &Connections{IPA: conIpa, Database: conPg, SSH: conSSH}, err
}
