package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type ClientPg struct {
	condb *sql.DB
}

func (c *ClientPg) CloseDB() {
	c.condb.Close()
}

func NewClient() (*ClientPg, error) {

	dbhost, dbport, dbusername, dbuserpass, dbname := getСredentials()
	conn, err := ConnectDB(dbhost, dbport, dbusername, dbuserpass, dbname)
	if err != nil {
		return nil, err
	}
	return &ClientPg{condb: conn}, nil
}

func getСredentials() (string, string, string, string, string) {
	dbhost, _ := os.LookupEnv("DB_HOST")
	dbport, _ := os.LookupEnv("DB_PORT")
	dbname, _ := os.LookupEnv("DB_NAME")
	dbussername, _ := os.LookupEnv("DB_USER")
	dbuserpass, _ := os.LookupEnv("DB_PASS")

	log.Println("get credentials in file .env")

	return dbhost, dbport, dbussername, dbuserpass, dbname
}

func ConnectDB(host, port, user, password, dbname string) (*sql.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to", host)
	return db, nil
}

// func (c *ClientPg) GetDeployService() (map[int]string, error) {

// 	deployServiceList := make(map[int]string)

// 	var id int
// 	var name string

// 	result, err := c.condb.Query("SELECT id, name FROM public.uds__deployed_service order by id ")
// 	if err != nil {
// 		return nil, err
// 	}

// 	for result.Next() {
// 		if err := result.Scan(&id, &name); err != nil {
// 			return nil, err
// 		}
// 		deployServiceList[id] = name
// 	}
// 	return deployServiceList, err
// }
