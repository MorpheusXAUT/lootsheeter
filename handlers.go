// handlers
package main

import (
	"net"
	"net/http"
	"strconv"
)

func HandleRequests(host string, port int) {
	http.Handle("/", router)
	http.ListenAndServe(net.JoinHostPort(host, strconv.Itoa(port)), nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

}

func TrustRequestHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetListHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDetailsHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetEditHandler(w http.ResponseWriter, r *http.Request) {

}

func FleetDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerListHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDetailsHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerEditHandler(w http.ResponseWriter, r *http.Request) {

}

func PlayerDeleteHandler(w http.ResponseWriter, r *http.Request) {

}
