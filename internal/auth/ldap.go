package auth

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

const (
	DOMAIN = "pengwin"
	DC     = "local"
)

type ClientLdap struct {
	con *ldap.Conn
}

func NewLdapClient() (*ClientLdap, error) {

	// Подключение к серверу LDAP
	c, err := ldap.Dial("tcp", viper.GetString("freeIpa.host")+":"+viper.GetString("freeIpa.port"))
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
