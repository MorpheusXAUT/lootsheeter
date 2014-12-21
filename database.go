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
	db           *sql.DB
	corporations map[int64]*models.Corporation
	players      map[int64]*models.Player
}

func NewDatabase(d *sql.DB) *Database {
	database := &Database{
		db:           d,
		corporations: make(map[int64]*models.Corporation),
		players:      make(map[int64]*models.Player),
	}

	return database
}

func InitialiseDatabase() {
	logger.Infof("Trying to connect to MySQL database at %q...", net.JoinHostPort(config.MySqlHost, strconv.Itoa(config.MySqlPort)))

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", config.MySqlUser, config.MySqlPassword, net.JoinHostPort(config.MySqlHost, strconv.Itoa(config.MySqlPort)), config.MySqlDatabase))
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

	corp, ok := db.corporations[id]
	if ok {
		logger.Tracef("Corporation with cid = %d found in cache, returning...", id)
		return corp, nil
	}

	row := db.db.QueryRow("SELECT id, corporation_id, name, ticker, corporation_cut, api_keyid, api_keycode FROM corporations WHERE id = ?", id)

	var cid, corporationID, corporationAPIKeyID int64
	var corporationName, corporationTicker, corporationAPIKeyCode string
	var corporationCut float64

	err := row.Scan(&cid, &corporationID, &corporationName, &corporationTicker, &corporationCut, &corporationAPIKeyID, &corporationAPIKeyCode)
	if err != nil {
		return &models.Corporation{}, fmt.Errorf("Received error while scanning corporation row: [%v]", err)
	}

	corp = models.NewCorporation(cid, corporationID, corporationName, corporationTicker, corporationCut, corporationAPIKeyID, corporationAPIKeyCode)

	db.corporations[id] = corp

	return corp, nil
}

func (db *Database) LoadCorporationFromName(name string) (*models.Corporation, error) {
	logger.Tracef("Querying database for corporation with name = %q...", name)

	for _, corp := range db.corporations {
		if strings.EqualFold(name, corp.Name) {
			logger.Tracef("Corporation with name %q found in cache, returning...", name)
			return corp, nil
		}
	}

	row := db.db.QueryRow("SELECT id, corporation_id, name, ticker, corporation_cut, api_keyid, api_keycode FROM corporations WHERE name LIKE ?", name)

	var cid, corporationID, corporationAPIKeyID int64
	var corporationName, corporationTicker, corporationAPIKeyCode string
	var corporationCut float64

	err := row.Scan(&cid, &corporationID, &corporationName, &corporationTicker, &corporationCut, &corporationAPIKeyID, &corporationAPIKeyCode)
	if err != nil {
		return &models.Corporation{}, fmt.Errorf("Received error while scanning corporation name row: [%v]", err)
	}

	corp := models.NewCorporation(cid, corporationID, corporationName, corporationTicker, corporationCut, corporationAPIKeyID, corporationAPIKeyCode)

	db.corporations[corp.ID] = corp

	return corp, nil
}

func (db *Database) LoadAllCorporations() ([]*models.Corporation, error) {
	logger.Tracef("Querying database for all corporations...")

	var corporations []*models.Corporation

	rows, err := db.db.Query("SELECT id, corporation_id, name, ticker, corporation_cut, api_keyid, api_keycode FROM corporations ORDER BY name")
	if err != nil {
		return corporations, fmt.Errorf("Received error while querying for all corporations: [%v]", err)
	}

	for rows.Next() {
		var cid, corporationID, corporationAPIKeyID int64
		var corporationName, corporationTicker, corporationAPIKeyCode string
		var corporationCut float64

		err := rows.Scan(&cid, &corporationID, &corporationName, &corporationTicker, &corporationCut, &corporationAPIKeyID, &corporationAPIKeyCode)
		if err != nil {
			return corporations, fmt.Errorf("Received error while scanning corporations name row: [%v]", err)
		}

		corp := models.NewCorporation(cid, corporationID, corporationName, corporationTicker, corporationCut, corporationAPIKeyID, corporationAPIKeyCode)

		db.corporations[corp.ID] = corp

		corporations = append(corporations, corp)
	}

	return corporations, nil
}

func (db *Database) SaveCorporation(corporation *models.Corporation) (*models.Corporation, error) {
	logger.Tracef("Saving corporation #%d to database...", corporation.ID)

	_, err := db.LoadCorporation(corporation.ID)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO corporations(corporation_id, name, ticker, api_keyid, api_keycode) VALUES (?, ?, ?, ?, ?)", corporation.CorporationID, corporation.Name, corporation.Ticker, corporation.APIID, corporation.APICode)
		if err != nil {
			return corporation, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return corporation, err
		}

		corporation.ID = id
	} else {
		_, err := db.db.Exec("UPDATE corporations SET corporation_id=?, name=?, ticker=?, api_keyid=?, api_keycode=? WHERE id=?", corporation.CorporationID, corporation.Name, corporation.Ticker, corporation.ID, corporation.APIID, corporation.APICode)
		if err != nil {
			return corporation, err
		}
	}

	db.corporations[corporation.ID] = corporation

	return corporation, nil
}

func (db *Database) LoadPlayer(id int64) (*models.Player, error) {
	logger.Tracef("Querying database for player with pid = %d...", id)

	player, ok := db.players[id]
	if ok {
		logger.Tracef("Player with pid = %d found in cache, returning...", id)
		return player, nil
	}

	row := db.db.QueryRow("SELECT id, player_id, name, corporation_id, accessmask FROM players WHERE id = ?", id)

	var pid, playerID, cid int64
	var playerAccessMask int
	var playerName string

	err := row.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
	if err != nil {
		return &models.Player{}, fmt.Errorf("Received error while scanning player row: [%v]", err)
	}

	corp, err := db.LoadCorporation(cid)
	if err != nil {
		return &models.Player{}, err
	}

	player = models.NewPlayer(pid, playerID, playerName, corp, models.AccessMask(playerAccessMask))

	db.players[id] = player

	return player, nil
}

func (db *Database) LoadPlayerFromName(name string) (*models.Player, error) {
	logger.Tracef("Querying database for player with player_name = %q...", name)

	for _, player := range db.players {
		if strings.EqualFold(name, player.Name) {
			logger.Tracef("Player with name %q found in cache, returning...", name)
			return player, nil
		}
	}

	row := db.db.QueryRow("SELECT id, player_id, name, corporation_id, accessmask FROM players WHERE name LIKE ?", name)

	var pid, playerID, cid int64
	var playerAccessMask int
	var playerName string

	err := row.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
	if err != nil {
		return &models.Player{}, fmt.Errorf("Received error while scanning player name row: [%v]", err)
	}

	corp, err := db.LoadCorporation(cid)
	if err != nil {
		return &models.Player{}, err
	}

	player := models.NewPlayer(pid, playerID, playerName, corp, models.AccessMask(playerAccessMask))

	db.players[player.ID] = player

	return player, nil
}

func (db *Database) LoadAllPlayers() ([]*models.Player, error) {
	logger.Tracef("Querying database for all players...")

	var players []*models.Player

	rows, err := db.db.Query("SELECT id, player_id, name, corporation_id, accessmask FROM players ORDER BY name")
	if err != nil {
		return players, fmt.Errorf("Received error while querying for all players: [%v]", err)
	}

	for rows.Next() {
		var pid, playerID, cid int64
		var playerAccessMask int
		var playerName string

		err := rows.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
		if err != nil {
			return players, fmt.Errorf("Received error while scanning player rows: [%v]", err)
		}

		corp, err := db.LoadCorporation(cid)
		if err != nil {
			return players, err
		}

		player := models.NewPlayer(pid, playerID, playerName, corp, models.AccessMask(playerAccessMask))

		db.players[player.ID] = player

		players = append(players, player)
	}

	return players, nil
}

func (db *Database) LoadAvailablePlayers(fleedID int64, corporationID int64) ([]*models.Player, error) {
	logger.Tracef("Querying database for available players with cid = %d...", corporationID)

	var players []*models.Player

	rows, err := db.db.Query("SELECT id, player_id, name, corporation_id, accessmask FROM players WHERE corporation_id = ? AND id NOT IN (SELECT player_id FROM fleetmembers WHERE fleet_id = ?) ORDER BY name", corporationID, fleedID)
	if err != nil {
		return players, fmt.Errorf("Received error while querying for available players: [%v]", err)
	}

	for rows.Next() {
		var pid, playerID, cid int64
		var playerAccessMask int
		var playerName string

		err := rows.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
		if err != nil {
			return players, fmt.Errorf("Received error while scanning available player rows: [%v]", err)
		}

		corp, err := db.LoadCorporation(cid)
		if err != nil {
			return players, err
		}

		player := models.NewPlayer(pid, playerID, playerName, corp, models.AccessMask(playerAccessMask))

		db.players[player.ID] = player

		players = append(players, player)
	}

	return players, nil
}

func (db *Database) SavePlayer(player *models.Player) (*models.Player, error) {
	logger.Tracef("Saving player #%d to database...", player.ID)

	_, err := db.LoadPlayer(player.ID)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO players(player_id, name, corporation_id, accessmask) VALUES (?, ?, ?, ?)", player.PlayerID, player.Name, player.Corp.ID, player.AccessMask)
		if err != nil {
			return player, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return player, err
		}

		player.ID = id
	} else {
		_, err := db.db.Exec("UPDATE players SET player_id=?, name=?, corporation_id=?, accessmask=? WHERE id=?", player.PlayerID, player.Name, player.Corp.ID, player.AccessMask, player.ID)
		if err != nil {
			return player, err
		}
	}

	db.players[player.ID] = player

	return player, nil
}

func (db *Database) LoadFleetMember(fleetID int64, id int64) (*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet member with fid = %d and pid = %d...", fleetID, id)

	row := db.db.QueryRow("SELECT id, fleet_id, player_id, role, ship, site_modifier, payment_modifier, payout, payout_complete, report_id FROM fleetmembers WHERE fleet_id = ? AND id = ?", fleetID, id)

	var fmid, fid, pid, rid int64
	var sqlRid sql.NullInt64
	var fleetmemberRole, fleetmemberSiteModifier int
	var fleetmemberPaymentModifier, fleetmemberPayout float64
	var fleetmemberPayoutCompleteEnum, fleetMemberShip string
	var fleetmemberPayoutComplete bool

	err := row.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetMemberShip, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum, &sqlRid)
	if err != nil {
		return &models.FleetMember{}, fmt.Errorf("Received error while scanning fleet member row: [%v]", err)
	}

	if strings.EqualFold(fleetmemberPayoutCompleteEnum, "Y") {
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

	return models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetMemberShip, fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid), nil
}

func (db *Database) LoadAllFleetMembers(fleetID int64) ([]*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet members with fid = %d...", fleetID)

	var fleetMembers []*models.FleetMember

	rows, err := db.db.Query("SELECT f.id, fleet_id, f.player_id, role, ship, site_modifier, payment_modifier, payout, payout_complete, report_id FROM fleetmembers AS f INNER JOIN players AS p ON f.player_id = p.id WHERE fleet_id = ? ORDER BY p.Name", fleetID)
	if err != nil {
		return fleetMembers, err
	}

	for rows.Next() {
		var fmid, fid, pid, rid int64
		var sqlRid sql.NullInt64
		var fleetmemberRole, fleetmemberSiteModifier int
		var fleetmemberPaymentModifier, fleetmemberPayout float64
		var fleetmemberPayoutCompleteEnum, fleetMemberShip string
		var fleetmemberPayoutComplete bool

		err = rows.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetMemberShip, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum, &sqlRid)
		if err != nil {
			return fleetMembers, fmt.Errorf("Received error while scanning fleet member rows: [%v]", err)
		}

		if strings.EqualFold(fleetmemberPayoutCompleteEnum, "Y") {
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

		fleetMembers = append(fleetMembers, models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetMemberShip, fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid))
	}

	return fleetMembers, nil
}

func (db *Database) SaveFleetMember(fleetID int64, member *models.FleetMember) (*models.FleetMember, error) {
	logger.Tracef("Saving fleet member #%d to database...", member.ID)

	var fleetmemberReportID sql.NullInt64

	if member.ReportID > 0 {
		fleetmemberReportID.Int64 = member.ReportID
		fleetmemberReportID.Valid = true
	}

	var fleetmemberPayoutCompleteEnum string

	if member.PayoutComplete {
		fleetmemberPayoutCompleteEnum = "Y"
	} else {
		fleetmemberPayoutCompleteEnum = "N"
	}

	_, err := db.LoadFleetMember(fleetID, member.ID)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO fleetmembers(fleet_id, player_id, role, site_modifier, payment_modifier, payout, payout_complete, report_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", fleetID, member.Player.ID, member.Role, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum, fleetmemberReportID)
		if err != nil {
			return member, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return member, err
		}

		member.ID = id
	} else {
		_, err := db.db.Exec("UPDATE fleetmembers SET fleet_id=?, player_id=?, role=?, site_modifier=?, payment_modifier=?, payout=?, payout_complete=?, report_id=? WHERE id=?", fleetID, member.Player.ID, member.Role, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum, fleetmemberReportID, member.ID)
		if err != nil {
			return member, err
		}
	}

	return member, nil
}

func (db *Database) DeleteFleetMember(fleetID int64, memberID int64) error {
	logger.Tracef("Deleting member #%d from fleet #%d from database...")

	_, err := db.db.Exec("DELETE FROM fleetmembers WHERE fleet_id = ? AND id = ?", fleetID, memberID)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) LoadFleet(id int64) (*models.Fleet, error) {
	logger.Tracef("Querying database for fleet with fid = %d...", id)

	row := db.db.QueryRow("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id FROM fleets WHERE id = ?", id)

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

	corporation, err := db.LoadCorporation(cid)
	if err != nil {
		return &models.Fleet{}, err
	}

	fleetMembers, err := db.LoadAllFleetMembers(fid)
	if err != nil {
		return &models.Fleet{}, err
	}

	fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id FROM fleets")
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

		corporation, err := db.LoadCorporation(cid)
		if err != nil {
			return fleets, err
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadAllFleetsForCorporation(corporationID int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets with cid = %d...", corporationID)

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id FROM fleets WHERE corporation_id = ?", corporationID)
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

		corporation, err := db.LoadCorporation(cid)
		if err != nil {
			return fleets, err
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadAllFleetsForReport(reportID int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets with rid = %d...", reportID)

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id FROM fleets WHERE report_id = ?", reportID)
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

		corporation, err := db.LoadCorporation(cid)
		if err != nil {
			return fleets, err
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) LoadAllFleetsWithoutReports() ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets without reports...")

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id FROM fleets WHERE report_id IS NULL AND endtime IS NOT NULL")
	if err != nil {
		return fleets, fmt.Errorf("Received error while querying for all fleets without reports: [%v]", err)
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
			return fleets, fmt.Errorf("Received error while scanning fleet rows without reports: [%v]", err)
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

		corporation, err := db.LoadCorporation(cid)
		if err != nil {
			return fleets, err
		}

		fleetMembers, err := db.LoadAllFleetMembers(fid)
		if err != nil {
			return fleets, err
		}

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, rid)

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

func (db *Database) SaveFleet(fleet *models.Fleet) (*models.Fleet, error) {
	logger.Tracef("Saving fleet #%d to database...", fleet.ID)

	var fleetPayoutCompleteEnumString string

	if fleet.PayoutComplete {
		fleetPayoutCompleteEnumString = "Y"
	} else {
		fleetPayoutCompleteEnumString = "N"
	}

	var fleetReportID sql.NullInt64

	if fleet.ReportID > 0 {
		fleetReportID.Int64 = fleet.ReportID
		fleetReportID.Valid = true
	}

	var fleetEndTime *time.Time
	if !fleet.EndTime.IsZero() {
		fleetEndTime = &fleet.EndTime
	} else {
		fleetEndTime = nil
	}

	_, err := db.LoadFleet(fleet.ID)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO fleets(name, corporation_id, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, report_id) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", fleet.Name, fleet.Corporation.ID, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleetEndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString, fleetReportID)
		if err != nil {
			return fleet, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fleet, err
		}

		fleet.ID = id
	} else {
		_, err := db.db.Exec("UPDATE fleets SET name=?, corporation_id=?, system=?, system_nickname=?, profit=?, losses=?, sites_finished=?, starttime=?, endtime=?, corporation_payout=?, payout_complete=?, report_id=? WHERE id=?", fleet.Name, fleet.Corporation.ID, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleetEndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString, fleetReportID, fleet.ID)
		if err != nil {
			return fleet, err
		}
	}

	for _, member := range fleet.Members {
		m, err := db.SaveFleetMember(fleet.ID, member)
		if err != nil {
			return fleet, err
		}

		member = m
	}

	return fleet, nil
}

func (db *Database) LoadReport(id int64) (*models.Report, error) {
	logger.Tracef("Querying database for report with rid = %d...", id)

	row := db.db.QueryRow("SELECT id, creator, total_payout, starttime, endtime, payout_complete FROM reports WHERE id=?", id)

	var rid, pid int64
	var recordTotalPayout float64
	var recordPayoutCompleteEnumString string
	var recordStartTime, recordEndTime time.Time
	var recordPayoutComplete bool

	err := row.Scan(&rid, &pid, &recordTotalPayout, &recordStartTime, &recordEndTime, &recordPayoutCompleteEnumString)
	if err != nil {
		return &models.Report{}, fmt.Errorf("Received error while scanning report row: [%v]", err)
	}

	if strings.EqualFold(recordPayoutCompleteEnumString, "y") {
		recordPayoutComplete = true
	} else {
		recordPayoutComplete = false
	}

	fleets, err := database.LoadAllFleetsForReport(rid)
	if err != nil {
		return &models.Report{}, err
	}

	player, err := database.LoadPlayer(pid)
	if err != nil {
		return &models.Report{}, err
	}

	report := models.NewReport(rid, recordTotalPayout, recordStartTime, recordEndTime, recordPayoutComplete, player, fleets)

	return report, nil
}

func (db *Database) LoadAllReports() ([]*models.Report, error) {
	logger.Tracef("Querying database for all reports...")

	var reports []*models.Report

	rows, err := db.db.Query("SELECT id, creator, total_payout, starttime, endtime, payout_complete FROM reports")
	if err != nil {
		return reports, fmt.Errorf("Received error while querying for all reports: [%v]", err)
	}

	for rows.Next() {
		var rid, pid int64
		var recordTotalPayout float64
		var recordPayoutCompleteEnumString string
		var recordStartTime, recordEndTime time.Time
		var recordPayoutComplete bool

		err := rows.Scan(&rid, &pid, &recordTotalPayout, &recordStartTime, &recordEndTime, &recordPayoutCompleteEnumString)
		if err != nil {
			return reports, fmt.Errorf("Received error while scanning report rows: [%v]", err)
		}

		if strings.EqualFold(recordPayoutCompleteEnumString, "y") {
			recordPayoutComplete = true
		} else {
			recordPayoutComplete = false
		}

		fleets, err := database.LoadAllFleetsForReport(rid)
		if err != nil {
			return reports, err
		}

		player, err := database.LoadPlayer(pid)
		if err != nil {
			return reports, err
		}

		report := models.NewReport(rid, recordTotalPayout, recordStartTime, recordEndTime, recordPayoutComplete, player, fleets)

		reports = append(reports, report)
	}

	return reports, nil
}

func (db *Database) SaveReport(report *models.Report) (*models.Report, error) {
	logger.Tracef("Saving report #%d to database...", report.ID)

	var reportPayoutCompleteEnum string

	if report.PayoutComplete {
		reportPayoutCompleteEnum = "Y"
	} else {
		reportPayoutCompleteEnum = "N"
	}

	_, err := db.LoadReport(report.ID)
	if err != nil {
		result, err := db.db.Exec("INSERT INTO reports(creator, total_payout, starttime, endtime, payout_complete) VALUES (?, ?, ?, ?, ?)", report.Creator.ID, report.TotalPayout, report.StartRange, report.EndRange, reportPayoutCompleteEnum)
		if err != nil {
			return report, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return report, err
		}

		report.ID = id

		for _, fleet := range report.Fleets {
			fleet.ReportID = report.ID

			for _, member := range fleet.Members {
				member.ReportID = report.ID
			}

			f, err := database.SaveFleet(fleet)
			if err != nil {
				return report, err
			}

			fleet = f
		}
	} else {
		_, err := db.db.Exec("UPDATE reports SET creator=?, total_payout=?, starttime=?, endtime=?, payout_complete=? WHERE id = ?", report.Creator.ID, report.TotalPayout, report.StartRange, report.EndRange, reportPayoutCompleteEnum, report.ID)
		if err != nil {
			return report, err
		}

		for _, fleet := range report.Fleets {
			var payoutCompleteEnum string

			if fleet.PayoutComplete {
				payoutCompleteEnum = "Y"
			} else {
				payoutCompleteEnum = "N"
			}

			_, err := db.db.Exec("UPDATE fleets SET payout_complete = ? WHERE id = ? AND report_id = ?", payoutCompleteEnum, fleet.ID, report.ID)
			if err != nil {
				return report, err
			}
		}

		fleetPayoutsComplete := make(map[int64][]bool)

		for _, reportPayout := range report.Payouts {
			for _, payout := range reportPayout.Payouts {
				if _, ok := fleetPayoutsComplete[payout.FleetID]; !ok {
					fleetPayoutsComplete[payout.FleetID] = make([]bool, 0)
				}

				fleetPayoutsComplete[payout.FleetID] = append(fleetPayoutsComplete[payout.FleetID], payout.PayoutComplete)

				var payoutCompleteEnum string

				if payout.PayoutComplete {
					payoutCompleteEnum = "Y"
				} else {
					payoutCompleteEnum = "N"
				}

				_, err := db.db.Exec("UPDATE fleetmembers SET payout_complete = ? WHERE fleet_id = ? AND player_id = ? AND report_id = ?", payoutCompleteEnum, payout.FleetID, payout.PlayerID, report.ID)
				if err != nil {
					return report, err
				}
			}
		}

		for fleetID, payoutsComplete := range fleetPayoutsComplete {
			payoutComplete := true

			for _, complete := range payoutsComplete {
				if !complete {
					payoutComplete = false
				}
			}

			var payoutCompleteEnum string

			if payoutComplete {
				payoutCompleteEnum = "Y"
			} else {
				payoutCompleteEnum = "N"
			}

			_, err := db.db.Exec("UPDATE fleets SET payout_complete = ? WHERE id = ? AND report_id = ?", payoutCompleteEnum, fleetID, report.ID)
			if err != nil {
				return report, err
			}
		}
	}

	return report, nil
}

func (db *Database) QueryShipRole(ship string) (models.FleetRole, error) {
	logger.Tracef("Querying database for role for ship %q...", ship)

	row := db.db.QueryRow("SELECT fleet_role FROM fleetroles WHERE ship LIKE ?", strings.ToLower(ship))

	var fleetMemberRole int

	err := row.Scan(&fleetMemberRole)
	if err != nil {
		return models.FleetRoleUnknown, err
	}

	return models.FleetRole(fleetMemberRole), nil
}

func (db *Database) SaveLootPaste(fleetID int64, playerID int64, rawPaste string, value float64, pasteType string) error {
	logger.Tracef("Saving raw loot paste to database...")

	_, err := db.db.Exec("INSERT INTO lootpastes(fleet_id, pasted_by, raw_paste, value, paste_type) VALUES(?, ?, ?, ?, ?)", fleetID, playerID, rawPaste, value, pasteType)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemovePlayerFromCache(id int64) {
	_, ok := db.players[id]
	if ok {
		delete(db.players, id)
	}
}
