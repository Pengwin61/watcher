package core

import (
	"log"
	"os/exec"
)

func SetQuota(softQuota, hardQuota string, userList []string) error {
	for _, users := range userList {
		_, err := exec.Command("setquota", "-u", users, softQuota, hardQuota, "0", "0", "/").Output()
		if err != nil {
			return err
		}
		log.Println("User quota:", users, "is set", hardQuota)
	}
}
