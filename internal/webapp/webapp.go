package webapp

import (
	"fmt"
	"net/http"
	"watcher/internal/core"
	"watcher/internal/webapp/controllers"
	"watcher/internal/webapp/middleware"

	"github.com/gin-gonic/gin"
)

// type Client struct {
// 	con *http.Server
// 	mux *http.ServeMux
// }

// func NewClient(webPort string) *Client {
// 	mux, srv := initWeb(webPort)

// 	return &Client{con: srv, mux: mux}
// }

// func initWeb(webPort string) (*http.ServeMux, *http.Server) {
// 	mux := http.NewServeMux()
// 	srv := &http.Server{
// 		Addr:         fmt.Sprint("0.0.0.0:", webPort),
// 		Handler:      mux,
// 		IdleTimeout:  time.Minute,
// 		ReadTimeout:  10 * time.Second,
// 		WriteTimeout: 30 * time.Second,
// 	}
// 	log.Printf("starting server on %s", srv.Addr)

// 	return mux, srv
// }

// func (c *Client) RunWeb(user, pass, sslpub, sslpriv string) {

// 	c.runHandlers(user, pass)

// 	err := c.con.ListenAndServeTLS(sslpub, sslpriv)
// 	if err != nil {
// 		log.Printf("%s", err.Error())
// 	}
// }

// func (c *Client) runHandlers(webUser, webPass string) {

// 	app := new(handlers.Application)
// 	app.Auth.Username = webUser
// 	app.Auth.Password = webPass

// 	fs := http.FileServer(http.Dir("templates"))

// 	c.mux.Handle("/", fs)
// 	c.mux.HandleFunc("/status", app.BasicAuth(app.ProtectedHandler))
// 	c.mux.HandleFunc("/status/terminate/", app.TerminateSession)

// }

func InitGin() {

	r := gin.Default()

	r.Static("/css", "./templates/css")
	r.Static("/images", "./templates/images")
	r.Static("/js", "./templates/js")

	r.LoadHTMLGlob("templates/html/*.html")

	RunHandlers(r)

	err := r.Run(":7777")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running on port 7777")

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
	// r.GET("/actors", middleware.RequireAuth, statusActors)
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

// func statusActors(c *gin.Context) {

// 	user, _ := c.Get("username")
// 	isAuthorized := c.GetBool("isAuthorized")
// 	isAdmin := c.GetBool("isAdmin")

// 	c.HTML(http.StatusOK, "actors.html", gin.H{
// 		"title":        "Users Sessions",
// 		"username":     user,
// 		"isAuthorized": isAuthorized,
// 		"isAdmin":      isAdmin,
// 		"actors":       core.GetServerView(),
// 	})
// }
