// handlers
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/morpheusxaut/lootsheeter/models"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Index"
	data["PageType"] = 1
	data["LoggedIn"] = session.IsLoggedIn(w, r)

	templates.ExecuteTemplate(w, "index", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)
	if loggedIn {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Login"
	data["PageType"] = 5
	data["LoggedIn"] = loggedIn

	state := GenerateRandomString(32)

	session.SetSSOState(w, r, state)

	data["SSOState"] = state

	templates.ExecuteTemplate(w, "login", data)
}

func LoginSSOHandler(w http.ResponseWriter, r *http.Request) {
	authorizationCode := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if len(authorizationCode) == 0 || len(state) == 0 {
		logger.Errorf("Received empty authorization code or state in LoginSSOHandler...")

		http.Error(w, "Received empty authorization code or state for SSO sign-on", http.StatusInternalServerError)
		return
	}

	savedState := session.GetSSOState(r)
	if !strings.EqualFold(savedState, state) {
		logger.Errorf("Failed to verify SSO state...")

		http.Redirect(w, r, "/login?error=sso_state", http.StatusSeeOther)
		return
	}

	t, err := FetchSSOToken(authorizationCode)
	if err != nil {
		logger.Errorf("Received error while fetching SSO token in LoginSSOHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v, err := FetchSSOVerification(t)
	if err != nil {
		logger.Errorf("Received error while fetching SSO verification in LoginSSOHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a, err := FetchCharacterAffiliation(v)
	if err != nil {
		logger.Errorf("Received error while fetching character association in LoginSSOHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sh, err := FetchCorporationSheet(a)
	if err != nil {
		logger.Errorf("Received error while fetching corporation sheet in LoginSSOHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.SetIdentity(w, r, a, sh)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FleetListHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Active Fleets"
	data["PageType"] = 2
	data["LoggedIn"] = loggedIn

	var fleets []*models.Fleet

	if HasAccessMask(r, int(models.AccessMaskAdmin)) {
		f, err := database.LoadAllFleets()
		if err != nil {
			logger.Errorf("Failed to load all fleets in FleetListHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fleets = f
	} else {
		corporationName := session.GetCorporationName(r)

		corporation, err := database.LoadCorporationFromName(corporationName)
		if err != nil {
			logger.Errorf("Failed to load corporation in FleetListHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := database.LoadAllFleetsFromCorpId(corporation.Id)
		if err != nil {
			logger.Errorf("Failed to load all fleets in FleetListHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fleets = f
	}

	data["Fleets"] = fleets
	data["ShowAll"] = false

	err := templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListHandler: [%v]", err)
	}
}

func FleetListAllHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "All Fleets"
	data["PageType"] = 2
	data["LoggedIn"] = loggedIn

	var fleets []*models.Fleet

	if HasAccessMask(r, int(models.AccessMaskAdmin)) {
		f, err := database.LoadAllFleets()
		if err != nil {
			logger.Errorf("Failed to load all fleets in FleetListAllHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fleets = f
	} else {
		corporationName := session.GetCorporationName(r)

		corporation, err := database.LoadCorporationFromName(corporationName)
		if err != nil {
			logger.Errorf("Failed to load corporation in FleetListAllHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := database.LoadAllFleetsFromCorpId(corporation.Id)
		if err != nil {
			logger.Errorf("Failed to load all fleets in FleetListAllHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fleets = f
	}

	data["Fleets"] = fleets
	data["ShowAll"] = true

	err := templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListAllHandler: [%v]", err)
	}
}

func FleetCreateHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Create Fleet"
	data["PageType"] = 2
	data["LoggedIn"] = loggedIn

	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Errorf("Failed to load all players in FleetCreateHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Players"] = players

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleetcreate", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetCreateHandler: [%v]", err)
	}
}

func FleetCreateFormHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse POST form in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleetCommanderId, err := strconv.ParseInt(r.FormValue("selectFleetCommander"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse commander ID in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fleetName := r.FormValue("textFleetName")
	fleetSystem := r.FormValue("textFleetSystem")
	fleetSystemNickname := r.FormValue("textFleetSystemNickname")

	if len(fleetName) == 0 || len(fleetSystem) == 0 {
		logger.Warnf("Content of POST form in FleetCreateFormHandler was empty...")

		http.Redirect(w, r, "/fleets/create", http.StatusSeeOther)
		return
	}

	corporationName := session.GetCorporationName(r)

	corporation, err := database.LoadCorporationFromName(corporationName)
	if err != nil {
		logger.Errorf("Failed to load corporation in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleet := models.NewFleet(-1, corporation.Id, fleetName, fleetSystem, fleetSystemNickname, 0, 0, 0, time.Now(), time.Time{}, 0, false, -1)

	player, err := database.LoadPlayer(fleetCommanderId)
	if err != nil {
		logger.Errorf("Failed to load fleet commander in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Printf("%#+v", fleet)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Printf("%#+v", fleet)

	commander := models.NewFleetMember(fleetCommanderId, fleet.Id, player, models.FleetRoleFleetCommander, 0, 0, 0, false, -1)

	fleet.AddMember(commander)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet commander in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Printf("%#+v", fleet)

	http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleet.Id), http.StatusSeeOther)
}

func FleetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	fleetId, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetDetailsHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = fmt.Sprintf("Details Fleet #%d", fleetId)
	data["PageType"] = 2
	data["LoggedIn"] = loggedIn

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load details for fleet #%d in FleetDetailsHandler: [%v]", fleetId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleet"] = fleet

	availablePlayers, err := database.LoadAvailablePlayers(fleetId, fleet.CorporationId)
	if err != nil {
		logger.Errorf("Failed to load available players in FleetDetailsHandler: [%v]", fleetId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["AvailablePlayers"] = availablePlayers

	logger.Printf("%#+v", availablePlayers)

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleetdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetDetailsHandler: [%v]", err)
	}
}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func FleetFinishHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	fleetId, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetFinishHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetFinishHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleet.FinishFleet()

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetFinishHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
}

func FleetAddProfitHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	fleetId, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetAddProfitHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !strings.Contains(strings.ToLower(r.Referer()), fmt.Sprintf("/fleet/%d", fleetId)) {
		logger.Warnf("Received request to FleetAddProfitHandler without proper referrer: %q", r.Referer())

		http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
		return
	}

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetAddProfitHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse POST form in FleetAddProfitHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawProfit := r.FormValue("addprofit_textarea")
	if len(rawProfit) == 0 {
		logger.Warnf("Content of POST form in FleetAddProfitHandler was empty...")

		http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
		return
	}

	var profit float64

	profit = 0

	if strings.Contains(strings.ToLower(rawProfit), "evepraisal") {
		rowSplit := strings.Split(rawProfit, "\r\n")

		for _, row := range rowSplit {
			p, err := GetEvepraisalValue(row)
			if err != nil {
				logger.Errorf("Failed to parse evepraisal row in FleetAddProfitHandler: [%v]", err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			profit += p
		}
	} else {
		p, err := GetPasteValue(rawProfit)
		if err != nil {
			logger.Errorf("Failed to parse paste in FleetAddProfitHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		profit = p
	}

	fleet.AddProfit(profit)

	_, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetAddProfitHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
	return
}

func FleetAddLossHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	fleetId, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetAddLossHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !strings.Contains(strings.ToLower(r.Referer()), fmt.Sprintf("/fleet/%d", fleetId)) {
		logger.Warnf("Received request to FleetAddLossHandler without proper referrer: %q", r.Referer())

		http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
		return
	}

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetAddLossHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse POST form in FleetAddLossHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawLoss := r.FormValue("addloss_textarea")
	if len(rawLoss) == 0 {
		logger.Warnf("Content of POST form in FleetAddLossHandler was empty...")

		http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
		return
	}

	var loss float64

	loss = 0

	if strings.Contains(strings.ToLower(rawLoss), "evepraisal") {
		rowSplit := strings.Split(rawLoss, "\r\n")

		for _, row := range rowSplit {
			l, err := GetEvepraisalValue(row)
			if err != nil {
				logger.Errorf("Failed to parse evepraisal row in FleetAddLossHandler: [%v]", err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			loss += l
		}
	} else if strings.Contains(strings.ToLower(rawLoss), "zkillboard") {
		rowSplit := strings.Split(rawLoss, "\r\n")

		for _, row := range rowSplit {
			l, err := GetzKillboardValue(row)
			if err != nil {
				logger.Errorf("Failed to parse zKillboard row in FleetAddLossHandler: [%v]", err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			loss += l
		}
	} else {
		l, err := GetPasteValue(rawLoss)
		if err != nil {
			logger.Errorf("Failed to parse paste in FleetAddLossHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		loss = l
	}

	fleet.AddLoss(loss)

	_, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetAddLossHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetId), http.StatusSeeOther)
}

func FleetDeleteHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func PlayerListHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Players"
	data["PageType"] = 3
	data["LoggedIn"] = loggedIn

	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Errorf("Failed to load all players in PlayerListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Players"] = players

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "players", data)
	if err != nil {
		logger.Errorf("Failed to execute template in PlayerListHandler: [%v]", err)
	}
}

func PlayerDetailsHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	playerId, err := strconv.ParseInt(vars["playerid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse player ID %q in PlayerDetailsHandler: [%v]", vars["playerid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = fmt.Sprintf("Details Player #%d", playerId)
	data["PageType"] = 3
	data["LoggedIn"] = loggedIn

	player, err := database.LoadPlayer(playerId)
	if err != nil {
		logger.Errorf("Failed to load details for player #%d in PlayerDetailsHandler: [%v]", playerId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Player"] = player

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "playerdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in PlayerDetailsHandler: [%v]", err)
	}
}

func PlayerEditHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func PlayerDeleteHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func ReportListHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Reports"
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn
	data["ShowAll"] = false

	reports, err := database.LoadAllReports()
	if err != nil {
		logger.Errorf("Failed to load all reports in ReportListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Reports"] = reports

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "reports", data)
	if err != nil {
		logger.Errorf("Failed to execute template in ReportListHandler: [%v]", err)
	}
}

func ReportListAllHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Reports"
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn
	data["ShowAll"] = true

	reports, err := database.LoadAllReports()
	if err != nil {
		logger.Errorf("Failed to load all reports in ReportListAllHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Reports"] = reports

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "reports", data)
	if err != nil {
		logger.Errorf("Failed to execute template in ReportListAllHandler: [%v]", err)
	}
}

func ReportDetailsHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	reportId, err := strconv.ParseInt(vars["reportid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse report ID %q in ReportDetailsHandler: [%v]", vars["reportid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = fmt.Sprintf("Report #%d", reportId)
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn

	report, err := database.LoadReport(reportId)
	if err != nil {
		logger.Errorf("Failed to load details for report #%d in ReportDetailsHandler: [%v]", reportId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	report.CalculatePayouts()

	data["Report"] = report

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "reportdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in ReportDetailsHandler: [%v]", err)
	}
}

func ReportCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func ReportCreateFormHandler(w http.ResponseWriter, r *http.Request) {

}

func AdminMenuHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Admin Menu"
	data["PageType"] = 6
	data["LoggedIn"] = loggedIn

	err := templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "adminmenu", data)
	if err != nil {
		logger.Errorf("Failed to execute template in AdminMenuHandler: [%v]", err)
	}
}
