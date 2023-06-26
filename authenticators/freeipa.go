package authenticators

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/ccin2p3/go-freeipa/freeipa"
)

type Client struct {
	con *freeipa.Client
}

func NewClient(ipaHost, ipaUser, ipaPassword string) (*Client, error) {
	conn, err := ConIpa(ipaHost, ipaUser, ipaPassword)

	if err != nil {
		return nil, err
	}
	return &Client{con: conn}, nil
}

func ConIpa(ipaHost, ipaUser, ipaPasswd string) (*freeipa.Client, error) {

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c, err := freeipa.Connect(ipaHost, tspt, ipaUser, ipaPasswd)
	if err != nil {
		return nil, err
	}
	return c, nil

}

func (c *Client) GetUser(ipaGroup string) ([]string, error) {
	res, err := c.con.GroupShow(&freeipa.GroupShowArgs{Cn: *freeipa.String(ipaGroup)}, &freeipa.GroupShowOptionalArgs{})
	if err != nil {
		if ipaE, ok := err.(*freeipa.Error); ok {
			log.Printf("FreeIPA error %v: %v\n", ipaE.Code, ipaE.Message)
			if ipaE.Code == freeipa.NotFoundCode {
				log.Println("(matched expected error code)")
			}
		} else {
			log.Printf("Other error: %v", err)
		}
		return nil, err
	}
	userFreeIpaList := res.Result.MemberUser
	return *userFreeIpaList, nil
}

func (c *Client) GetUserID(userlist []string) (map[string]int, error) {
	employee := make(map[string]int)

	for _, user := range userlist {
		res2, err := c.con.UserShow(&freeipa.UserShowArgs{}, &freeipa.UserShowOptionalArgs{UID: freeipa.String(user)})

		if err != nil {
			return nil, err
		}

		username := res2.Result.UID
		uidnumber := res2.Result.Uidnumber

		employee[username] = *uidnumber
	}

	return employee, nil
}
