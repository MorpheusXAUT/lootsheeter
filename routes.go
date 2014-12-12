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
		Methods:     []string{"POST"},
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
		Name:        "ReportEdit",
		Methods:     []string{"GET"},
		Pattern:     "/report/{reportid:[0-9]+}/edit",
		HandlerFunc: ReportEditHandler,
	},
}
