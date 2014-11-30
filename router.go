// router
package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
)

func SetupRouter(strictSlash bool) {
	r := mux.NewRouter().StrictSlash(strictSlash)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = WebLogger(handler, route.Name)

		r.Methods(route.Methods...).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/assets")))

	router = r
}
