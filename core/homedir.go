package core

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func CreatePath(pathFlag string) string {
	basePath := filepath.Join(pathFlag)
	return basePath
}

func CreateDirectory(basePath string, users []string, employeeList map[string]int) error {
	dir, err := os.Open(basePath)
	if err != nil {
		return err
	}

	var userlist []string
	userlist = append(userlist, users...)

	for _, user := range userlist {
		fullPathUser := filepath.Join(basePath, user)

		err = os.Mkdir(fullPathUser, 0700)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}
	changeOwner(basePath, employeeList)

	defer dir.Close()

	return err
}

func changeOwner(basePath string, employeeList map[string]int) {

	for key, value := range employeeList {
		fullPath := filepath.Join(basePath, key)
		e := os.Chown(fullPath, value, value)
		if e != nil {
			log.Println(e)
		}
	}
}

func DirExpired(basePath string, daysRotation string, usersList []string) error {
	var err error

	days, err := strconv.Atoi(daysRotation)
	if err != nil {
		return err
	}

	daysRotationInMinuts := days * 24

	nowTime := time.Now()
	then := nowTime.Add(time.Duration(-daysRotationInMinuts) * time.Hour)

	for _, user := range usersList {
		fullPathUser := basePath + "/" + user

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
			log.Println("Folder", user, "delete", "last modify:", then)
		}
	}
	return err
}
