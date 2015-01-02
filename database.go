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
	fleetMembers map[int64]*models.FleetMember
	fleets       map[int64]*models.Fleet
	reports      map[int64]*models.Report
}

func NewDatabase(d *sql.DB) *Database {
	database := &Database{
		db:           d,
		corporations: make(map[int64]*models.Corporation),
		players:      make(map[int64]*models.Player),
		fleetMembers: make(map[int64]*models.FleetMember),
		fleets:       make(map[int64]*models.Fleet),
		reports:      make(map[int64]*models.Report),
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
		return &models.Corporation{}, err
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
		return &models.Corporation{}, err
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
		return corporations, err
	}

	for rows.Next() {
		var cid, corporationID, corporationAPIKeyID int64
		var corporationName, corporationTicker, corporationAPIKeyCode string
		var corporationCut float64

		err := rows.Scan(&cid, &corporationID, &corporationName, &corporationTicker, &corporationCut, &corporationAPIKeyID, &corporationAPIKeyCode)
		if err != nil {
			return corporations, err
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
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO corporations(corporation_id, name, ticker, api_keyid, api_keycode) VALUES (?, ?, ?, ?, ?)", corporation.CorporationID, corporation.Name, corporation.Ticker, corporation.APIID, corporation.APICode)
		if err != nil {
			return corporation, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return corporation, err
		}

		corporation.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE corporations SET corporation_id=?, name=?, ticker=?, api_keyid=?, api_keycode=? WHERE id=?", corporation.CorporationID, corporation.Name, corporation.Ticker, corporation.ID, corporation.APIID, corporation.APICode)
		if err != nil {
			return corporation, err
		}
	} else {
		return corporation, err
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
		return &models.Player{}, err
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
		return &models.Player{}, err
	}

	corp, err := db.LoadCorporation(cid)
	if err != nil {
		return &models.Player{}, err
	}

	player := models.NewPlayer(pid, playerID, playerName, corp, models.AccessMask(playerAccessMask))

	db.players[player.ID] = player

	return player, nil
}

func (db *Database) LoadAllPlayers(corporationID int64) ([]*models.Player, error) {
	logger.Tracef("Querying database for all players for corporation #%d...", corporationID)

	var players []*models.Player

	rows, err := db.db.Query("SELECT id, player_id, name, corporation_id, accessmask FROM players WHERE corporation_id = ? ORDER BY name", corporationID)
	if err != nil {
		return players, err
	}

	for rows.Next() {
		var pid, playerID, cid int64
		var playerAccessMask int
		var playerName string

		err := rows.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
		if err != nil {
			return players, err
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
		return players, err
	}

	for rows.Next() {
		var pid, playerID, cid int64
		var playerAccessMask int
		var playerName string

		err := rows.Scan(&pid, &playerID, &playerName, &cid, &playerAccessMask)
		if err != nil {
			return players, err
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
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO players(player_id, name, corporation_id, accessmask) VALUES (?, ?, ?, ?)", player.PlayerID, player.Name, player.Corp.ID, player.AccessMask)
		if err != nil {
			return player, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return player, err
		}

		player.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE players SET player_id=?, name=?, corporation_id=?, accessmask=? WHERE id=?", player.PlayerID, player.Name, player.Corp.ID, player.AccessMask, player.ID)
		if err != nil {
			return player, err
		}
	} else {
		return player, err
	}

	db.players[player.ID] = player

	return player, nil
}

func (db *Database) LoadFleetMember(fleetID int64, id int64) (*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet member with fid = %d and pid = %d...", fleetID, id)

	fleetMember, ok := db.fleetMembers[id]
	if ok {
		logger.Tracef("FleetMember with fid = %d and pid = %d found in cache, returning...", fleetID, id)
		return fleetMember, nil
	}

	row := db.db.QueryRow("SELECT id, fleet_id, player_id, role, ship, site_modifier, payment_modifier, payout, payout_complete, report_id FROM fleetmembers WHERE fleet_id = ? AND id = ?", fleetID, id)

	var fmid, fid, pid, rid int64
	var sqlRid sql.NullInt64
	var fleetmemberRole, fleetmemberSiteModifier int
	var fleetmemberPaymentModifier, fleetmemberPayout float64
	var fleetmemberPayoutCompleteEnum, fleetMemberShip string
	var fleetmemberPayoutComplete bool

	err := row.Scan(&fmid, &fid, &pid, &fleetmemberRole, &fleetMemberShip, &fleetmemberSiteModifier, &fleetmemberPaymentModifier, &fleetmemberPayout, &fleetmemberPayoutCompleteEnum, &sqlRid)
	if err != nil {
		return &models.FleetMember{}, err
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

	fleetMember = models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetMemberShip, fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid)

	db.fleetMembers[fleetMember.ID] = fleetMember
	if _, ok := db.fleets[fleetMember.FleetID]; ok {
		db.fleets[fleetMember.FleetID].UpdateMember(fleetMember)
	}

	return fleetMember, nil
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
			return fleetMembers, err
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

		fleetMember := models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetMemberShip, fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid)

		db.fleetMembers[fleetMember.ID] = fleetMember
		if _, ok := db.fleets[fleetMember.FleetID]; ok {
			db.fleets[fleetMember.FleetID].UpdateMember(fleetMember)
		}

		fleetMembers = append(fleetMembers, fleetMember)
	}

	return fleetMembers, nil
}

func (db *Database) LoadAllFleetMembersForReportPlayer(reportID int64, playerID int64) ([]*models.FleetMember, error) {
	logger.Tracef("Querying database for fleet members with rid = %d and pid = %d...", reportID, playerID)

	var fleetMembers []*models.FleetMember

	rows, err := db.db.Query("SELECT f.id, fleet_id, f.player_id, role, ship, site_modifier, payment_modifier, payout, payout_complete, report_id FROM fleetmembers AS f INNER JOIN players AS p ON f.player_id = p.id WHERE report_id = ? AND p.id = ? ORDER BY p.Name", reportID, playerID)
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
			return fleetMembers, err
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

		fleetMember := models.NewFleetMember(fmid, fid, player, models.FleetRole(fleetmemberRole), fleetMemberShip, fleetmemberSiteModifier, fleetmemberPaymentModifier, fleetmemberPayout, fleetmemberPayoutComplete, rid)

		db.fleetMembers[fleetMember.ID] = fleetMember
		if _, ok := db.fleets[fleetMember.FleetID]; ok {
			db.fleets[fleetMember.FleetID].UpdateMember(fleetMember)
		}

		fleetMembers = append(fleetMembers, fleetMember)
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
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO fleetmembers(fleet_id, player_id, role, ship, site_modifier, payment_modifier, payout, payout_complete, report_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", fleetID, member.Player.ID, member.Role, member.Ship, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum, fleetmemberReportID)
		if err != nil {
			return member, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return member, err
		}

		member.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE fleetmembers SET fleet_id=?, player_id=?, role=?, ship=?, site_modifier=?, payment_modifier=?, payout=?, payout_complete=?, report_id=? WHERE id=?", fleetID, member.Player.ID, member.Role, member.Ship, member.SiteModifier, member.PaymentModifier, member.Payout, fleetmemberPayoutCompleteEnum, fleetmemberReportID, member.ID)
		if err != nil {
			return member, err
		}
	} else {
		return member, err
	}

	db.fleetMembers[member.ID] = member
	if _, ok := db.fleets[member.FleetID]; ok {
		db.fleets[member.FleetID].UpdateMember(member)
	}

	return member, nil
}

func (db *Database) DeleteFleetMember(fleetID int64, memberID int64) error {
	logger.Tracef("Deleting member #%d from fleet #%d from database...")

	_, err := db.db.Exec("DELETE FROM fleetmembers WHERE fleet_id = ? AND id = ?", fleetID, memberID)
	if err != nil {
		return err
	}

	delete(db.fleetMembers, memberID)

	return nil
}

func (db *Database) LoadFleet(id int64) (*models.Fleet, error) {
	logger.Tracef("Querying database for fleet with fid = %d...", id)

	fleet, ok := db.fleets[id]
	if ok {
		logger.Tracef("Fleet with fid = %d found in cache, returning...", id)
		return fleet, nil
	}

	row := db.db.QueryRow("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, notes, report_id FROM fleets WHERE id = ?", id)

	var fid, cid, rid int64
	var sqlRid sql.NullInt64
	var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString, fleetNotes string
	var fleetProfit, fleetLosses, fleetCorporationPayout float64
	var fleetSitesFinished int
	var fleetStart, fleetEnd *time.Time
	var fleetPayoutComplete bool

	err := row.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &fleetNotes, &sqlRid)
	if err != nil {
		return &models.Fleet{}, err
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

	fleet = models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, fleetNotes, rid)

	for _, member := range fleetMembers {
		err = fleet.AddMember(member)
		if err != nil {
			return &models.Fleet{}, err
		}
	}

	db.fleets[fleet.ID] = fleet

	return fleet, nil
}

func (db *Database) LoadAllFleets(corporationID int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets for corporation #%d...", corporationID)

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, notes, report_id FROM fleets WHERE corporation_id = ?", corporationID)
	if err != nil {
		return fleets, err
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString, fleetNotes string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &fleetNotes, &sqlRid)
		if err != nil {
			return fleets, err
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

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, fleetNotes, rid)

		for _, member := range fleetMembers {
			err = fleet.AddMember(member)
			if err != nil {
				return fleets, err
			}
		}

		db.fleets[fleet.ID] = fleet

		fleets = append(fleets, fleet)
	}

	return fleets, nil
}

func (db *Database) LoadAllFleetsForReport(reportID int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets with rid = %d...", reportID)

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, notes, report_id FROM fleets WHERE report_id = ?", reportID)
	if err != nil {
		return fleets, err
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString, fleetNotes string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &fleetNotes, &sqlRid)
		if err != nil {
			return fleets, err
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

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, fleetNotes, rid)

		for _, member := range fleetMembers {
			err = fleet.AddMember(member)
			if err != nil {
				return fleets, err
			}
		}

		db.fleets[fleet.ID] = fleet

		fleets = append(fleets, fleet)
	}

	return fleets, nil
}

func (db *Database) LoadAllFleetsWithoutReports(corporationID int64) ([]*models.Fleet, error) {
	logger.Tracef("Querying database for all fleets for corporation #%d without reports...", corporationID)

	var fleets []*models.Fleet

	rows, err := db.db.Query("SELECT id, corporation_id, name, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, notes, report_id FROM fleets WHERE corporation_id = ? AND report_id IS NULL AND endtime IS NOT NULL", corporationID)
	if err != nil {
		return fleets, err
	}

	for rows.Next() {
		var fid, cid, rid int64
		var sqlRid sql.NullInt64
		var fleetName, fleetSystem, fleetSystemNickname, fleetPayoutCompleteEnumString, fleetNotes string
		var fleetProfit, fleetLosses, fleetCorporationPayout float64
		var fleetSitesFinished int
		var fleetStart, fleetEnd *time.Time
		var fleetPayoutComplete bool

		err := rows.Scan(&fid, &cid, &fleetName, &fleetSystem, &fleetSystemNickname, &fleetProfit, &fleetLosses, &fleetSitesFinished, &fleetStart, &fleetEnd, &fleetCorporationPayout, &fleetPayoutCompleteEnumString, &fleetNotes, &sqlRid)
		if err != nil {
			return fleets, err
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

		fleet := models.NewFleet(fid, corporation, fleetName, fleetSystem, fleetSystemNickname, fleetProfit, fleetLosses, fleetSitesFinished, *fleetStart, *fleetEnd, fleetCorporationPayout, fleetPayoutComplete, fleetNotes, rid)

		for _, member := range fleetMembers {
			err = fleet.AddMember(member)
			if err != nil {
				return fleets, err
			}
		}

		db.fleets[fleet.ID] = fleet

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
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO fleets(name, corporation_id, system, system_nickname, profit, losses, sites_finished, starttime, endtime, corporation_payout, payout_complete, notes, report_id) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", fleet.Name, fleet.Corporation.ID, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleetEndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString, fleet.Notes, fleetReportID)
		if err != nil {
			return fleet, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return fleet, err
		}

		fleet.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE fleets SET name=?, corporation_id=?, system=?, system_nickname=?, profit=?, losses=?, sites_finished=?, starttime=?, endtime=?, corporation_payout=?, payout_complete=?, notes=?, report_id=? WHERE id=?", fleet.Name, fleet.Corporation.ID, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleetEndTime, fleet.CorporationPayout, fleetPayoutCompleteEnumString, fleet.Notes, fleetReportID, fleet.ID)
		if err != nil {
			return fleet, err
		}
	} else {
		return fleet, err
	}

	for _, member := range fleet.Members {
		m, err := db.SaveFleetMember(fleet.ID, member)
		if err != nil {
			return fleet, err
		}

		member = m
	}

	db.fleets[fleet.ID] = fleet

	return fleet, nil
}

func (db *Database) LoadReportPayout(reportPayoutID int64) (*models.ReportPayout, error) {
	logger.Tracef("Querying database for report payout with rpid = %d...", reportPayoutID)

	row := db.db.QueryRow("SELECT id, report_id, player_id, payout, payout_complete FROM reportpayouts WHERE id = ?", reportPayoutID)

	var rpid, rid, pid int64
	var recordPayoutPayout float64
	var recordPayoutPayoutCompleteEnumString string
	var recordPayoutPayoutComplete bool

	err := row.Scan(&rpid, &rid, &pid, &recordPayoutPayout, &recordPayoutPayoutCompleteEnumString)
	if err != nil {
		return &models.ReportPayout{}, err
	}

	if strings.EqualFold(recordPayoutPayoutCompleteEnumString, "y") {
		recordPayoutPayoutComplete = true
	} else {
		recordPayoutPayoutComplete = false
	}

	player, err := database.LoadPlayer(pid)
	if err != nil {
		return &models.ReportPayout{}, err
	}

	reportPayout := models.NewReportPayout(rpid, rid, player, recordPayoutPayout, recordPayoutPayoutComplete)

	return reportPayout, nil
}

func (db *Database) LoadAllReportPayouts(reportID int64) ([]*models.ReportPayout, error) {
	logger.Tracef("Querying database for all report payouts with rid = %d...", reportID)

	var reportPayouts []*models.ReportPayout

	rows, err := db.db.Query("SELECT id, report_id, player_id, payout, payout_complete FROM reportpayouts WHERE report_id = ?", reportID)
	if err != nil {
		return reportPayouts, err
	}

	for rows.Next() {
		var rpid, rid, pid int64
		var recordPayoutPayout float64
		var recordPayoutPayoutCompleteEnumString string
		var recordPayoutPayoutComplete bool

		err := rows.Scan(&rpid, &rid, &pid, &recordPayoutPayout, &recordPayoutPayoutCompleteEnumString)
		if err != nil {
			return reportPayouts, err
		}

		if strings.EqualFold(recordPayoutPayoutCompleteEnumString, "y") {
			recordPayoutPayoutComplete = true
		} else {
			recordPayoutPayoutComplete = false
		}

		player, err := database.LoadPlayer(pid)
		if err != nil {
			return reportPayouts, err
		}

		reportPayout := models.NewReportPayout(rpid, rid, player, recordPayoutPayout, recordPayoutPayoutComplete)

		reportPayouts = append(reportPayouts, reportPayout)
	}

	return reportPayouts, nil
}

func (db *Database) SaveReportPayout(reportPayout *models.ReportPayout) (*models.ReportPayout, error) {
	logger.Tracef("Saving report payout #%d to database...", reportPayout.ID)

	var recordPayoutCompleteEnumString string

	if reportPayout.PayoutComplete {
		recordPayoutCompleteEnumString = "Y"
	} else {
		recordPayoutCompleteEnumString = "N"
	}

	_, err := db.LoadReportPayout(reportPayout.ID)
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO reportpayouts(report_id, player_id, payout, payout_complete) VALUES (?, ?, ?, ?)", reportPayout.ReportID, reportPayout.Player.ID, reportPayout.Payout, recordPayoutCompleteEnumString)
		if err != nil {
			return reportPayout, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return reportPayout, err
		}

		reportPayout.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE reportpayouts SET payout = ?, payout_complete = ? WHERE id = ?", reportPayout.Payout, recordPayoutCompleteEnumString, reportPayout.ID)
		if err != nil {
			return reportPayout, err
		}
	} else {
		return reportPayout, err
	}

	return reportPayout, nil
}

func (db *Database) LoadReport(id int64) (*models.Report, error) {
	logger.Tracef("Querying database for report with rid = %d...", id)

	report, ok := db.reports[id]
	if ok {
		logger.Tracef("Report with rid = %d found in cache, returning...", id)
		return report, nil
	}

	row := db.db.QueryRow("SELECT id, corporation_id, creator, total_payout, starttime, endtime, payout_complete FROM reports WHERE id=?", id)

	var rid, cid, pid int64
	var recordTotalPayout float64
	var recordPayoutCompleteEnumString string
	var recordStartTime, recordEndTime time.Time
	var recordPayoutComplete bool

	err := row.Scan(&rid, &cid, &pid, &recordTotalPayout, &recordStartTime, &recordEndTime, &recordPayoutCompleteEnumString)
	if err != nil {
		return &models.Report{}, err
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

	corporation, err := database.LoadCorporation(cid)
	if err != nil {
		return &models.Report{}, err
	}

	player, err := database.LoadPlayer(pid)
	if err != nil {
		return &models.Report{}, err
	}

	reportPayouts, err := db.LoadAllReportPayouts(rid)
	if err != nil {
		return &models.Report{}, err
	}

	report = models.NewReport(rid, recordTotalPayout, recordStartTime, recordEndTime, recordPayoutComplete, corporation, player, fleets)

	for _, payout := range reportPayouts {
		report.Payouts[payout.Player.Name] = payout
	}

	db.reports[report.ID] = report

	return report, nil
}

func (db *Database) LoadAllReports(corporationID int64) ([]*models.Report, error) {
	logger.Tracef("Querying database for all reports for corporation #%d...", corporationID)

	var reports []*models.Report

	rows, err := db.db.Query("SELECT id, corporation_id, creator, total_payout, starttime, endtime, payout_complete FROM reports WHERE corporation_id = ?", corporationID)
	if err != nil {
		return reports, err
	}

	for rows.Next() {
		var rid, cid, pid int64
		var recordTotalPayout float64
		var recordPayoutCompleteEnumString string
		var recordStartTime, recordEndTime time.Time
		var recordPayoutComplete bool

		err := rows.Scan(&rid, &cid, &pid, &recordTotalPayout, &recordStartTime, &recordEndTime, &recordPayoutCompleteEnumString)
		if err != nil {
			return reports, err
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

		corporation, err := database.LoadCorporation(cid)
		if err != nil {
			return reports, err
		}

		player, err := database.LoadPlayer(pid)
		if err != nil {
			return reports, err
		}

		reportPayouts, err := db.LoadAllReportPayouts(rid)
		if err != nil {
			return reports, err
		}

		report := models.NewReport(rid, recordTotalPayout, recordStartTime, recordEndTime, recordPayoutComplete, corporation, player, fleets)

		for _, payout := range reportPayouts {
			report.Payouts[payout.Player.Name] = payout
		}

		db.reports[report.ID] = report

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
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO reports(corporation_id, creator, total_payout, starttime, endtime, payout_complete) VALUES (?, ?, ?, ?, ?, ?)", report.Corporation.ID, report.Creator.ID, report.TotalPayout, report.StartRange, report.EndRange, reportPayoutCompleteEnum)
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
	} else if err == nil {
		_, err := db.db.Exec("UPDATE reports SET corporation_id=?, creator=?, total_payout=?, starttime=?, endtime=?, payout_complete=? WHERE id = ?", report.Corporation.ID, report.Creator.ID, report.TotalPayout, report.StartRange, report.EndRange, reportPayoutCompleteEnum, report.ID)
		if err != nil {
			return report, err
		}

		for _, reportPayout := range report.Payouts {
			reportPayout, err = db.SaveReportPayout(reportPayout)
			if err != nil {
				return report, err
			}
		}

		if report.PayoutComplete {
			for _, reportPayout := range report.Payouts {
				fleetMembers, err := db.LoadAllFleetMembersForReportPlayer(reportPayout.ReportID, reportPayout.Player.ID)
				if err != nil {
					return report, err
				}

				for _, member := range fleetMembers {
					member.PayoutComplete = true

					_, err := db.SaveFleetMember(member.FleetID, member)
					if err != nil {
						return report, err
					}
				}
			}

		}
	} else {
		return report, err
	}

	db.reports[report.ID] = report

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

func (db *Database) LoadLootPaste(id int64) (*models.LootPaste, error) {
	logger.Tracef("Querying database for loot paste with id = %d...", id)

	row := db.db.QueryRow("SELECT id, fleet_id, pasted_by, raw_paste, value, paste_type FROM lootpastes WHERE id = ?", id)

	var lid, lootPasteFleetID, lootPastePastedBy int64
	var lootPasteRawPaste string
	var lootPasteValue float64
	var lootPastePasteType int

	err := row.Scan(&lid, &lootPasteFleetID, &lootPastePastedBy, &lootPasteRawPaste, &lootPasteValue, &lootPastePasteType)
	if err != nil {
		return &models.LootPaste{}, err
	}

	return models.NewLootPaste(lid, lootPasteFleetID, lootPastePastedBy, lootPasteRawPaste, lootPasteValue, models.LootPasteType(lootPastePasteType)), nil
}

func (db *Database) SaveLootPaste(paste *models.LootPaste) (*models.LootPaste, error) {
	logger.Tracef("Saving loot paste #%d to database...", paste.ID)

	_, err := db.LoadLootPaste(paste.ID)
	if err == sql.ErrNoRows {
		result, err := db.db.Exec("INSERT INTO lootpastes(fleet_id, pasted_by, raw_paste, value, paste_type) VALUES(?, ?, ?, ?, ?)", paste.FleetID, paste.PastedBy, paste.RawPaste, paste.Value, paste.PasteType)
		if err != nil {
			return paste, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return paste, err
		}

		paste.ID = id
	} else if err == nil {
		_, err := db.db.Exec("UPDATE lootpastes SET fleet_id=?, pasted_by=?, raw_paste=?, value=?, paste_type=? WHERE id=?", paste.FleetID, paste.PastedBy, paste.RawPaste, paste.Value, paste.PasteType, paste.ID)
		if err != nil {
			return paste, err
		}
	} else {
		return paste, nil
	}

	return paste, nil
}

func (db *Database) RemovePlayerFromCache(id int64) {
	_, ok := db.players[id]
	if ok {
		delete(db.players, id)
	}
}
