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

	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Fatalf("Failed to load all players from database: [%v]", err)
		return
	}

	logger.Printf("%#+v", players)
}

func (db *Database) LoadAllPlayers() ([]models.Player, error) {
	players := make([]models.Player, 0)

	rows, err := db.db.Query("SELECT p.id AS id, p.player_id AS player_id, p.name AS player_name, p.access AS player_access, c.id AS corporation_id, c.corp_id AS corp_id, c.name AS corporation_name FROM players AS p INNER JOIN corporations AS c ON c.id = p.corporation_id WHERE p.active = 'Y'")
	if err != nil {
		return players, err
	}

	for rows.Next() {
		var id, playerId, corporationId, corpId int64
		var access int
		var playerName, corporationName string

		err := rows.Scan(&id, &playerId, &playerName, &access, &corporationId, &corpId, &corporationName)
		if err != nil {
			return players, err
		}

		players = append(players, models.NewPlayer(id, playerId, playerName, models.NewCorporation(corporationId, corpId, corporationName), models.AccessMask(access)))
	}

	return players, nil
}

func (db *Database) LoadPlayer(name string) (models.Player, error) {
	row := db.db.QueryRow("SELECT p.id AS id, p.player_id AS player_id, p.name AS player_name, p.access AS player_access, c.id AS corporation_id, c.corp_id AS corp_id, c.name AS corporation_name FROM players AS p INNER JOIN corporations AS c ON c.id = p.corporation_id WHERE p.active = 'Y' AND p.name LIKE ?", name)

	var id, playerId, corporationId, corpId int64
	var access int
	var playerName, corporationName string

	err := row.Scan(&id, &playerId, &playerName, &access, &corporationId, &corpId, &corporationName)
	if err != nil {
		return models.Player{}, err
	}

	return models.NewPlayer(id, playerId, playerName, models.NewCorporation(corporationId, corpId, corporationName), models.AccessMask(access)), nil
}
