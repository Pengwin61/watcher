package auth

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

const (
	HOST   = "ipa.pengwin.local"
	DOMAIN = "pengwin"
	DC     = "local"
	PORT   = "389"
)

type ClientLdap struct {
	con *ldap.Conn
}

func NewLdapClient() (*ClientLdap, error) {

	// Подключение к серверу LDAP
	c, err := ldap.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		log.Println("Ошибка подключения к серверу FreeIPA:", err)
		return nil, err
	}
	return &ClientLdap{con: c}, nil
}

func (c *ClientLdap) CheckUser(username string, password string) bool {
	// Проверка логина и пароля пользователя
	userDN := fmt.Sprintf("uid=%s,cn=users,cn=accounts,dc=%s,dc=%s", username, DOMAIN, DC)
	err := c.con.Bind(userDN, password)
	if err != nil {
		log.Println("Ошибка проверки логина и пароля:", err)
		return false
	}
	log.Printf("Пользователь %s успешно аутентифицирован", username)

	return true
}
