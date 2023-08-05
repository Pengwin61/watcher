package db

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

type Person struct {
	Id                int
	SrcIP             string
	State             string
	InUse             bool
	InUseDate         string
	DeployServiceID   int
	UserID            int
	DeployServiceName string
	Username          string
}

func (c *ClientPg) ReqTest() {

	var tmp Person
	fmt.Println("TMP:", tmp)

	res := squirrel.Select("public.uds__user_service.id, src_ip, public.uds__user_service.state, in_use , in_use_date, deployed_service_id, user_id, public.uds__deployed_service.name, public.uds_user.name").
		From("public.uds__user_service").
		Join("uds__deployed_service on deployed_service_id = public.uds__deployed_service.id").
		Join("uds_user on user_id = uds_user.id").
		Where("public.uds__user_service.state = 'U'")

	fmt.Println("RES:", res)

}
