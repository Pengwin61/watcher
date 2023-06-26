package connectors

import (
	"fmt"
	"strings"
)

type User struct {
	AgentPid, SessionPid, Port, Hostname,
	SessionState, StartDateSession, field7, UserIp,
	field9, field10, StopDateSession, UserSession,
	field13, field14, field15, field16 string
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
			AgentPid:         v[0],
			SessionPid:       v[1],
			Port:             v[2],
			Hostname:         v[3],
			SessionState:     v[4],
			StartDateSession: v[5],
			field7:           v[6],
			UserIp:           v[7],
			field9:           v[8],
			field10:          v[9],
			StopDateSession:  v[10],
			UserSession:      v[11],
			field13:          v[12],
			field14:          v[13],
			field15:          v[14],
			field16:          v[15],
		}

		storage[tmp.UserSession] = tmp
	}

	return storage, err
}

func (c *Client) TerminateSession(sessionPid, host, cmd string, conssh *Client) {

	hostsList := make(map[string]string)
	hostsList[host] = host

	c.ConnectHost(fmt.Sprint(cmd+" "+sessionPid), hostsList)
}
