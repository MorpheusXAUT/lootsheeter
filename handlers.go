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
		http.Redirect(w, r, session.GetLoginRedirect(r), http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Login"
	data["PageType"] = 2
	data["LoggedIn"] = loggedIn

	state := GenerateRandomString(32)

	session.SetSSOState(w, r, state)

	data["SSOState"] = state
	data["SSOClientID"] = config.SSOClientID
	data["SSOCallbackURL"] = config.SSOCallbackURL

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

		http.Redirect(w, r, "/login?error=ssoState", http.StatusSeeOther)
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

	err = session.SetIdentity(w, r, a, sh)
	if err != nil {
		logger.Errorf("Received error while setting identity in LoginSSOHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, session.GetLoginRedirect(r), http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FleetListGetHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/fleets")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in FleetListGetHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	showAll := false
	if len(r.FormValue("showAll")) > 0 {
		showAll = true
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Active Fleets"
	data["PageType"] = 3
	data["LoggedIn"] = loggedIn
	data["ShowAll"] = showAll

	corporationID := session.GetCorpID(r)

	fleets, err := database.LoadAllFleets(corporationID)
	if err != nil {
		logger.Errorf("Failed to load all fleets in FleetListGetHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleets"] = fleets

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListGetHandler: [%v]", err)
	}
}

func FleetCreateHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/fleets/create")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Create Fleet"
	data["PageType"] = 3
	data["LoggedIn"] = loggedIn

	corporationID := session.GetCorpID(r)

	players, err := database.LoadAllPlayers(corporationID)
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
		session.SetLoginRedirect(w, r, "/fleets/create")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse POST form in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleetCommanderID, err := strconv.ParseInt(r.FormValue("selectFleetCommander"), 10, 64)
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

	fleet := models.NewFleet(-1, corporation, fleetName, fleetSystem, fleetSystemNickname, 0, 0, 0, time.Now(), time.Time{}, 0, false, "", -1)

	player, err := database.LoadPlayer(fleetCommanderID)
	if err != nil {
		logger.Errorf("Failed to load fleet commander in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commander := models.NewFleetMember(fleetCommanderID, fleet.ID, player, models.FleetRoleFleetCommander, "", 0, 1, 0, false, -1)

	fleet.AddMember(commander)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet commander in FleetCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleet.ID), http.StatusSeeOther)
}

func FleetGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetGetHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = fmt.Sprintf("Details Fleet #%d", fleetID)
	data["PageType"] = 3
	data["LoggedIn"] = loggedIn

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load details for fleet #%d in FleetGetHandler: [%v]", fleetID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	corporationID := session.GetCorpID(r)

	if fleet.Corporation.ID != corporationID {
		http.Redirect(w, r, "/fleets", http.StatusSeeOther)
		return
	}

	data["Fleet"] = fleet

	availablePlayers, err := database.LoadAvailablePlayers(fleetID, fleet.Corporation.ID)
	if err != nil {
		logger.Errorf("Failed to load available players for fleet #%d in FleetGetHandler: [%v]", fleetID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["AvailablePlayers"] = availablePlayers

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleetdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetGetHandler: [%v]", err)
	}
}

func FleetPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetPutHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in FleetPutHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := r.FormValue("command")
	if len(command) == 0 {
		logger.Errorf("Received empty command in FleetPutHandler...")

		http.Error(w, "Received empty command", http.StatusBadRequest)
		return
	}

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetPutHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	corporationID := session.GetCorpID(r)

	if fleet.Corporation.ID != corporationID {
		http.Redirect(w, r, "/fleets", http.StatusSeeOther)
		return
	}

	switch strings.ToLower(command) {
	case "editdetails":
		FleetPutEditDetailsHandler(w, r, fleet)
		break
	case "addprofit":
		FleetPutAddProfitHandler(w, r, fleet)
		break
	case "addloss":
		FleetPutAddLossHandler(w, r, fleet)
		break
	case "calculatepayouts":
		FleetPutCalculatePayoutsHandler(w, r, fleet)
		break
	case "finishfleet":
		FleetPutFinishFleetHandler(w, r, fleet)
		break
	default:
		response := make(map[string]interface{})
		response["result"] = "error"
		response["error"] = "Invalid command"

		SendJSONResponse(w, response)
	}
}

func FleetPutEditDetailsHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetPutEditDetailsHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", r.FormValue("fleetDetailsStartTimeEdit"))
	if err != nil {
		logger.Errorf("Failed to parse startTime in FleetPutEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	var endTime time.Time

	if len(r.FormValue("fleetDetailsEndTimeEdit")) > 0 &&
		!strings.EqualFold(r.FormValue("fleetDetailsEndTimeEdit"), "YYYY-MM-DD HH:MM:SS +0000 UTC") &&
		!strings.EqualFold(r.FormValue("fleetDetailsEndTimeEdit"), "---") {
		e, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", r.FormValue("fleetDetailsEndTimeEdit"))
		if err != nil {
			logger.Errorf("Failed to parse endTime in FleetPutEditDetailsHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		endTime = e
	}

	sitesFinished, err := strconv.ParseInt(r.FormValue("fleetDetailsSitesFinishedEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse sitesFinished in FleetPutEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	payoutComplete, err := strconv.ParseBool(r.FormValue("fleetDetailsPayoutCompleteEdit"))
	if err != nil {
		logger.Errorf("Failed to parse payoutComplete in FleetPutEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	notes := r.FormValue("fleetDetailsNotesEdit")

	fleet.StartTime = startTime
	fleet.EndTime = endTime
	fleet.SitesFinished = int(sitesFinished)
	fleet.PayoutComplete = payoutComplete
	fleet.Notes = notes

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetPutEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetPutAddProfitHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasFleetRole(r, fleet, 8) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetPutAddProfitHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	rawProfit := r.FormValue("addProfitRaw")
	if len(rawProfit) == 0 {
		logger.Errorf("Content of rawProfit in FleetPutAddProfitHandler was empty...")

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Content of rawProfit was empty")

		SendJSONResponse(w, response)
		return
	}

	player := session.GetPlayerFromRequest(r)
	if player == nil {
		logger.Errorf("Failed to get player from request in FleetPutAddProfitHandler...")

		response["result"] = "error"
		response["error"] = "Failed to load player, cannot submit loot paste"

		SendJSONResponse(w, response)
		return
	}

	lootPaste := models.NewLootPaste(-1, fleet.ID, player.ID, rawProfit, 0, models.LootPasteTypeProfit)

	lootPaste, err := database.SaveLootPaste(lootPaste)
	if err != nil {
		logger.Errorf("Failed to save loot paste in FleetPutAddProfitHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	var profit float64

	profit = 0

	if strings.Contains(strings.ToLower(rawProfit), "evepraisal") {
		rowSplit := strings.Split(rawProfit, "\r\n")

		for _, row := range rowSplit {
			p, err := GetEvepraisalValue(row)
			if err != nil {
				logger.Errorf("Failed to parse evepraisal row in FleetPutAddProfitHandler: [%v]", err)

				response["result"] = "error"
				response["error"] = err.Error()

				SendJSONResponse(w, response)
				return
			}

			profit += p
		}
	} else {
		p, err := GetPasteValue(rawProfit)
		if err != nil {
			logger.Errorf("Failed to parse paste in FleetPutAddProfitHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		profit = p
	}

	lootPaste.Value = profit

	lootPaste, err = database.SaveLootPaste(lootPaste)
	if err != nil {
		logger.Errorf("Failed to update loot paste in FleetPutAddProfitHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.AddProfit(profit)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetPutAddProfitHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetPutAddLossHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasFleetRole(r, fleet, 8) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetPutAddLossHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	rawLoss := r.FormValue("addLossRaw")
	if len(rawLoss) == 0 {
		logger.Errorf("Content of rawLoss in FleetPutAddLossHandler was empty...")

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Content of rawLoss was empty")

		SendJSONResponse(w, response)
		return
	}

	player := session.GetPlayerFromRequest(r)
	if player == nil {
		logger.Errorf("Failed to get player from request in FleetPutAddLossHandler...")

		response["result"] = "error"
		response["error"] = "Failed to load player, cannot submit loot paste"

		SendJSONResponse(w, response)
		return
	}

	lootPaste := models.NewLootPaste(-1, fleet.ID, player.ID, rawLoss, 0, models.LootPasteTypeLoss)

	lootPaste, err := database.SaveLootPaste(lootPaste)
	if err != nil {
		logger.Errorf("Failed to save loot paste in FleetPutAddLossHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	var loss float64

	loss = 0

	if strings.Contains(strings.ToLower(rawLoss), "evepraisal") {
		rowSplit := strings.Split(rawLoss, "\r\n")

		for _, row := range rowSplit {
			l, err := GetEvepraisalValue(row)
			if err != nil {
				logger.Errorf("Failed to parse evepraisal row in FleetPutAddLossHandler: [%v]", err)

				response["result"] = "error"
				response["error"] = err.Error()

				SendJSONResponse(w, response)
				return
			}

			loss += l
		}
	} else if strings.Contains(strings.ToLower(rawLoss), "zkillboard") {
		rowSplit := strings.Split(rawLoss, "\r\n")

		for _, row := range rowSplit {
			l, err := GetZKillboardValue(row)
			if err != nil {
				logger.Errorf("Failed to parse zKillboard row in FleetPutAddLossHandler: [%v]", err)

				response["result"] = "error"
				response["error"] = err.Error()

				SendJSONResponse(w, response)
				return
			}

			loss += l
		}
	} else {
		l, err := GetPasteValue(rawLoss)
		if err != nil {
			logger.Errorf("Failed to parse paste in FleetPutAddLossHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		loss = l
	}

	lootPaste.Value = loss

	lootPaste, err = database.SaveLootPaste(lootPaste)
	if err != nil {
		logger.Errorf("Failed to update loot paste in FleetPutAddLossHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.AddLoss(loss)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetPutAddLossHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetPutCalculatePayoutsHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetPutCalculatePayoutsHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleet.CalculatePayouts()

	fleet, err := database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetPutCalculatePayoutsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetPutFinishFleetHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetPutFinishFleetHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleet.FinishFleet()

	fleet, err := database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetPutFinishFleetHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetMembersGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetMembersGetHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse fleet ID"

		SendJSONResponse(w, response)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fleetMembers, err := database.LoadAllFleetMembers(fleetID)
	if err != nil {
		logger.Errorf("Failed to load all fleet members for fleet #%d in FleetMembersGetHandler: [%v]", fleetID, err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleetMembers"] = fleetMembers

	SendJSONResponse(w, response)
	return
}

func FleetMembersPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	var errors []string

	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetMembersPostHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse fleet ID"

		SendJSONResponse(w, response)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetMembersPostHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	corporationID := session.GetCorpID(r)

	if fleet.Corporation.ID != corporationID {
		http.Redirect(w, r, "/fleets", http.StatusSeeOther)
		return
	}

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetMembersPostHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleetCommanders := fleet.FleetCommanders()

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in FleetPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetComposition := r.FormValue("addMemberFleetComposition")
	if len(fleetComposition) > 0 {
		fleetCompositionRows := strings.Split(fleetComposition, "\r\n")

		members, err := ParseFleetCompositionRows(fleet.ID, fleetCompositionRows)
		if err != nil {
			logger.Errorf("Failed to parse fleet composition rows in FleetMembersPostHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		for _, member := range members {
			if member.Role == models.FleetRoleFleetCommander && len(fleetCommanders) > 0 {
				secondCommander := true

				for _, commander := range fleetCommanders {
					if strings.EqualFold(member.Name, commander.Name) {
						secondCommander = false
						break
					}
				}

				if secondCommander {
					logger.Errorf("Tried to add second fleet commander to fleet in FleetMembersPostHandler...")

					errors = append(errors, "Cannot add two fleet commanders to the same fleet!")
					continue
				}
			}

			fleet.AddMember(member)
		}
	} else {
		memberID, err := strconv.ParseInt(r.FormValue("addMemberSelectMember"), 10, 64)
		if err != nil {
			logger.Errorf("Failed to parse memberID in FleetMembersPostHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		fleetRole, err := strconv.ParseInt(r.FormValue("addMemberSelectRole"), 10, 64)
		if err != nil {
			logger.Errorf("Failed to parse fleetRole in FleetMembersPostHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		if models.FleetRole(fleetRole) == models.FleetRoleFleetCommander && len(fleetCommanders) > 0 {
			logger.Errorf("Tried to add second fleet commander to fleet in FleetMembersPostHandler...")

			response["result"] = "error"
			response["error"] = "Cannot add two fleet commanders to the same fleet!"

			SendJSONResponse(w, response)
			return
		}

		ship := r.FormValue("addMemberShip")

		player, err := database.LoadPlayer(memberID)
		if err != nil {
			logger.Errorf("Failed to load player in FleetMembersPostHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		fleetMember := models.NewFleetMember(-1, fleet.ID, player, models.FleetRole(fleetRole), ship, 0, 1, 0, false, -1)

		fleet.AddMember(fleetMember)
	}

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetMembersPostHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	if len(errors) > 0 {
		response["result"] = "error"
		response["error"] = strings.Join(errors, ";")

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetMembersPutHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetMembersPutHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse fleet ID"

		SendJSONResponse(w, response)
		return
	}

	memberID, err := strconv.ParseInt(vars["memberid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse member ID %q in FleetMembersPutHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse member ID"

		SendJSONResponse(w, response)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	corporationID := session.GetCorpID(r)

	if fleet.Corporation.ID != corporationID {
		http.Redirect(w, r, "/fleets", http.StatusSeeOther)
		return
	}

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetMembersPutHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetRole, err := strconv.ParseInt(r.FormValue("fleetMemberRoleEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleetRole in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	siteModifier, err := strconv.ParseInt(r.FormValue("fleetMemberSiteModiferEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse siteModifier in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	paymentModifier, err := strconv.ParseFloat(r.FormValue("fleetMemberPaymentModifierEdit"), 64)
	if err != nil {
		logger.Errorf("Failed to parse paymentModifier in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	payoutComplete, err := strconv.ParseBool(r.FormValue("fleetMemberPayoutCompleteEdit"))
	if err != nil {
		logger.Errorf("Failed to parse payoutComplete in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetMember, err := database.LoadFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to load fleet member in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetCommanders := fleet.FleetCommanders()

	if fleetMember.Role == models.FleetRoleFleetCommander && models.FleetRole(fleetRole) != models.FleetRoleFleetCommander && len(fleetCommanders) <= 1 {
		logger.Errorf("Tried to remove fleet commander without replacement in FleetMembersPutHandler...")

		response["result"] = "error"
		response["error"] = "Cannot remove the fleet commander without replacement from the member list!"

		SendJSONResponse(w, response)
		return
	}

	fleetMember.Role = models.FleetRole(fleetRole)
	fleetMember.SiteModifier = int(siteModifier)
	fleetMember.PaymentModifier = paymentModifier
	fleetMember.PayoutComplete = payoutComplete

	fleet.Members[fleetMember.Name] = fleetMember

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetMembersPutHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetMembersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetMembersDeleteHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse fleet ID"

		SendJSONResponse(w, response)
		return
	}

	memberID, err := strconv.ParseInt(vars["memberid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse member ID %q in FleetMembersDeleteHandler: [%v]", vars["fleetid"], err)

		response["result"] = "error"
		response["error"] = "Failed to parse member ID"

		SendJSONResponse(w, response)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetMembersDeleteHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	corporationID := session.GetCorpID(r)

	if fleet.Corporation.ID != corporationID {
		http.Redirect(w, r, "/fleets", http.StatusSeeOther)
		return
	}

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetMembersDeleteHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	member, err := database.LoadFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to load fleet member in FleetMembersDeleteHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.RemoveMember(member.Name)

	fleetCommanders := fleet.FleetCommanders()

	if len(fleetCommanders) == 0 {
		logger.Errorf("Tried to remove fleet commander in FleetMembersDeleteHandler...")

		response["result"] = "error"
		response["error"] = "Cannot remove the fleet commander from the member list!"

		SendJSONResponse(w, response)
		return
	}

	err = database.DeleteFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to remove fleet member in FleetMembersDeleteHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func ReportListHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/reports")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Reports"
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn
	data["ShowAll"] = false

	corporationID := session.GetCorpID(r)

	reports, err := database.LoadAllReports(corporationID)
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
		session.SetLoginRedirect(w, r, "/reports/all")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Reports"
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn
	data["ShowAll"] = true

	corporationID := session.GetCorpID(r)

	reports, err := database.LoadAllReports(corporationID)
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

func ReportCreateHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/reports/create")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Create Report"
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn

	corpID := session.GetCorpID(r)

	fleets, err := database.LoadAllFleetsWithoutReports(corpID)
	if err != nil {
		logger.Errorf("Failed to load all reports in ReportCreateHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleets"] = fleets

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "reportcreate", data)
	if err != nil {
		logger.Errorf("Failed to execute template in ReportCreateHandler: [%v]", err)
	}
}

func ReportCreateFormHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/reports/create")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse POST form in ReportCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fleetsInclude := []string(r.Form["fleetsInclude"])
	if len(fleetsInclude) == 0 {
		logger.Warnf("Content of POST form in ReportCreateFormHandler was empty...")

		http.Redirect(w, r, "/reports/create", http.StatusSeeOther)
		return
	}

	var fleets []*models.Fleet
	startTime := time.Now()
	endTime := time.Time{}

	for _, fleet := range fleetsInclude {
		fleetID, err := strconv.ParseInt(fleet, 10, 64)
		if err != nil {
			logger.Errorf("Failed to parse fleet ID in ReportCreateFormHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := database.LoadFleet(fleetID)
		if err != nil {
			logger.Errorf("Failed to load fleet in ReportCreateFormHandler: [%v]", err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if f.StartTime.Before(startTime) {
			startTime = f.StartTime
		}

		if f.EndTime.After(endTime) {
			endTime = f.EndTime
		}

		fleets = append(fleets, f)
	}

	corporationID := session.GetCorpID(r)

	corporation, err := database.LoadCorporation(corporationID)
	if err != nil {
		logger.Errorf("Failed to load corporation in ReportCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player := session.GetPlayerFromRequest(r)

	report := models.NewReport(-1, 0, startTime, endTime, false, corporation, player, fleets)

	report, err = database.SaveReport(report)
	if err != nil {
		logger.Errorf("Failed to save report in ReportCreateFormHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/report/%d", report.ID), http.StatusSeeOther)
}

func ReportDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID, err := strconv.ParseInt(vars["reportid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse report ID %q in ReportDetailsHandler: [%v]", vars["reportid"], err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/report/%d", reportID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = fmt.Sprintf("Report #%d", reportID)
	data["PageType"] = 4
	data["LoggedIn"] = loggedIn

	report, err := database.LoadReport(reportID)
	if err != nil {
		logger.Errorf("Failed to load details for report #%d in ReportDetailsHandler: [%v]", reportID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	corporationID := session.GetCorpID(r)

	if report.Corporation.ID != corporationID {
		http.Redirect(w, r, "/reports", http.StatusSeeOther)
		return
	}

	report.CalculatePayouts()

	report, err = database.SaveReport(report)
	if err != nil {
		logger.Errorf("Failed to save report #%d in ReportDetailsHandler: [%v]", reportID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Report"] = report

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "reportdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in ReportDetailsHandler: [%v]", err)
	}
}

func ReportEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reportID, err := strconv.ParseInt(vars["reportid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse report ID %q in ReportEditHandler: [%v]", vars["reportID"], err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/report/%d", reportID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !strings.Contains(strings.ToLower(r.Referer()), fmt.Sprintf("/report/%d", reportID)) {
		logger.Warnf("Received request to ReportEditHandler without proper referrer: %q", r.Referer())

		http.Redirect(w, r, fmt.Sprintf("/report/%d", reportID), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in ReportEditHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := r.FormValue("command")
	if len(command) == 0 {
		logger.Errorf("Received empty command int ReportEditHandler...")

		http.Error(w, "Received empty command", http.StatusBadRequest)
		return
	}

	report, err := database.LoadReport(reportID)
	if err != nil {
		logger.Errorf("Failed to load report in ReportEditHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	corporationID := session.GetCorpID(r)

	if report.Corporation.ID != corporationID {
		http.Redirect(w, r, "/reports", http.StatusSeeOther)
		return
	}

	report.CalculatePayouts()

	report, err = database.SaveReport(report)
	if err != nil {
		logger.Errorf("Failed to save report #%d in ReportEditHandler: [%v]", reportID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch strings.ToLower(command) {
	case "poll":
		ReportEditPollHandler(w, r, report)
		break
	case "editdetails":
		ReportEditDetailsHandler(w, r, report)
		break
	case "addfleet":
		ReportEditAddFleetHandler(w, r, report)
		break
	case "editfleet":
		ReportEditEditFleetHandler(w, r, report)
		break
	case "removefleet":
		ReportEditRemoveFleetHandler(w, r, report)
		break
	case "playerpaid":
		ReportEditPlayerPaidHandler(w, r, report)
		break
	case "finish":
		ReportEditFinishHandler(w, r, report)
		break
	default:
		response := make(map[string]interface{})
		response["result"] = "error"
		response["error"] = "Invalid command"

		SendJSONResponse(w, response)
	}
}

func ReportEditPollHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditDetailsHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditAddFleetHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditEditFleetHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditRemoveFleetHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditPlayerPaidHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	if !IsReportCreator(r, report) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to ReportEditPlayerPaidHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask"

		SendJSONResponse(w, response)
		return
	}

	playerName := r.FormValue("playerName")
	if len(playerName) == 0 {
		logger.Errorf("Content of playerName in ReportEditPlayerPaidHandler was empty...")

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Content of playerName was empty")

		SendJSONResponse(w, response)
		return
	}

	reportPayout, ok := report.Payouts[playerName]
	if !ok {
		logger.Errorf("Failed to find ReportPayout for player %q in ReportEditPlayerPaidHandler...", playerName)

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Failed to find report payout for player")

		SendJSONResponse(w, response)
		return
	}

	reportPayout.PayoutComplete = true

	for _, payout := range reportPayout.Payouts {
		payout.PayoutComplete = true
	}

	report.Payouts[playerName] = reportPayout

	report, err := database.SaveReport(report)
	if err != nil {
		logger.Errorf("Failed to save report in ReportEditPlayerPaidHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditFinishHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	if !IsReportCreator(r, report) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to ReportEditFinishHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask"

		SendJSONResponse(w, response)
		return
	}

	for _, reportPayout := range report.Payouts {
		reportPayout.PayoutComplete = true

		for _, payout := range reportPayout.Payouts {
			payout.PayoutComplete = true
		}
	}

	report.PayoutComplete = true

	report, err := database.SaveReport(report)
	if err != nil {
		logger.Errorf("Failed to save report in ReportEditFinishHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}
