package controllers

import (
	"net/http"
	"os"
	"watcher/internal/connections"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Login(c *gin.Context) {

	username := c.PostForm("username")
	pass := c.PostForm("password")

	ok := connections.Conn.LDAP.CheckUser(username, pass)
	// connections.Conn.LDAP.
	if !ok {
		c.HTML(http.StatusUnauthorized, "unauthorized.html", gin.H{})
	}

	user, isAdmin, err := connections.Conn.IPA.CheckUser(username)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "unauthorized.html", gin.H{})
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
		"isAdmin":  isAdmin,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	c.Redirect(http.StatusFound, "/")

}
