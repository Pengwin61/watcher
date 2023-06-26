package connectors

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	con *ssh.ClientConfig
}

func NewClient(username, password string) (*Client, error) {

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &Client{con: config}, nil
}

func (c *Client) ConnectHost(cmd string, actorlist map[string]string) string {

	var fullResult string
	results := make(chan string, 10)
	timeout := time.After(30 * time.Second)

	for _, ip := range actorlist {
		go func(ip string) {
			results <- executeCmd(cmd, ip, c.con)
		}(ip)

	}
	for i := 0; i < len(actorlist); i++ {
		select {
		case res := <-results:
			fullResult += res
		case <-timeout:
			fmt.Println("Timed out:")
			// return
		}
	}

	fullResult = strings.ReplaceAll(fullResult, "\n", " ")

	return fullResult
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {

	conn, _ := ssh.Dial("tcp", hostname+":22", config)
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	res := stdoutBuf.String()

	// return hostname + ": " + stdoutBuf.String()
	return res
}

// func (c *Client) ConnectHost(cmd string, hosts []string) string {
//
//
// 	var fullResult string
// 	results := make(chan string, 10)
// 	timeout := time.After(30 * time.Second)

// 	for _, hostname := range hosts {
// 		go func(hostname string) {
// 			results <- executeCmd(cmd, hostname, c.con)
// 		}(hostname)
// 	}

// 	for i := 0; i < len(hosts); i++ {
// 		select {
// 		case res := <-results:
// 			fullResult += res
// 		case <-timeout:
// 			fmt.Println("Timed out:")
// 			// return
// 		}
// 	}
// 	fullResult = strings.ReplaceAll(fullResult, "\n", " ")

// 	return fullResult
// }

//
//
// func CreateSSH(username, password, cmd string, hosts []string) string {
// 	var fullResult string

// 	results := make(chan string, 10)

// 	timeout := time.After(30 * time.Second)
// 	// cmd := "x2golistsessions_root"

// 	config := &ssh.ClientConfig{
// 		User: username,
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password(password),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	for _, hostname := range hosts {
// 		go func(hostname string) {
// 			results <- executeCmd(cmd, hostname, config)
// 		}(hostname)
// 	}

// 	for i := 0; i < len(hosts); i++ {
// 		select {
// 		case res := <-results:
// 			fullResult += res
// 		case <-timeout:
// 			fmt.Println("Timed out:")
// 			// return
// 		}
// 	}
// 	fullResult = strings.ReplaceAll(fullResult, "\n", " ")

// 	return fullResult
// }
