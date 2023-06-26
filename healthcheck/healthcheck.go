package healthcheck

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func CheckPort(host string, port int) {

	localport := strconv.Itoa(port)

	timeout := time.Duration(3 * time.Second)
	_, err := net.DialTimeout("tcp", host+":"+localport, timeout)
	if err != nil {
		//fmt.Printf("%s %s %s\n", host, "not responding", err.Error())
		log.Println("", err.Error())

	} else {
		//fmt.Printf("%s %s %s\n", host, "responding on port:", localport)
		log.Println("")
	}
}

func CmdPing(host string) string {
	var cmdstr string

	cmd := "> /dev/null && echo true || echo false"
	p := "ping -c "
	pingcount := "1 "
	cmdstr += p + pingcount + host + cmd
	// fmt.Println(cmdstr)

	os := runtime.GOOS
	switch os {
	case "windows":
		fmt.Println("Windows")
		output, err := exec.Command("cmd.exe", "/C", "ping ", host).Output()
		if err != nil {
			panic(err)
		}
		result := string(output)
		fmt.Println("result:", result)
		return result
		//rbool, _ := strconv.ParseBool(result)
		//fmt.Println("rbool:", rbool)
	case "darwin":
		fmt.Println("MacOs")
		output, err := exec.Command("/bin/bash", "-c", "ping", host).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(output)
		result := string(output)
		return result
	case "linux":
		fmt.Println("Linux")
		output, err := exec.Command("/bin/sh", "-c", cmdstr).Output()
		if err != nil {
			panic(err)
		}
		result := string(output)
		fmt.Println("result:", result)
		rbool, _ := strconv.ParseBool(result)
		fmt.Println("rbool:", rbool)
	default:
		fmt.Printf("%s.\n\n", os)
	}
	return ""
}
