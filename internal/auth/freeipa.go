package auth

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/spf13/viper"
)

type Employee struct {
	Username   string
	UidNumber  int
	GuidNumber int
}

type ClientFreeIpa struct {
	con *freeipa.Client
}

func NewClient() (*ClientFreeIpa, error) {
	conn, err := ConIpa()

	if err != nil {
		return nil, err
	}
	return &ClientFreeIpa{con: conn}, nil
}

func ConIpa() (*freeipa.Client, error) {

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c, err := freeipa.Connect(viper.GetString("freeIpa.host"), tspt,
		viper.GetString("freeIpa.username"), viper.GetString("freeIpa.password"))

	if err != nil {
		return nil, err
	}
	return c, nil

}
func (c *ClientFreeIpa) GetGroups(ipaGroup string) ([]string, error) {

	r, err := c.con.GroupShow(&freeipa.GroupShowArgs{Cn: *freeipa.String(ipaGroup)}, &freeipa.GroupShowOptionalArgs{})

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
	groupFreeIpaList := r.Result.MemberGroup
	return *groupFreeIpaList, nil
}

func (c *ClientFreeIpa) GetUser(ipaGroup string) ([]string, error) {
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
	if res.Result.MemberUser != nil {
		result := res.Result.MemberUser

		return *result, nil
	}

	return nil, nil
}

func (c *ClientFreeIpa) GetUserID(userlist []string) (map[string]Employee, error) {

	employeeList := map[string]Employee{}

	for _, user := range userlist {
		res2, err := c.con.UserShow(&freeipa.UserShowArgs{}, &freeipa.UserShowOptionalArgs{UID: freeipa.String(user)})

		if err != nil {
			return nil, err
		}

		employee := Employee{
			Username:   res2.Result.UID,
			UidNumber:  *res2.Result.Uidnumber,
			GuidNumber: *res2.Result.Gidnumber}

		employeeList[res2.Result.UID] = employee
	}
	return employeeList, nil
}

func (c *ClientFreeIpa) CheckUser(username string) (*string, bool, error) {
	var isAdmin bool

	res, err := c.con.UserShow(&freeipa.UserShowArgs{}, &freeipa.UserShowOptionalArgs{UID: freeipa.String(username)})

	if err != nil {
		if ipaE, ok := err.(*freeipa.Error); ok {
			log.Printf("FreeIPA error %v: %v\n", ipaE.Code, ipaE.Message)
			if ipaE.Code == freeipa.NotFoundCode {
				log.Println("(matched expected error code)")
			}
		} else {
			log.Printf("Other error: %v", err)
		}
		return nil, false, err
	}

	userGroups := res.Result.MemberofGroup

	for _, group := range *userGroups {
		if strings.Contains(group, "admins") {
			isAdmin = true
			break
		}

	}

	return &res.Result.UID, isAdmin, err
}
