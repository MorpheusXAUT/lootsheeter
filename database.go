// database
package main

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	database *Database
)

type Database struct {
	db *sql.DB
}

func NewDatabase(d *sql.DB) *Database {
	database := &Database{
		db: d,
	}

	return database
}

func InitialiseDatabase(host string, port int, user string, password string, data string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, net.JoinHostPort(host, strconv.Itoa(port)), data))
	if err != nil {
		return err
	}

	database = NewDatabase(db)

	err = database.db.Ping()
	if err != nil {
		return err
	}

	return nil
}
