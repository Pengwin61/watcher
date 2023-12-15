package webapp

import (
	"fmt"
	"net/http"
	"watcher/internal/core"
	"watcher/internal/webapp/controllers"
	"watcher/internal/webapp/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitGin() {

	r := gin.Default()

	r.Static("/css", "./templates/css")
	r.Static("/images", "./templates/images")
	r.Static("/js", "./templates/js")

	r.LoadHTMLGlob("templates/html/*.html")

	RunHandlers(r)

	gin.SetMode(gin.ReleaseMode)

	if viper.GetString("web.port") == "443" {
		err := r.RunTLS(":443", viper.GetString("web.ssl.cert"), viper.GetString("web.ssl.key"))
		if err != nil {
			panic(err)
		}

	} else {
		err := r.Run(fmt.Sprint(":", viper.GetString("web.port")))
		if err != nil {
			panic(err)
		}
	}

}

func RunHandlers(r *gin.Engine) {
	r.GET("/", middleware.RequireAuth, func(c *gin.Context) {
		user := c.GetString("username")
		isAuthorized := c.GetBool("isAuthorized")
		isAdmin := c.GetBool("isAdmin")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":        "Main Page",
			"username":     user,
			"isAuthorized": isAuthorized,
			"isAdmin":      isAdmin,
		})
	})

	r.GET("/home", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})

	// sessions
	r.GET("/actors", middleware.RequireAuth, statusActors)
	r.GET("/sessions", middleware.RequireAuth, statusUsers)
	r.GET("/sessions/terminate/:id", controllers.TerminateSession)

	//Login
	r.POST("/login", controllers.Login)
	r.GET("/login", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/")
	})
}

func statusUsers(c *gin.Context) {
	user := c.GetString("username")
	isAuthorized := c.GetBool("isAuthorized")
	isAdmin := c.GetBool("isAdmin")

	c.HTML(http.StatusOK, "sessions.html", gin.H{
		"title":        "Users Sessions",
		"username":     user,
		"isAuthorized": isAuthorized,
		"isAdmin":      isAdmin,
		"users":        core.GetUsersView(),
		"persons":      core.GetPersonalView(user),
	})
}

func statusActors(c *gin.Context) {

	user, _ := c.Get("username")
	isAuthorized := c.GetBool("isAuthorized")
	isAdmin := c.GetBool("isAdmin")

	c.HTML(http.StatusOK, "actors.html", gin.H{
		"title":        "Users Sessions",
		"username":     user,
		"isAuthorized": isAuthorized,
		"isAdmin":      isAdmin,
		"actors":       core.GetServerView(),
	})
}
