package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dags-/goquery/query"
	"fmt"
	"strings"
)

func IpOnly(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]
	sendStatus(w, r, ip, "25565")
}

func IpAndPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]
	port := vars["port"]
	sendStatus(w, r, ip, port)
}

func sendStatus(w http.ResponseWriter, r *http.Request, ip string, port string) {
	var response interface{}

	status := getStatus(ip, port)
	response = status

	if r.FormValue("uuid") == "true" {
		players, _ := status["players"]
		profiles := goquery.Profiles(players)
		status.Put("players", profiles)
		response = status
	}

	if val := r.FormValue("keys"); val != "" {
		keys := strings.Split(val, ",")
		status = status.Retain(keys...)
		response = status
	}

	fmt.Fprintf(w, goquery.ToJson(response, r.FormValue("pretty") == "true"))
}

func getStatus(ip string, port string) goquery.Data {
	status := goquery.GetStatus(ip, port)
	return status
}