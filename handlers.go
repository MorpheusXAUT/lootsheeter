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

func FleetListHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, "/fleets")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Active Fleets"
	data["PageType"] = 3
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

		f, err := database.LoadAllFleetsForCorporation(corporation.ID)
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
		session.SetLoginRedirect(w, r, "/fleets/all")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "All Fleets"
	data["PageType"] = 3
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

		f, err := database.LoadAllFleetsForCorporation(corporation.ID)
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
		session.SetLoginRedirect(w, r, "/fleets/create")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	data["PageTitle"] = "Create Fleet"
	data["PageType"] = 3
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

	fleet := models.NewFleet(-1, corporation.ID, fleetName, fleetSystem, fleetSystemNickname, 0, 0, 0, time.Now(), time.Time{}, 0, false, -1)

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

func FleetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetDetailsHandler: [%v]", vars["fleetid"], err)

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
		logger.Errorf("Failed to load details for fleet #%d in FleetDetailsHandler: [%v]", fleetID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleet"] = fleet

	availablePlayers, err := database.LoadAvailablePlayers(fleetID, fleet.CorporationID)
	if err != nil {
		logger.Errorf("Failed to load available players for fleet #%d in FleetDetailsHandler: [%v]", fleetID, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["AvailablePlayers"] = availablePlayers

	err = templates.Funcs(TemplateFunctions(r)).ExecuteTemplate(w, "fleetdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetDetailsHandler: [%v]", err)
	}
}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetEditHandler: [%v]", vars["fleetid"], err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedIn := session.IsLoggedIn(w, r)

	if !loggedIn {
		session.SetLoginRedirect(w, r, fmt.Sprintf("/fleet/%d", fleetID))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !strings.Contains(strings.ToLower(r.Referer()), fmt.Sprintf("/fleet/%d", fleetID)) {
		logger.Warnf("Received request to FleetEditHandler without proper referrer: %q", r.Referer())

		http.Redirect(w, r, fmt.Sprintf("/fleet/%d", fleetID), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		logger.Errorf("Failed to parse form in FleetEditHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := r.FormValue("command")
	if len(command) == 0 {
		logger.Errorf("Received empty command int FleetEditHandler...")

		http.Error(w, "Received empty command", http.StatusBadRequest)
		return
	}

	fleet, err := database.LoadFleet(fleetID)
	if err != nil {
		logger.Errorf("Failed to load fleet in FleetEditHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch strings.ToLower(command) {
	case "poll":
		FleetEditPollHandler(w, r, fleet)
		break
	case "editdetails":
		FleetEditEditDetailsHandler(w, r, fleet)
		break
	case "addmember":
		FleetEditAddMemberHandler(w, r, fleet)
		break
	case "editmember":
		FleetEditEditMemberHandler(w, r, fleet)
		break
	case "removemember":
		FleetEditRemoveMemberHandler(w, r, fleet)
		break
	case "addprofit":
		FleetEditAddProfitHandler(w, r, fleet)
		break
	case "addloss":
		FleetEditAddLossHandler(w, r, fleet)
		break
	case "calculate":
		FleetEditCalculateHandler(w, r, fleet)
		break
	case "finish":
		FleetEditFinishHandler(w, r, fleet)
		break
	default:
		response := make(map[string]interface{})
		response["result"] = "error"
		response["error"] = "Invalid command"

		SendJSONResponse(w, response)
	}
}

func FleetEditPollHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["fleet"] = fleet

	SendJSONResponse(w, response)
}

func FleetEditEditDetailsHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditEditDetailsHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	startTime, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", r.FormValue("fleetDetailsStartTimeEdit"))
	if err != nil {
		logger.Errorf("Failed to parse startTime in FleetEditEditDetailsHandler: [%v]", err)

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
			logger.Errorf("Failed to parse endTime in FleetEditEditDetailsHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		endTime = e
	}

	sitesFinished, err := strconv.ParseInt(r.FormValue("fleetDetailsSitesFinishedEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse sitesFinished in FleetEditEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	payoutComplete, err := strconv.ParseBool(r.FormValue("fleetDetailsPayoutCompleteEdit"))
	if err != nil {
		logger.Errorf("Failed to parse payoutComplete in FleetEditEditDetailsHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.StartTime = startTime
	fleet.EndTime = endTime
	fleet.SitesFinished = int(sitesFinished)
	fleet.PayoutComplete = payoutComplete

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditEditDetailsHandler: [%v]", err)

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

func FleetEditAddMemberHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditAddMemberHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleetComposition := r.FormValue("addMemberFleetComposition")
	if len(fleetComposition) > 0 {
		fleetCompositionRows := strings.Split(fleetComposition, "\r\n")

		members, err := ParseFleetCompositionRows(fleet.ID, fleetCompositionRows)
		if err != nil {
			logger.Errorf("Failed to parse fleet composition rows in FleetEditAddMemberHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		for _, member := range members {
			fleet.AddMember(member)
		}
	} else {
		memberID, err := strconv.ParseInt(r.FormValue("addMemberSelectMember"), 10, 64)
		if err != nil {
			logger.Errorf("Failed to parse memberID in FleetEditAddMemberHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		fleetRole, err := strconv.ParseInt(r.FormValue("addMemberSelectRole"), 10, 64)
		if err != nil {
			logger.Errorf("Failed to parse fleetRole in FleetEditAddMemberHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		ship := r.FormValue("addMemberShip")

		player, err := database.LoadPlayer(memberID)
		if err != nil {
			logger.Errorf("Failed to load player in FleetEditAddMemberHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		fleetMember := models.NewFleetMember(-1, fleet.ID, player, models.FleetRole(fleetRole), ship, 0, 1, 0, false, -1)

		fleet.AddMember(fleetMember)
	}

	fleet, err := database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditAddMemberHandler: [%v]", err)

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

func FleetEditEditMemberHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditEditMemberHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	memberID, err := strconv.ParseInt(r.FormValue("fleetMemberMemberID"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse memberID in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetRole, err := strconv.ParseInt(r.FormValue("fleetMemberRoleEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleetRole in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	siteModifier, err := strconv.ParseInt(r.FormValue("fleetMemberSiteModiferEdit"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse siteModifier in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	paymentModifier, err := strconv.ParseFloat(r.FormValue("fleetMemberPaymentModifierEdit"), 64)
	if err != nil {
		logger.Errorf("Failed to parse paymentModifier in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	payoutComplete, err := strconv.ParseBool(r.FormValue("fleetMemberPayoutCompleteEdit"))
	if err != nil {
		logger.Errorf("Failed to parse payoutComplete in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleetMember, err := database.LoadFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to load fleet member in FleetEditEditMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

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
		logger.Errorf("Failed to save fleet in FleetEditEditMemberHandler: [%v]", err)

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

func FleetEditRemoveMemberHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditRemoveMemberHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	memberID, err := strconv.ParseInt(r.FormValue("removeMemberID"), 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse memberID in FleetEditRemoveMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	member, err := database.LoadFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to load fleet member in FleetEditRemoveMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.RemoveMember(member.Name)

	err = database.DeleteFleetMember(fleet.ID, memberID)
	if err != nil {
		logger.Errorf("Failed to remove fleet member in FleetEditRemoveMemberHandler: [%v]", err)

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

func FleetEditAddProfitHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasFleetRole(r, fleet, 8) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditAddProfitHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	rawProfit := r.FormValue("addProfitRaw")
	if len(rawProfit) == 0 {
		logger.Errorf("Content of rawProfit in FleetEditAddProfitHandler was empty...")

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Content of rawProfit was empty")

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
				logger.Errorf("Failed to parse evepraisal row in FleetEditAddMemberHandler: [%v]", err)

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
			logger.Errorf("Failed to parse paste in FleetEditAddMemberHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		profit = p
	}

	err := database.SaveLootPaste(fleet.ID, rawProfit, profit, "P")
	if err != nil {
		logger.Errorf("Failed to save loot paste in FleetEditAddMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.AddProfit(profit)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditAddMemberHandler: [%v]", err)

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

func FleetEditAddLossHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasFleetRole(r, fleet, 8) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditAddLossHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	rawLoss := r.FormValue("addLossRaw")
	if len(rawLoss) == 0 {
		logger.Errorf("Content of rawLoss in FleetEditAddLossHandler was empty...")

		response["result"] = "error"
		response["error"] = fmt.Sprintf("Content of rawLoss was empty")

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
				logger.Errorf("Failed to parse evepraisal row in FleetEditAddLossHandler: [%v]", err)

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
				logger.Errorf("Failed to parse zKillboard row in FleetEditAddLossHandler: [%v]", err)

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
			logger.Errorf("Failed to parse paste in FleetEditAddLossHandler: [%v]", err)

			response["result"] = "error"
			response["error"] = err.Error()

			SendJSONResponse(w, response)
			return
		}

		loss = l
	}

	err := database.SaveLootPaste(fleet.ID, rawLoss, loss, "L")
	if err != nil {
		logger.Errorf("Failed to save loot paste in FleetEditAddMemberHandler: [%v]", err)

		response["result"] = "error"
		response["error"] = err.Error()

		SendJSONResponse(w, response)
		return
	}

	fleet.AddLoss(loss)

	fleet, err = database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditAddLossHandler: [%v]", err)

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

func FleetEditCalculateHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditCalculateHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleet.CalculatePayouts()

	fleet, err := database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditCalculateHandler: [%v]", err)

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

func FleetEditFinishHandler(w http.ResponseWriter, r *http.Request, fleet *models.Fleet) {
	response := make(map[string]interface{})

	if !IsFleetCommander(r, fleet) && !HasAccessMask(r, int(models.AccessMaskAdmin)) {
		logger.Warnf("Received request to FleetEditFinishHandler without proper access...")

		response["result"] = "error"
		response["error"] = "Unauthorised access: cannot perform this operation with your current access mask or fleet role"

		SendJSONResponse(w, response)
		return
	}

	fleet.FinishFleet()

	fleet, err := database.SaveFleet(fleet)
	if err != nil {
		logger.Errorf("Failed to save fleet in FleetEditFinishHandler: [%v]", err)

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
		session.SetLoginRedirect(w, r, "/reports/all")
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

	fleets, err := database.LoadAllFleetsWithoutReports()
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

	player := session.GetPlayerFromRequest(r)

	report := models.NewReport(-1, 0, startTime, endTime, false, player, fleets)

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

	if !HasHigherAccessMask(r, int(models.AccessMaskJuniorFleetCommander)) {
		logger.Warnf("Received request to ReportEditHandler without proper access...")

		http.Redirect(w, r, fmt.Sprintf("/report/%d", reportID), http.StatusUnauthorized)
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
	case "calculate":
		ReportEditCalculateHandler(w, r, report)
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

func ReportEditCalculateHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}

func ReportEditFinishHandler(w http.ResponseWriter, r *http.Request, report *models.Report) {
	response := make(map[string]interface{})

	response["result"] = "success"
	response["error"] = nil
	response["report"] = report

	SendJSONResponse(w, response)
}
