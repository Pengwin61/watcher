package connectors

import (
	"fmt"
	"strings"
)

type User struct {
	SessionID string
	Hostname  string
	State     string
	InitTime  string
	LastTime  string
	Username  string
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

func ParseSession(stdout string) (map[string]*User, error) {

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
			SessionID: v[1],
			Hostname:  v[3],
			State:     v[4],
			InitTime:  v[5],
			LastTime:  v[10],
			Username:  v[11],
		}

		storage[tmp.Username] = tmp
	}

	return storage, err
}

func (c *Client) TerminateSession(sessionPid, host string) {
	cmdTerminated := "sudo x2goterminate-session"

	c.ExecuteCmd(cmdTerminated+""+sessionPid, host)
}
