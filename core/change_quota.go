package core

import (
	"log"
	"os/exec"
)

func SetQuota(softQuota, hardQuota string, userList []string) {
	for _, users := range userList {
		out, err := exec.Command("setquota", "-u", users, softQuota, hardQuota, "0", "0", "/").Output()
		if err != nil {
			panic("Is the system exactly Linux?")
		}
		log.Println("User quota:", users, "is set", hardQuota)
		log.Println(string(out))
	}
}
