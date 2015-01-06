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
		Name:        "Legal",
		Methods:     []string{"GET"},
		Pattern:     "/legal",
		HandlerFunc: LegalHandler,
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
		Name:        "FleetListGet",
		Methods:     []string{"GET"},
		Pattern:     "/fleets",
		HandlerFunc: FleetListGetHandler,
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
		Name:        "FleetGet",
		Methods:     []string{"GET"},
		Pattern:     "/fleet/{fleetid:[0-9]+}",
		HandlerFunc: FleetGetHandler,
	},
	Route{
		Name:        "FleetPut",
		Methods:     []string{"PUT"},
		Pattern:     "/fleet/{fleetid:[0-9]+}",
		HandlerFunc: FleetPutHandler,
	},
	Route{
		Name:        "FleetMembersGet",
		Methods:     []string{"GET"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/members",
		HandlerFunc: FleetMembersGetHandler,
	},
	Route{
		Name:        "FleetMembersPost",
		Methods:     []string{"POST"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/members",
		HandlerFunc: FleetMembersPostHandler,
	},
	Route{
		Name:        "FleetMembersPut",
		Methods:     []string{"PUT"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/members/{memberid:[0-9]+}",
		HandlerFunc: FleetMembersPutHandler,
	},
	Route{
		Name:        "FleetMembersDelete",
		Methods:     []string{"DELETE"},
		Pattern:     "/fleet/{fleetid:[0-9]+}/members/{memberid:[0-9]+}",
		HandlerFunc: FleetMembersDeleteHandler,
	},
	Route{
		Name:        "ReportListGet",
		Methods:     []string{"GET"},
		Pattern:     "/reports",
		HandlerFunc: ReportListGetHandler,
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
		Name:        "ReportGet",
		Methods:     []string{"GET"},
		Pattern:     "/report/{reportid:[0-9]+}",
		HandlerFunc: ReportGetHandler,
	},
	Route{
		Name:        "ReportPut",
		Methods:     []string{"PUT"},
		Pattern:     "/report/{reportid:[0-9]+}",
		HandlerFunc: ReportPutHandler,
	},
	Route{
		Name:        "ReportPlayersPut",
		Methods:     []string{"PUT"},
		Pattern:     "/report/{reportid:[0-9]+}/players",
		HandlerFunc: ReportPlayersPutHandler,
	},
}
