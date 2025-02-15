package connectors

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	con *ssh.ClientConfig
}

func NewClient() (*Client, error) {

	config := &ssh.ClientConfig{
		User: viper.GetString("ssh.username"),
		Auth: []ssh.AuthMethod{
			ssh.Password(viper.GetString("ssh.password")),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &Client{con: config}, nil
}

func (c *Client) GetSessionX2go(cmd string, actorlist map[string]string) string {

	var fullResult string
	results := make(chan string, 10)
	timeout := time.After(20 * time.Second)

	for _, ip := range actorlist {
		go func(ip string) {
			results <- c.ExecuteCmd(cmd, ip)
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

func (c *Client) ExecuteCmd(cmd, hostname string) string {

	var res string
	var stdoutBuf bytes.Buffer

	conn, err := ssh.Dial("tcp", hostname+":22", c.con)
	if err != nil {
		log.Printf("host is not available:%s\n", err.Error())
		return res
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Println("can`t create session:", err.Error())
	}

	defer session.Close()

	session.Stdout = &stdoutBuf
	err = session.Run(cmd)
	if err != nil {
		log.Printf("can`t run cmd: %s", err.Error())
	}

	res = stdoutBuf.String()

	return res
}
