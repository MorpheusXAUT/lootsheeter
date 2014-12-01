// handlers
package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/morpheusxaut/lootsheeter/models"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
	}{}

	templates.ExecuteTemplate(w, "index", data)
}

func TrustRequestHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetListHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Fleets []*models.Fleet
	}{}

	fleets, err := database.LoadAllFleets()
	if err != nil {
		logger.Errorf("Failed to load all fleets in FleetListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Fleets = fleets

	err = templates.ExecuteTemplate(w, "fleets", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetListHandler: [%v]", err)
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

	data := struct {
		Fleet *models.Fleet
	}{}

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load details for fleet #%d in FleetDetailsHandler: [%v]", fleetId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Fleet = fleet

	err = templates.ExecuteTemplate(w, "fleetdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in FleetDetailsHandler: [%v]", err)
	}
}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerListHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Players []*models.Player
	}{}

	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Errorf("Failed to load all players in PlayerListHandler: [%v]", err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Players = players

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

	data := struct {
		Player *models.Player
	}{}

	player, err := database.LoadPlayer(playerId)
	if err != nil {
		logger.Errorf("Failed to load details for player #%d in PlayerDetailsHandler: [%v]", playerId, err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data.Player = player

	err = templates.ExecuteTemplate(w, "playerdetails", data)
	if err != nil {
		logger.Errorf("Failed to execute template in PlayerDetailsHandler: [%v]", err)
	}
}

func PlayerEditHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
