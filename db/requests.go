package db

import (
	"database/sql"
	"fmt"
)

type UserService struct {
	User_service_id int
	SrcIP           string
	State           string
	InUse           bool
	InUseDate       string
	DepSvcID        int
	DepSvcName      string
	UserID          int
	Username        string
}

func (c *ClientPg) GetNewRequest() (map[string]UserService, error) {

	var tmp UserService

	storage := map[string]UserService{}

	sqlselect := "SELECT public.uds__user_service.id, src_ip, public.uds__user_service.state, in_use , in_use_date, deployed_service_id, user_id, public.uds__deployed_service.name, public.uds_user.name"
	sqlfrom := " FROM public.uds__user_service"
	sqljoin := " left join uds__deployed_service on deployed_service_id = public.uds__deployed_service.id left join uds_user on user_id = uds_user.id"
	sqlWhere := " where public.uds__user_service.state = 'U'"
	result, err := c.condb.Query(fmt.Sprintf(sqlselect + sqlfrom + sqljoin + sqlWhere))
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		if err := result.Scan(&tmp.User_service_id, &tmp.SrcIP, &tmp.State, &tmp.InUse,
			&tmp.InUseDate, &tmp.DepSvcID, &tmp.UserID, &tmp.DepSvcName, &tmp.Username); err != nil {
			return nil, err
		}
		storage[tmp.Username] = tmp
	}

	return storage, result.Err()
}

func (c *ClientPg) GetEntity(entity string) (map[string]string, error) {

	var ip string
	var hostname string
	entityList := make(map[string]string)

	sqlStr := fmt.Sprintf("SELECT ip, hostname FROM public.%s order by ip", entity)

	result, err := c.condb.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		if err := result.Scan(&ip, &hostname); err != nil {
			return nil, err
		}
		entityList[hostname] = ip
	}

	return entityList, result.Err()

}
func (c *ClientPg) UpdateTab(User_service_id int) error {
	_, err := c.condb.Exec("update public.uds__user_service set state = $1, in_use= $2 where id = $3",
		"S", "false", User_service_id)

	return err
}

func (c *ClientPg) UpdateDB() (sql.Result, error) {
	// обновляем строку с где state U
	result, err := c.condb.Exec("update public.uds__user_service set state = $1, in_use= $2 where state = $3",
		"S", "false", "U")
	if err != nil {
		return nil, err
	}

	return result, err
}
