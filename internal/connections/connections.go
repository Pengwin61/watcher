package connections

import (
	"watcher/internal/auth"
	"watcher/internal/connectors"
	"watcher/internal/db"
)

type Connections struct {
	IPA      *auth.ClientFreeIpa
	Database *db.ClientPg
	SSH      *connectors.Client
	LDAP     *auth.ClientLdap
}

var Conn *Connections

func InitConnections() error {
	var err error

	Conn, err = getConnections()
	if err != nil {
		return err
	}
	return err
}

func getConnections() (*Connections, error) {

	conIpa, err := auth.NewClient()
	if err != nil {
		return nil, err
	}

	conLdap, err := auth.NewLdapClient()
	if err != nil {
		return nil, err
	}

	conPg, err := db.NewClient()
	if err != nil {
		return nil, err
	}

	conSSH, err := connectors.NewClient()
	if err != nil {
		return nil, err
	}

	return &Connections{IPA: conIpa, Database: conPg, SSH: conSSH, LDAP: conLdap}, err
}
