// router
package main

import (
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
)

func SetupRouter() {
	logger.Infof("Setting up new router...")

	router = mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = WebLogger(handler, route.Name)

		router.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/assets")))

	logger.Infof("Successfully set up new router!")
}

func HandleRequests() {
	logger.Infof("Listening for requests on %q...", net.JoinHostPort(config.HTTPHost, strconv.Itoa(config.HTTPPort)))

	http.Handle("/", router)
	err := http.ListenAndServe(net.JoinHostPort(config.HTTPHost, strconv.Itoa(config.HTTPPort)), nil)

	logger.Fatalf("Received error while listening for requests: [%v]", err)
}
