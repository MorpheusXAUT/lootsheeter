// database
package main

import (
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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
	logger.Infof("Trying to connect to MySQL database at %q...", net.JoinHostPort(host, strconv.Itoa(port)))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, net.JoinHostPort(host, strconv.Itoa(port)), data))
	if err != nil {
		logger.Fatalf("Failed to connect to database: [%v]", err)
		return
	}

	database = NewDatabase(db)

	logger.Infof("Successfully connected to database, trying to ping...")

	err = database.db.Ping()
	if err != nil {
		logger.Fatalf("Failed to ping database: [%v]", err)
		return
	}

	logger.Infof("Successfully pinged database, initialisation completed!")
}

func (db *Database) LoadCorporation(id int64) (*models.Corporation, error) {
	logger.Tracef("Querying database for corporation with cid = %d...", id)

	row := db.db.QueryRow("SELECT c.id as cid, c.corporation_id AS corporation_id, c.name as corporation_name, c.ticker AS corporation_ticker FROM corporations AS c WHERE c.active = 'Y' AND c.id = ?", id)

	var cid, corporationId int64
	var corporationName, corporationTicker string

	err := row.Scan(&cid, &corporationId, &corporationName, &corporationTicker)
	if err != nil {
		return &models.Corporation{}, fmt.Errorf("Received error while scanning corporation row: [%v]", err)
	}

	return models.NewCorporation(cid, corporationId, corporationName, corporationTicker), nil
}

func (db *Database) LoadPlayer(id int64) (*models.Player, error) {
	logger.Tracef("Querying database for player with pid = %d...", id)

	row := db.db.QueryRow("SELECT p.id AS pid, p.player_id AS player_id, p.name AS player_name, p.corporation_id AS cid, p.access AS player_access FROM players AS p WHERE p.active = 'Y' AND p.id = ?", id)

	var pid, playerId, cid int64
	var playerAccess int
	var playerName string

	err := row.Scan(&pid, &playerId, &playerName, &cid, &playerAccess)
	if err != nil {
		return &models.Player{}, fmt.Errorf("Received error while scanning player row: [%v]", err)
	}

	corp, err := db.LoadCorporation(cid)
	if err != nil {
		return &models.Player{}, err
	}

	return models.NewPlayer(pid, playerId, playerName, corp, models.AccessMask(playerAccess)), nil
}

func (db *Database) LoadAllPlayers() ([]*models.Player, error) {
	logger.Tracef("Querying database for all players...")

	players := make([]*models.Player, 0)

	rows, err := db.db.Query("SELECT p.id AS pid, p.player_id AS player_id, p.name AS player_name, p.corporation_id AS cid, p.access AS player_access FROM players AS p WHERE p.active = 'Y'")
	if err != nil {
		return players, fmt.Errorf("Received error while querying for all players: [%v]", err)
	}

	for rows.Next() {
		var pid, playerId, cid int64
		var playerAccess int
		var playerName string

		err := rows.Scan(&pid, &playerId, &playerName, &cid, &playerAccess)
		if err != nil {
			return players, fmt.Errorf("Received error while scanning player rows: [%v]", err)
		}

		corp, err := db.LoadCorporation(cid)
		if err != nil {
			return players, err
		}

		players = append(players, models.NewPlayer(pid, playerId, playerName, corp, models.AccessMask(playerAccess)))
	}

	return players, nil
}

func (db *Database) LoadFleetMembers(id int64) ([]*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet members with fid = %d...", id)

	fleetMembers := make([]*models.FleetMember, 0)

	rows, err := db.db.Query("SELECT fm.id AS fmid, fm.fleet_id AS fid, fm.player_id AS pid, fm.role AS fleetmember_role, fm.site_modifier AS fleetmember_site_modifier, fm.payment_modifier AS fleetmember_payment_modifier, fm.payout AS fleetmember_payout, fm.payout_complete AS fleetmember_payout_complete FROM fleetmembers AS fm WHERE fm.fleet_id = ?", id)
	if err != nil {
		return fleetMembers, err
	}

	for rows.Next() {
		var fmid, fid, pid int64
		var fleetmemberRole, fleetmemberSiteModifier int
		var fleetmemberPaymentModifier, fleetmemberPayout float64
		var fleetmemberPayoutCompleteEnum string
		var fleetmemberPayoutComplete bool

		err = rows.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum)
		if err != nil {
			return fleetMembers, fmt.Errorf("Received error while scanning fleet member row: [%v]", err)
		}

		if strings.EqualFold(fleetmemberPayoutCompleteEnum, "y") {
			fleetmemberPayoutComplete = true
		} else {
			fleetmemberPayoutComplete = false
		}

		player, err := db.LoadPlayer(pid)
		if err != nil {
			return fleetMembers, err
		}

		fleetMembers = append(fleetMembers, models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete))
	}

	return fleetMembers, nil
}

func (db *Database) LoadFleet(id int64) (*models.Fleet, error) {
	logger.Tracef("Querying database for fleet with fid = %d...", id)

	row := db.db.QueryRow("SELECT f.id AS fid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.payout_complete AS fleet_payout_complete FROM fleets AS f WHERE f.active = 'Y' AND f.id = ?", id)

	var fid int64
	var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
	var fleetProfit, fleetLosses float64
	var fleetSitesFinished int
	var fleetStart, fleetEnd *time.Time
	var fleetPayoutComplete bool

	err := row.Scan(&fid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetPayoutCompleteEnumString)
	if err != nil {
		return &models.Fleet{}, fmt.Errorf("Received error while scanning fleet row: [%v]", err)
	}

	if strings.EqualFold(fleetPayoutCompleteEnumString, "y") {
		fleetPayoutComplete = true
	} else {
		fleetPayoutComplete = false
	}

	if fleetEnd == nil {
		fleetEnd = &time.Time{}
	}

	fleetMembers, err := db.LoadFleetMembers(fid)
	if err != nil {
		return &models.Fleet{}, err
	}

	fleet := models.NewFleet(fid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetPayoutComplete)

	for _, member := range fleetMembers {
		err = fleet.AddMember(member)
		if err != nil {
			return &models.Fleet{}, err
		}
	}

	return fleet, nil
}

func (db *Database) LoadAllFleets() ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets...")

	fleets := make([]*models.Fleet, 0)

	rows, err := db.db.Query("SELECT f.id AS fid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.payout_complete AS fleet_payout_complete FROM fleets AS f WHERE f.active = 'Y'")
	if err != nil {
		return fleets, fmt.Errorf("Received error while querying for all fleets: [%v]", err)
	}

	for rows.Next() {
		var fid int64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
		var fleetProfit, fleetLosses float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetPayoutCompleteEnumString)
		if err != nil {
			return fleets, fmt.Errorf("Received error while scanning fleet rows: [%v]", err)
		}

		if fleetEnd == nil {
			fleetEnd = &time.Time{}
		}

		if strings.EqualFold(fleetPayoutCompleteEnumString, "y") {
			fleetPayoutComplete = true
		} else {
			fleetPayoutComplete = false
		}

		fleetMembers, err := db.LoadFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetPayoutComplete)

		for _, member := range fleetMembers {
			err = fleet.AddMember(member)
			if err != nil {
				return fleets, err
			}
		}

		fleets = append(fleets, fleet)
	}

	return fleets, nil
}
