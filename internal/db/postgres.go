package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type ClientPg struct {
	condb *sql.DB
}

func (c *ClientPg) CloseDB() {
	c.condb.Close()
}

func NewClient() (*ClientPg, error) {

	conn, err := connectDB(viper.GetString("database.host"), viper.GetString("database.port"),
		viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.name"))
	if err != nil {
		return nil, err
	}
	return &ClientPg{condb: conn}, nil
}

func connectDB(host, port, user, password, dbname string) (*sql.DB, error) {

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
