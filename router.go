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

func SetupRouter(strictSlash bool) {
	logger.Infof("Setting up new routers (StrictSlash: %v)...", strictSlash)

	router = mux.NewRouter().StrictSlash(strictSlash)

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
	logger.Infof("Listening for requests on %q...", net.JoinHostPort(config.HttpHost, strconv.Itoa(config.HttpPort)))

	http.Handle("/", router)
	err := http.ListenAndServe(net.JoinHostPort(config.HttpHost, strconv.Itoa(config.HttpPort)), nil)

	logger.Fatalf("Received error while listening for requests: [%v]", err)
}
