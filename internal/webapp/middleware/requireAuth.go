package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"watcher/internal/connections"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {

	//Get the cookie from the request

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println(err)
	}
	if tokenString != "" {
		// Decode-Validate
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			c.HTML(http.StatusUnauthorized, "unauthorized.html", gin.H{
				"error": "unauthorized, token is not valid",
			})
		}

		// Check the expiration
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.HTML(http.StatusUnauthorized, "unauthorized.html", gin.H{
					"error": "unauthorized, token is expired",
				})
			}
		}
		raw := token.Claims.(jwt.MapClaims)["username"]

		user, isAdmin, err := connections.Conn.IPA.CheckUser(raw.(string))
		if err != nil {
			c.HTML(http.StatusUnauthorized, "unauthorized.html", gin.H{
				"error": "unauthorized, user is not found in LDAP",
			})
		}

		ok := strings.Contains(raw.(string), *user)
		if ok {
			c.Set("username", *user)
			c.Set("isAdmin", isAdmin)
			c.Set("isAuthorized", true)
		}

		c.Next()
	}
}
