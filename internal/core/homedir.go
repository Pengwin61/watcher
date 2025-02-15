package core

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"watcher/internal/auth"
)

const folderPerm = 0700

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
			err = os.Mkdir(fullPathGroup, folderPerm)
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
	employeeList map[string]auth.Employee) error {

	dir, err := os.Open(basePath)
	if err != nil {
		return err
	}

	for _, user := range users {
		fullPathUser := filepath.Join(basePath, group, user)

		if _, err := os.Stat(fullPathUser); os.IsNotExist(err) {

			err = os.Mkdir(fullPathUser, folderPerm)
			if err != nil {
				return err
			}
			log.Printf("folder is created %s ", fullPathUser)
		}
	}
	err = changeOwner(basePath, group, employeeList)

	defer dir.Close()

	return err
}

func changeOwner(basePath, group string, employeeList map[string]auth.Employee) error {

	for username, value := range employeeList {
		fullPath := filepath.Join(basePath, group, username)
		err := os.Chown(fullPath, value.UidNumber, value.GuidNumber)
		if err != nil {
			err = errors.New("err: " + "can not change owner folder: " + err.Error())
			return err
		}
	}
	return nil
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

func DeleteFolders(basePath, group string, diffListFolder []string) (err error) {
	for _, user := range diffListFolder {
		fullPathUser := filepath.Join(basePath, group, user)

		err = os.RemoveAll(fullPathUser)
		if err != nil {
			return err
		}
		log.Printf("folder delete %s, watcher did not find the user %s in the group %s", fullPathUser, user, group)
	}
	return err
}

func FindHomeFolder(basePath, group string) ([]string, error) {
	var userList []string
	var err error
	fullPath := filepath.Join(basePath, group)

	dir, err := os.Open(fullPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer dir.Close()

	folder, err := dir.ReadDir(-1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, user := range folder {
		if !user.IsDir() {
			continue
		}
		userList = append(userList, user.Name())
	}
	return userList, err
}

func DiffDirectory(folderList, users []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range users {
		m[item] = true
	}

	for _, item := range folderList {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
