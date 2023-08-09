package connectors

import (
	"fmt"
	"strings"
)

type User struct {
	AgentPid, SessionID, Port, Hostname,
	State, InitTime, sessionCookie, ClientIP,
	grPort, sndPort, LastTime, Username,
	ageInSec, sshfsPort, tekictrlPort, tekidataPort string
}

func chunk(input []string, chunkSize int) ([][]string, error) {
	q, reminder := len(input)/chunkSize, len(input)%chunkSize
	if reminder != 0 {
		return nil, fmt.Errorf("wrong input")
	}

	var result [][]string

	for i := 0; i < q; i++ {
		chunk := input[i*chunkSize : i*chunkSize+chunkSize]
		result = append(result, chunk)
	}

	return result, nil

}

func GetSessionX2go(stdout string) (map[string]*User, error) {

	stdout = strings.TrimSpace(stdout)
	stdout = strings.ReplaceAll(stdout, "\n", "|")
	stdout = strings.ReplaceAll(stdout, " ", "|")

	storage := map[string]*User{}

	//
	//
	// Разбираем строку |
	slc := strings.Split(stdout, "|")

	chunks, err := chunk(slc, 16)
	if err != nil {
		return nil, err
	}

	for _, v := range chunks {
		tmp := &User{
			AgentPid:      v[0],
			SessionID:     v[1],
			Port:          v[2],
			Hostname:      v[3],
			State:         v[4],
			InitTime:      v[5],
			sessionCookie: v[6],
			ClientIP:      v[7],
			grPort:        v[8],
			sndPort:       v[9],
			LastTime:      v[10],
			Username:      v[11],
			ageInSec:      v[12],
			sshfsPort:     v[13],
			tekictrlPort:  v[14],
			tekidataPort:  v[15],
		}

		storage[tmp.Username] = tmp
	}

	return storage, err
}

func (c *Client) TerminateSession(sessionPid, host, cmd string) {

	hostsList := make(map[string]string)
	hostsList[host] = host

	c.ConnectHost(fmt.Sprint(cmd+" "+sessionPid), hostsList)
}
