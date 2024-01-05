package controllers

import (
	"log"
	"net/http"
	"strings"
	"watcher/internal/connections"
	"watcher/internal/core"
	"watcher/internal/utils"

	"github.com/gin-gonic/gin"
)

func TerminateSession(c *gin.Context) {

	user := strings.SplitAfterN(c.Param("id"), "-", 2)
	u := strings.TrimRight(user[0], "-")

	for k, v := range core.GetUsersView() {
		if u != v.Username {
			continue
		} else {
			connections.Conn.SSH.TerminateSession(v.SessionID, v.Hostname)
			log.Printf("the session %s was terminated by the administrator", user)

			err := connections.Conn.Database.UpdateTab(v.DbID)
			if err != nil {
				log.Println(err)
			}

			core.SetUserView(utils.RemoveSlice(core.GetUsersView(), k))
		}
	}

	c.Redirect(http.StatusFound, "/sessions")
}
