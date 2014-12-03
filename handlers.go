// handlers
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Index"
	data["PageType"] = 1

	templates.ExecuteTemplate(w, "index", data)
}

func TrustRequestHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Trust requested"
	data["PageType"] = 1

	templates.ExecuteTemplate(w, "index", data)
}

func FleetListHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Active Fleets"
	data["PageType"] = 2

	fleets, err := database.LoadAllFleets()
	if err != nil {
		logger.Errorf("Failed to load all fleets in FleetListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleets"] = fleets
	data["ShowAll"] = false

	err = templates.ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListHandler: [%v]", err)
	}
}

func FleetListAllHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "All Fleets"
	data["PageType"] = 2

	fleets, err := database.LoadAllFleets()
	if err != nil {
		logger.Errorf("Failed to load all fleets in FleetListAllHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleets"] = fleets
	data["ShowAll"] = true

	err = templates.ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListAllHandler: [%v]", err)
	}
}

func FleetCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDetailsHandler(w http.ResponseWriter, r *http.Request) {
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

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load details for fleet #%d in FleetDetailsHandler: [%v]", fleetId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Fleet"] = fleet

	err = templates.ExecuteTemplate(w, "fleetdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetDetailsHandler: [%v]", err)
	}
}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetFinishHandler(w http.ResponseWriter, r *http.Request) {
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

func FleetDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerListHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Players"
	data["PageType"] = 3

	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Errorf("Failed to load all players in PlayerListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Players"] = players

	err = templates.ExecuteTemplate(w, "players", data)
	if err != nil {
		logger.Errorf("Failed to execute template in PlayerListHandler: [%v]", err)
	}
}

func PlayerCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDetailsHandler(w http.ResponseWriter, r *http.Request) {
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

	player, err := database.LoadPlayer(playerId)
	if err != nil {
		logger.Errorf("Failed to load details for player #%d in PlayerDetailsHandler: [%v]", playerId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data["Player"] = player

	err = templates.ExecuteTemplate(w, "playerdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in PlayerDetailsHandler: [%v]", err)
	}
}

func PlayerEditHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func AdminMenuHandler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	data["PageTitle"] = "Admin Menu"
	data["PageType"] = 4

	err := templates.ExecuteTemplate(w, "adminmenu", data)
	if err != nil {
		logger.Errorf("Failed to execute template in AdminMenuHandler: [%v]", err)
	}
}
