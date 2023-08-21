package watch

import (
	"log"
	"strings"
	"time"

	"watcher/authenticators"
	"watcher/configs"
	"watcher/connectors"
	"watcher/core"
	"watcher/db"
)

func RunWatcher(params configs.Params) {

	c, err := authenticators.NewClient(params.HostIpa, params.UserIpa, params.UserPassIpa)
	if err != nil {
		log.Fatalf("can not create freeIpa client; err: %s", err.Error())
	}

	conPg, err := db.NewClient()
	if err != nil {
		log.Fatalf("can not create Postgres SQL client; err: %s", err.Error())
	}
	defer conPg.CloseDB()

	conSSH, err := connectors.NewClient(params.ActorsUser, params.ActorsPaswd)
	if err != nil {
		log.Fatalf("can not create SSH connection to hosts: %s", err.Error())
	}

	for {

		if params.Mode == "production" {

			/*  Get list uds_actors */
			actorsList, err := conPg.GetEntity("uds_actortoken")
			if err != nil {
				log.Fatalf("can not get list actors: %s", err.Error())
			}

			groupsList, err := c.GetGroups(params.GroupIpa)
			if err != nil {
				log.Printf("can not get groups list in FreeIPA; err: %s", err.Error())
			}

			err = core.CreateRootDirectory(params.PathHome, groupsList)
			if err != nil {
				log.Printf("can not create root directory; err: %s", err.Error())
			}

			for _, group := range groupsList {

				usersList, err := c.GetUser(group)
				if err != nil {
					log.Printf("can not get user list in FreeIPA; err: %s", err.Error())
				}

				if usersList != nil {

					userListID, err := c.GetUserID(usersList)
					if err != nil {
						log.Printf("can not get user list ID; err: %s", err.Error())
					}

					err = core.CreateUserDirectory(params.PathHome, group, usersList, userListID)
					if err != nil {
						log.Printf("can not create directory; err: %s", err.Error())
					}

					folderList, err := core.FindHomeFolder(params.PathHome, group)
					if err != nil {
						log.Printf("can not get list folder; err:%s", err)
					}

					diffListFolder := core.DiffDirectory(folderList, usersList)
					if diffListFolder != nil {
						err := core.DeleteFolders(params.PathHome, group, diffListFolder)
						if err != nil {
							log.Printf("can not delete folder; err:%s", err)
						}
					}

				}

				/* Удаление папки */
				err = core.DirExpired(params.PathHome, group, params.DaysRotation, usersList)
				if err != nil {
					log.Printf("can not delete directory; err: %s", err.Error())
				}
			}

			sshstdout := conSSH.ConnectHost("sudo x2golistsessions_root", actorsList)
			if sshstdout == "" {
				core.ShowSession(nil)
				time.Sleep(params.Schedule)
			}

			x2gosession, err := connectors.GetSessionX2go(sshstdout)
			if err != nil {
				if strings.Contains(err.Error(), "wrong input") {
					continue
				}
				log.Printf("list session x2go is empty: %s", err.Error())
			}

			udssession, err := conPg.GetNewRequest()
			if err != nil {
				log.Fatalf("can not; err: %s", err.Error())
			}

			core.ManageSession(x2gosession, udssession,
				conPg, conSSH, params.TimeExpiration)

		} else {

			log.Printf("APP MODE:%s", params.Mode)

		}

		time.Sleep(params.Schedule)
	}
}
