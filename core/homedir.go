package core

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"watcher/authenticators"
)

func CreatePath(pathFlag string) string {
	basePath := filepath.Join(pathFlag)
	return basePath
}

func CreateRootDirectory(basePath string, listGroups []string) error {
	dir, err := os.Open(basePath)
	if err != nil {
		return err
	}
	for _, group := range listGroups {
		fullPathGroup := filepath.Join(basePath, group)

		if _, err := os.Stat(group); os.IsNotExist(err) {
			err = os.Mkdir(fullPathGroup, 0700)
			if err != nil {
				if strings.Contains(err.Error(), "file exists") {
					continue
				}
				return err
			}
		}
	}

	defer dir.Close()
	return err
}

func CreateUserDirectory(basePath, group string, users []string,
	employeeList map[string]authenticators.Employee) error {

	dir, err := os.Open(basePath)
	if err != nil {
		return err
	}

	for _, user := range users {
		fullPathUser := filepath.Join(basePath, group, user)

		if _, err := os.Stat(fullPathUser); os.IsNotExist(err) {

			err = os.Mkdir(fullPathUser, 0700)
			if err != nil {
				return err
			}
		}
	}
	changeOwner(basePath, group, employeeList)

	defer dir.Close()

	return err
}

func changeOwner(basePath, group string, employeeList map[string]authenticators.Employee) {

	for username, value := range employeeList {
		fullPath := filepath.Join(basePath, group, username)
		err := os.Chown(fullPath, value.UidNumber, value.GuidNumber)
		if err != nil {
			log.Println("can not change owner folder:", err)
		}
	}
}

func DirExpired(basePath, group, daysRotation string, usersList []string) error {
	var err error

	days, err := strconv.Atoi(daysRotation)
	if err != nil {
		return err
	}

	daysRotationInMinuts := days * 24 * 60

	nowTime := time.Now()
	then := nowTime.Add(time.Duration(-daysRotationInMinuts) * time.Hour)

	for _, user := range usersList {
		fullPathUser := filepath.Join(basePath, group, user)

		fileInfo, err := os.Stat(fullPathUser)
		if err != nil {
			return err
		}

		dirT := fileInfo.ModTime()

		if !fileInfo.IsDir() {

			continue
		}

		if then.Unix() > dirT.Unix() {
			err = os.RemoveAll(fullPathUser)
			if err != nil {
				return err
			}
			log.Printf("folder %s delete, last modify: %s", fullPathUser, then.Truncate(time.Minute))
		}
	}
	return err
}
