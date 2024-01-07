package db

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

type UserService struct {
	DbID       int
	State      string
	InUse      bool
	InUseDate  string
	DepSvcID   int
	DepSvcName string
	UserID     int
	Username   string
}

func (c *ClientPg) GetNewRequest() (map[string]UserService, error) {

	var tmp UserService

	storage := map[string]UserService{}

	// формируем sql string
	sql, args, err := squirrel.
		Select("public.uds__user_service.id, public.uds__user_service.state, in_use , in_use_date, deployed_service_id, user_id, public.uds__deployed_service.name, public.uds_user.name").
		From("public.uds__user_service").
		Join("uds__deployed_service on deployed_service_id = public.uds__deployed_service.id").
		Join("uds_user on user_id = uds_user.id").
		Where("public.uds__user_service.state = 'U'").PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	// делаем запрос к субд
	result, err := c.condb.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		if err := result.Scan(&tmp.DbID, &tmp.State, &tmp.InUse,
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

	sql, args, err := squirrel.Select("ip, hostname").From(entity).OrderBy("ip").ToSql()
	if err != nil {
		return nil, err
	}

	result, err := c.condb.Query(sql, args...)
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
func (c *ClientPg) UpdateTab(UserServiceId int) error {

	sql, args, err := squirrel.
		Update("public.uds__user_service").Set("state", "S").Set("in_use", false).
		Where(squirrel.Eq{"id": UserServiceId}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = c.condb.Exec(sql, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}
