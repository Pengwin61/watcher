package watch

import (
	"log"
	"strings"
	"time"

	"watcher/internal/configs"
	"watcher/internal/connections"
	"watcher/internal/connectors"
	"watcher/internal/core"
)

const (
	cmdListSession = "sudo x2golistsessions_root"
	cmdListActor   = "uds_actortoken"
)

func RunWatcher(params configs.Params, errCh chan error) error {

	defer close(errCh)

	for {

		/*  Get list uds_actors */
		actorsList, err := connections.Conn.Database.GetEntity(cmdListActor)
		if err != nil {
			errCh <- err
			// log.Fatalf("can not get list actors: %s", err.Error())
		}

		groupsList, err := connections.Conn.IPA.GetGroups(params.FreeIPA.Group)
		if err != nil {
			errCh <- err
			// log.Printf("can not get groups list in FreeIPA; err: %s", err.Error())
		}

		err = core.CreateRootDirectory(params.Paths.Home, groupsList)
		if err != nil {
			errCh <- err
			// log.Printf("can not create root directory; err: %s", err.Error())
		}

		for _, group := range groupsList {

			usersList, err := connections.Conn.IPA.GetUser(group)
			if err != nil {
				errCh <- err
				// log.Printf("can not get user list in FreeIPA; err: %s", err.Error())
			}

			if usersList != nil {

				userListID, err := connections.Conn.IPA.GetUserID(usersList)
				if err != nil {
					errCh <- err
					// log.Printf("can not get user list ID; err: %s", err.Error())
				}

				err = core.CreateUserDirectory(params.Paths.Home, group, usersList, userListID)
				if err != nil {
					errCh <- err
					// log.Printf("can not create directory; err: %s", err.Error())
				}

				folderList, err := core.FindHomeFolder(params.Paths.Home, group)
				if err != nil {
					errCh <- err
					// log.Printf("can not get list folder; err:%s", err)
				}

				diffListFolder := core.DiffDirectory(folderList, usersList)
				if diffListFolder != nil {
					err := core.DeleteFolders(params.Paths.Home, group, diffListFolder)
					if err != nil {
						errCh <- err
						// log.Printf("can not delete folder; err:%s", err)
					}
				}

			}

			/* Удаление папки */
			err = core.DirExpired(params.Paths.Home, group, params.DaysRotation, usersList)
			if err != nil {
				errCh <- err
				// log.Printf("can not delete directory; err: %s", err.Error())
			}
		}

		sshstdout := connections.Conn.SSH.GetSessionX2go(cmdListSession, actorsList)
		if sshstdout == "" {
			core.ShowSession(nil)
			time.Sleep(params.Schedule)
		}

		x2gosession, err := connectors.ParseSession(sshstdout)
		if err != nil {
			if strings.Contains(err.Error(), "wrong input") {
				continue
			}
			log.Printf("list session x2go is empty: %s", err.Error())
		}

		udssession, err := connections.Conn.Database.GetNewRequest()
		if err != nil {
			errCh <- err
			// log.Fatalf("can not; err: %s", err.Error())
		}

		err = core.ManageSession(x2gosession, udssession, params.TimeExpiration)
		if err != nil {
			errCh <- err
			// log.Fatalf("can not; err: %s", err.Error())
		}

		time.Sleep(params.Schedule)
	}
}
