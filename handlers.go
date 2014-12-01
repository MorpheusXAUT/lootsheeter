// handlers
package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {

}

func TrustRequestHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetListHandler(w http.ResponseWriter, r *http.Request) {
	fleets, err := database.LoadAllFleets()
	if err != nil {
		logger.Errorf("Failed to load all fleets in FleetListHandler: [%v]", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "<html><body>")

	for _, fleet := range fleets {
		fmt.Fprintf(w, "ID: <a href=\"fleet/%d\">%d</a>, Name: %s, System: %s, Sites Finished: %d, Start: %v, End: %v<br />", fleet.Id, fleet.Id, fleet.Name, fleet.System, fleet.SitesFinished, fleet.StartTime, fleet.EndTime)
		for _, member := range fleet.Members {
			fmt.Fprintf(w, "Fleet #%d Member - ID: %d, Name: %s, Role: %s<br />", fleet.Id, member.Id, member.Player.Name, member.Role)
		}
	}

	fmt.Fprintln(w, "</body></html>")
}

func FleetCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fleetId, err := strconv.ParseInt(vars["fleetid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse fleet ID %q in FleetDetailsHandler: [%v]", vars["fleetid"], err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fleet, err := database.LoadFleet(fleetId)
	if err != nil {
		logger.Errorf("Failed to load details for fleet #%d in FleetDetailsHandler: [%v]", fleetId, err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "ID: %d, Name: %s, System: %s, System Nickname: %s, Profit: %f, Losses: %f, Sites Finished: %d, Start: %v, End: %v\n", fleet.Id, fleet.Name, fleet.System, fleet.SystemNickname, fleet.Profit, fleet.Losses, fleet.SitesFinished, fleet.StartTime, fleet.EndTime)
	for _, member := range fleet.Members {
		fmt.Fprintf(w, "Fleet #%d Member - ID: %d, Name: %s, Player ID: %d, Site Modifier: %d, Payment Modifier: %f, Role: %s\n", fleet.Id, member.Id, member.Player.Name, member.Player.PlayerId, member.SiteModifier, member.PaymentModifier, member.Role)
	}
}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerListHandler(w http.ResponseWriter, r *http.Request) {
	players, err := database.LoadAllPlayers()
	if err != nil {
		logger.Errorf("Failed to load all players in PlayerListHandler: [%v]", err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	for _, player := range players {
		fmt.Fprintf(w, "ID: %d, Name: %s, Corporation: %s\n", player.Id, player.Name, player.Corporation.Name)
	}
}

func PlayerCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerId, err := strconv.ParseInt(vars["playerid"], 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse player ID %q in PlayerDetailsHandler: [%v]", vars["playerid"], err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player, err := database.LoadPlayer(playerId)
	if err != nil {
		logger.Errorf("Failed to load details for player #%d in PlayerDetailsHandler: [%v]", playerId, err)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "ID: %d, Name: %s, Player ID: %d, Corporation: %s, Access: %s\n", player.Id, player.Name, player.PlayerId, player.Corporation.Name, player.AccessMask)
}

func PlayerEditHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
