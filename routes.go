// routes
package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		Name:        "Index",
		Methods:     []string{"GET"},
		Pattern:     "/",
		HandlerFunc: IndexHandler,
	},
	Route{
		Name:        "TrustRequest",
		Methods:     []string{"GET"},
		Pattern:     "/trustrequest",
		HandlerFunc: TrustRequestHandler,
	},
	Route{
		Name:        "FleetList",
		Methods:     []string{"GET"},
		Pattern:     "/fleets",
		HandlerFunc: FleetListHandler,
	},
	Route{
		Name:        "FleetListAll",
		Methods:     []string{"GET"},
		Pattern:     "/fleets/all",
		HandlerFunc: FleetListAllHandler,
	},
	Route{
		Name:        "FleetCreate",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/fleets/create",
		HandlerFunc: FleetCreateHandler,
	},
	Route{
		Name:        "FleetDetails",
		Methods:     []string{"GET"},
		Pattern:     "/fleet/{fleetid:[0-9]+}",
		HandlerFunc: FleetDetailsHandler,
	},
	Route{
		Name:        "FleetEdit",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/edit",
		HandlerFunc: FleetEditHandler,
	},
	Route{
		Name:        "FleetDelete",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/delete",
		HandlerFunc: FleetDeleteHandler,
	},
	Route{
		Name:        "PlayerList",
		Methods:     []string{"GET"},
		Pattern:     "/players",
		HandlerFunc: PlayerListHandler,
	},
	Route{
		Name:        "PlayerCreate",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/players/create",
		HandlerFunc: PlayerCreateHandler,
	},
	Route{
		Name:        "PlayerDetails",
		Methods:     []string{"GET"},
		Pattern:     "/player/{playerid:[0-9]+}",
		HandlerFunc: PlayerDetailsHandler,
	},
	Route{
		Name:        "PlayerEdit",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/player/{playerid:[0-9]+}/edit",
		HandlerFunc: PlayerEditHandler,
	},
	Route{
		Name:        "PlayerDelete",
		Methods:     []string{"GET", "POST"},
		Pattern:     "/player/{playerid:[0-9]+}/delete",
		HandlerFunc: PlayerDeleteHandler,
	},
	Route{
		Name:        "AdminMenu",
		Methods:     []string{"GET"},
		Pattern:     "/admin",
		HandlerFunc: AdminMenuHandler,
	},
}
