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
		Name:        "Login",
		Methods:     []string{"GET"},
		Pattern:     "/login",
		HandlerFunc: LoginHandler,
	},
	Route{
		Name:        "LoginSSO",
		Methods:     []string{"GET"},
		Pattern:     "/login/sso",
		HandlerFunc: LoginSSOHandler,
	},
	Route{
		Name:        "Logout",
		Methods:     []string{"GET"},
		Pattern:     "/logout",
		HandlerFunc: LogoutHandler,
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
		Methods:     []string{"GET"},
		Pattern:     "/fleets/create",
		HandlerFunc: FleetCreateHandler,
	},
	Route{
		Name:        "FleetCreateForm",
		Methods:     []string{"GET"},
		Pattern:     "/fleets/create",
		HandlerFunc: FleetCreateFormHandler,
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
		Name:        "FleetFinish",
		Methods:     []string{"GET"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/finish",
		HandlerFunc: FleetFinishHandler,
	},
	Route{
		Name:        "FleetAddProfit",
		Methods:     []string{"POST"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/addprofit",
		HandlerFunc: FleetAddProfitHandler,
	},
	Route{
		Name:        "FleetAddLoss",
		Methods:     []string{"POST"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/addloss",
		HandlerFunc: FleetAddLossHandler,
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
		Name:        "ReportList",
		Methods:     []string{"GET"},
		Pattern:     "/reports",
		HandlerFunc: ReportListHandler,
	},
	Route{
		Name:        "ReportListAll",
		Methods:     []string{"GET"},
		Pattern:     "/reports/all",
		HandlerFunc: ReportListAllHandler,
	},
	Route{
		Name:        "ReportCreate",
		Methods:     []string{"GET"},
		Pattern:     "/reports/create",
		HandlerFunc: ReportCreateHandler,
	},
	Route{
		Name:        "ReportCreateForm",
		Methods:     []string{"POST"},
		Pattern:     "/reports/create",
		HandlerFunc: ReportCreateFormHandler,
	},
	Route{
		Name:        "ReportDetails",
		Methods:     []string{"GET"},
		Pattern:     "/report/{reportid:[0-9]+}",
		HandlerFunc: ReportDetailsHandler,
	},
	Route{
		Name:        "AdminMenu",
		Methods:     []string{"GET"},
		Pattern:     "/admin",
		HandlerFunc: AdminMenuHandler,
	},
}
