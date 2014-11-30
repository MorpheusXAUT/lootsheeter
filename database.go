// database
package main

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/morpheusxaut/lootsheeter/models"
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

func InitialiseDatabase(host string, port int, user string, password string, data string) {
	logger.Debugf("Trying to connect to MySQL database at %q...", net.JoinHostPort(host, strconv.Itoa(port)))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, net.JoinHostPort(host, strconv.Itoa(port)), data))
	if err != nil {
		logger.Fatalf("Failed to connect to database: [%v]", err)
		return
	}

	database = NewDatabase(db)

	logger.Debugln("Successfully connected to database, trying to ping...")

	err = database.db.Ping()
	if err != nil {
		logger.Fatalf("Failed to ping database: [%v]", err)
		return
	}

	logger.Debugln("Successfully pinged database, initialisation completed!")
}

func (db *Database) LoadCorporation(id int64) (models.Corporation, error) {
	row := db.db.QueryRow("SELECT c.id as cid, c.corporation_id AS corporation_id, c.name as corporation_name, c.ticker AS corporation_ticker FROM corporations AS c WHERE c.active = 'Y' AND c.id = ?", id)

	var cid, corporationId int64
	var corporationName, corporationTicker string

	err := row.Scan(&cid, &corporationId, &corporationName, &corporationTicker)
	if err != nil {
		return models.Corporation{}, err
	}

	return models.NewCorporation(cid, corporationId, corporationName, corporationTicker), nil
}

func (db *Database) LoadPlayer(name string) (models.Player, error) {
	row := db.db.QueryRow("SELECT p.id AS pid, p.player_id AS player_id, p.name AS player_name, p.corporation_id AS cid, p.access AS player_access FROM players AS p WHERE p.active = 'Y' AND p.name LIKE ?", name)

	var pid, playerId, cid int64
	var playerAccess int
	var playerName string

	err := row.Scan(&pid, &playerId, &playerName, &cid, &playerAccess)
	if err != nil {
		return models.Player{}, err
	}

	corp, err := db.LoadCorporation(cid)
	if err != nil {
		return models.Player{}, err
	}

	return models.NewPlayer(pid, playerId, playerName, corp, models.AccessMask(playerAccess)), nil
}
