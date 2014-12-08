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

func (db *Database) LoadCorporationFromName(name string) (*models.Corporation, error) {
	logger.Tracef("Querying database for corporation with name = %s...", name)

	row := db.db.QueryRow("SELECT c.id as cid, c.corporation_id AS corporation_id, c.name as corporation_name, c.ticker AS corporation_ticker FROM corporations AS c WHERE c.active = 'Y' AND c.name LIKE ?", name)

	var cid, corporationId int64
	var corporationName, corporationTicker string

	err := row.Scan(&cid, &corporationId, &corporationName, &corporationTicker)
	if err != nil {
		return &models.Corporation{}, fmt.Errorf("Received error while scanning corporation name row: [%v]", err)
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

func (db *Database) LoadPlayerFromName(name string) (*models.Player, error) {
	logger.Tracef("Querying database for player with player_name = %q...", name)

	row := db.db.QueryRow("SELECT p.id AS pid, p.player_id AS player_id, p.name AS player_name, p.corporation_id AS cid, p.access AS player_access FROM players AS p WHERE p.active = 'Y' AND p.name LIKE ?", name)

	var pid, playerId, cid int64
	var playerAccess int
	var playerName string

	err := row.Scan(&pid, &playerId, &playerName, &cid, &playerAccess)
	if err != nil {
		return &models.Player{}, fmt.Errorf("Received error while scanning player name row: [%v]", err)
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

func (db *Database) LoadAvailablePlayers(fleedId int64, corporationId int64) ([]*models.Player, error) {
	logger.Tracef("Querying database for available players with cid = %d...", corporationId)

	players := make([]*models.Player, 0)

	rows, err := db.db.Query("SELECT p.id AS pid, p.player_id AS player_id, p.name AS player_name, p.corporation_id AS cid, p.access AS player_access FROM players AS p WHERE p.active = 'Y' AND p.corporation_id = ? AND p.id NOT IN (SELECT player_id FROM fleetmembers WHERE fleet_id = ?)", corporationId, fleedId)
	if err != nil {
		return players, fmt.Errorf("Received error while querying for available players: [%v]", err)
	}

	for rows.Next() {
		var pid, playerId, cid int64
		var playerAccess int
		var playerName string

		err := rows.Scan(&pid, &playerId, &playerName, &cid, &playerAccess)
		if err != nil {
			return players, fmt.Errorf("Received error while scanning available player rows: [%v]", err)
		}

		corp, err := db.LoadCorporation(cid)
		if err != nil {
			return players, err
		}

		players = append(players, models.NewPlayer(pid, playerId, playerName, corp, models.AccessMask(playerAccess)))
	}

	return players, nil
}

func (db *Database) LoadFleetMember(fleetId int64, id int64) (*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet member with fid = %d and pid = %d...", fleetId, id)

	row := db.db.QueryRow("SELECT fm.id AS fmid, fm.fleet_id AS fid, fm.player_id AS pid, fm.role AS fleetmember_role, fm.site_modifier AS fleetmember_site_modifier, fm.payment_modifier AS fleetmember_payment_modifier, fm.payout AS fleetmember_payout, fm.payout_complete AS fleetmember_payout_complete, fm.report_id AS rid FROM fleetmembers AS fm WHERE fm.fleet_id = ? AND fm.id = ?", fleetId, id)

	var fmid, fid, pid, rid int64
	var sqlRid sql.NullInt64
	var fleetmemberRole, fleetmemberSiteModifier int
	var fleetmemberPaymentModifier, fleetmemberPayout float64
	var fleetmemberPayoutCompleteEnum string
	var fleetmemberPayoutComplete bool

	err := row.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum, &sqlRid)
	if err != nil {
		return &models.FleetMember{}, fmt.Errorf("Received error while scanning fleet member row: [%v]", err)
	}

	if strings.EqualFold(fleetmemberPayoutCompleteEnum, "y") {
		fleetmemberPayoutComplete = true
	} else {
		fleetmemberPayoutComplete = false
	}

	if sqlRid.Valid {
		rid = sqlRid.Int64
	} else {
		rid = -1
	}

	player, err := db.LoadPlayer(pid)
	if err != nil {
		return &models.FleetMember{}, err
	}

	return models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid), nil
}

func (db *Database) LoadAllFleetMembers(fleetId int64) ([]*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet members with fid = %d...", fleetId)

	fleetMembers := make([]*models.FleetMember, 0)

	rows, err := db.db.Query("SELECT fm.id AS fmid, fm.fleet_id AS fid, fm.player_id AS pid, fm.role AS fleetmember_role, fm.site_modifier AS fleetmember_site_modifier, fm.payment_modifier AS fleetmember_payment_modifier, fm.payout AS fleetmember_payout, fm.payout_complete AS fleetmember_payout_complete, fm.report_id AS rid FROM fleetmembers AS fm WHERE fm.fleet_id = ?", fleetId)
	if err != nil {
		return fleetMembers, err
	}

	for rows.Next() {
		var fmid, fid, pid, rid int64
		var sqlRid sql.NullInt64
		var fleetmemberRole, fleetmemberSiteModifier int
		var fleetmemberPaymentModifier, fleetmemberPayout float64
		var fleetmemberPayoutCompleteEnum string
		var fleetmemberPayoutComplete bool

		err = rows.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum, &sqlRid)
		if err != nil {
			return fleetMembers, fmt.Errorf("Received error while scanning fleet member rows: [%v]", err)
		}

		if strings.EqualFold(fleetmemberPayoutCompleteEnum, "y") {
			fleetmemberPayoutComplete = true
		} else {
			fleetmemberPayoutComplete = false
		}

		if sqlRid.Valid {
			rid = sqlRid.Int64
		} else {
			rid = -1
		}

		player, err := db.LoadPlayer(pid)
		if err != nil {
			return fleetMembers, err
		}

		fleetMembers = append(fleetMembers, models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid))
	}

	return fleetMembers, nil
}

func (db *Database) LoadFleet(id int64) (*models.Fleet, error) {
	logger.Tracef("Querying database for fleet with fid = %d...", id)

	row := db.db.QueryRow("SELECT f.id AS fid, f.corporation_id AS cid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.corporation_payout AS fleet_corporation_payout, f.payout_complete AS fleet_payout_complete, f.report_id AS rid FROM fleets AS f WHERE f.active = 'Y' AND f.id = ?", id)

	var fid, cid, rid int64
	var sqlRid sql.NullInt64
	var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
	var fleetProfit, fleetLosses, fleetCorporationPayout float64
	var fleetSitesFinished int
	var fleetStart, fleetEnd *time.Time
	var fleetPayoutComplete bool

	err := row.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &sqlRid)
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

	if sqlRid.Valid {
		rid = sqlRid.Int64
	} else {
		rid = -1
	}

	fleetMembers, err := db.LoadAllFleetMembers(fid)
	if err != nil {
		return &models.Fleet{}, err
	}

	fleet := models.NewFleet(fid, cid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

	rows, err := db.db.Query("SELECT f.id AS fid, f.corporation_id AS cid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.corporation_payout AS fleet_corporation_payout, f.payout_complete AS fleet_payout_complete, f.report_id AS rid FROM fleets AS f WHERE f.active = 'Y'")
	if err != nil {
		return fleets, fmt.Errorf("Received error while querying for all fleets: [%v]", err)
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &sqlRid)
		if err != nil {
			return fleets, fmt.Errorf("Received error while scanning fleet rows: [%v]", err)
		}

		if fleetEnd == nil {
			fleetEnd = &time.Time{}
		}

		if sqlRid.Valid {
			rid = sqlRid.Int64
		} else {
			rid = -1
		}

		if strings.EqualFold(fleetPayoutCompleteEnumString, "y") {
			fleetPayoutComplete = true
		} else {
			fleetPayoutComplete = false
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, cid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadAllFleetsFromCorpId(corporationId int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets with cid = %d...", corporationId)

	fleets := make([]*models.Fleet, 0)

	rows, err := db.db.Query("SELECT f.id AS fid, f.corporation_id AS cid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.corporation_payout AS fleet_corporation_payout, f.payout_complete AS fleet_payout_complete, f.report_id AS rid FROM fleets AS f WHERE f.active = 'Y' AND f.corporation_id = ?", corporationId)
	if err != nil {
		return fleets, fmt.Errorf("Received error while querying for all corporation fleets: [%v]", err)
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &sqlRid)
		if err != nil {
			return fleets, fmt.Errorf("Received error while scanning corporation fleet rows: [%v]", err)
		}

		if fleetEnd == nil {
			fleetEnd = &time.Time{}
		}

		if sqlRid.Valid {
			rid = sqlRid.Int64
		} else {
			rid = -1
		}

		if strings.EqualFold(fleetPayoutCompleteEnumString, "y") {
			fleetPayoutComplete = true
		} else {
			fleetPayoutComplete = false
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, cid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadAllReportFleets(recordId int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets with rid = %d...", recordId)

	fleets := make([]*models.Fleet, 0)

	rows, err := db.db.Query("SELECT f.id AS fid, f.corporation_id AS cid, f.name as fleet_name, f.system AS fleet_system, f.system_nickname AS fleet_system_nickname, f.profit AS fleet_profit, f.losses AS fleet_losses, f.sites_finished AS fleet_sites_finished, f.`start` AS fleet_start, f.`end` AS fleet_end, f.corporation_payout AS fleet_corporation_payout, f.payout_complete AS fleet_payout_complete, f.report_id AS rid FROM fleets AS f WHERE f.report_id = ? AND f.active = 'Y'", recordId)
	if err != nil {
		return fleets, fmt.Errorf("Received error while querying for all report fleets: [%v]", err)
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &sqlRid)
		if err != nil {
			return fleets, fmt.Errorf("Received error while scanning report fleet rows: [%v]", err)
		}

		if fleetEnd == nil {
			fleetEnd = &time.Time{}
		}

		if sqlRid.Valid {
			rid = sqlRid.Int64
		} else {
			rid = -1
		}

		if strings.EqualFold(fleetPayoutCompleteEnumString, "y") {
			fleetPayoutComplete = true
		} else {
			fleetPayoutComplete = false
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, cid, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadReport(id int64) (*models.Report, error) {
	logger.Tracef("Querying database for report with rid = %d...", id)

	row := db.db.QueryRow("SELECT r.id AS rid, r.created_by AS record_created_b, r.total_payout AS record_total_payout, r.startrange AS record_startrange, r.endrange AS record_endrange, r.payout_complete AS record_payout_complete FROM reports AS r")

	var rid, pid int64
	var recordTotalPayout float64
	var recordPayoutCompleteEnumString string
	var recordStartrange, recordEndrange time.Time
	var recordPayoutComplete bool

	err := row.Scan(&rid, &pid, &recordTotalPayout, &recordStartrange, &recordEndrange, &recordPayoutCompleteEnumString)
	if err != nil {
		return &models.Report{}, fmt.Errorf("Received error while scanning report row: [%v]", err)
	}

	if strings.EqualFold(recordPayoutCompleteEnumString, "y") {
		recordPayoutComplete = true
	} else {
		recordPayoutComplete = false
	}

	fleets, err := database.LoadAllReportFleets(rid)
	if err != nil {
		return &models.Report{}, err
	}

	player, err := database.LoadPlayer(pid)
	if err != nil {
		return &models.Report{}, err
	}

	report := models.NewReport(rid, recordTotalPayout, recordStartrange, recordEndrange, recordPayoutComplete, player, fleets)

	return report, nil
}

func (db *Database) LoadAllReports() ([]*models.Report, error) {
	logger.Tracef("Querying database for all reports...")

	reports := make([]*models.Report, 0)

	rows, err := db.db.Query("SELECT r.id AS rid, r.created_by AS record_created_b, r.total_payout AS record_total_payout, r.startrange AS record_startrange, r.endrange AS record_endrange, r.payout_complete AS record_payout_complete FROM reports AS r")
	if err != nil {
		return reports, fmt.Errorf("Received error while querying for all reports: [%v]", err)
	}

	for rows.Next() {
		var rid, pid int64
		var recordTotalPayout float64
		var recordPayoutCompleteEnumString string
		var recordStartrange, recordEndrange time.Time
		var recordPayoutComplete bool

		err := rows.Scan(&rid, &pid, &recordTotalPayout, &recordStartrange, &recordEndrange, &recordPayoutCompleteEnumString)
		if err != nil {
			return reports, fmt.Errorf("Received error while scanning report rows: [%v]", err)
		}

		if strings.EqualFold(recordPayoutCompleteEnumString, "y") {
			recordPayoutComplete = true
		} else {
			recordPayoutComplete = false
		}

		fleets, err := database.LoadAllReportFleets(rid)
		if err != nil {
			return reports, err
		}

		player, err := database.LoadPlayer(pid)
		if err != nil {
			return reports, err
		}

		report := models.NewReport(rid, recordTotalPayout, recordStartrange, recordEndrange, recordPayoutComplete, player, fleets)

		reports = append(reports, report)
	}

	return reports, nil
}

func (db *Database) SaveFleet(fleet *models.Fleet) (*models.Fleet, error) {
	logger.Tracef("Saving fleet #%d to database...", fleet.Id)

	var fleetPayoutCompleteEnumString string

	if fleet.PayoutComplete {
		fleetPayoutCompleteEnumString = "Y"
	} else {
		fleetPayoutCompleteEnumString = "N"
	}

	_, err := db.LoadFleet(fleet.Id)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO fleets(name, corporation_id, system, system_nickname, profit, losses, sites_finished, start, end, corporation_payout, payout_complete) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", fleet.Name, fleet.CorporationId, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleet.EndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString)
		if err != nil {
			return fleet, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fleet, err
		}

		fleet.Id = id
	} else {
		_, err := db.db.Exec("UPDATE fleets SET name=?, corporation_id = ?, system=?, system_nickname=?, profit=?, losses=?, sites_finished=?, start=?, end=?, corporation_payout=?, payout_complete=? WHERE id=?", fleet.Name, fleet.CorporationId, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleet.EndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString, fleet.Id)
		if err != nil {
			return fleet, err
		}
	}

	for _, member := range fleet.Members {
		m, err := db.SaveFleetMember(fleet.Id, member)
		if err != nil {
			return fleet, err
		}

		member = m
	}

	return fleet, nil
}

func (db *Database) SaveFleetMember(fleetId int64, member *models.FleetMember) (*models.FleetMember, error) {
	logger.Tracef("Saving fleet member #%d to database...", member.Id)

	var fleetmemberPayoutCompleteEnum string

	if member.PayoutComplete {
		fleetmemberPayoutCompleteEnum = "Y"
	} else {
		fleetmemberPayoutCompleteEnum = "N"
	}

	_, err := db.LoadFleetMember(fleetId, member.Id)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO fleetmembers(fleet_id, player_id, role, site_modifier, payment_modifier, payout, payout_complete) VALUES (?, ?, ?, ?, ?, ?, ?)", fleetId, member.Player.Id, member.Role, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum)
		if err != nil {
			return member, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return member, err
		}

		member.Id = id
	} else {
		_, err := db.db.Exec("UPDATE fleetmembers SET fleet_id=?, player_id=?, role=?, site_modifier=?, payment_modifier=?, payout=?, payout_complete=? WHERE id=?", fleetId, member.Player.Id, member.Role, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum, member.Id)
		if err != nil {
			return member, err
		}
	}

	return member, nil
}

func (db *Database) SavePlayer(player *models.Player) (*models.Player, error) {
	logger.Tracef("Saving player #%d to database...", player.Id)

	_, err := db.LoadPlayer(player.Id)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO players(player_id, name, corporation_id, access) VALUES (?, ?, ?, ?)", player.PlayerId, player.Name, player.Corporation.Id, player.AccessMask)
		if err != nil {
			return player, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return player, err
		}

		player.Id = id
	} else {
		_, err := db.db.Exec("UPDATE players SET player_id=?, name=?, corporation_id=?, access=? WHERE id=?", player.PlayerId, player.Name, player.Corporation.Id, player.AccessMask, player.Id)
		if err != nil {
			return player, err
		}
	}

	return player, nil
}

func (db *Database) SaveCorporation(corporation *models.Corporation) (*models.Corporation, error) {
	logger.Tracef("Saving corporation #%d to database...", corporation.Id)

	_, err := db.LoadCorporation(corporation.Id)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO corporations(corporation_id, name, ticker) VALUES (?, ?, ?)", corporation.CorpId, corporation.Name, corporation.Ticker)
		if err != nil {
			return corporation, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return corporation, err
		}

		corporation.Id = id
	} else {
		_, err := db.db.Exec("UPDATE corporations SET corporation_id=?, name=?, ticker=? WHERE id=?", corporation.CorpId, corporation.Name, corporation.Ticker, corporation.Id)
		if err != nil {
			return corporation, err
		}
	}

	return corporation, nil
}
