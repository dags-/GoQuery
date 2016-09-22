package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dags-/goquery/query"
	"fmt"
	"strings"
	"time"
)

type QueryManager struct {
	statusCache  map[string]Status
	statusExpire time.Duration
}

type Status struct {
	status    goquery.Data
	timestamp time.Time
}

var manager = QueryManager{make(map[string]Status), time.Duration(15 * time.Second)}

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
	uuid, keys := r.FormValue("uuid"), r.FormValue("keys")

	status := getStatus(ip, port, uuid == "true")
	response = status

	if keys != "" {
		k := strings.Split(keys, ",")
		status = status.Retain(k...)
		response = status
	}

	fmt.Fprintf(w, goquery.ToJson(response, r.FormValue("pretty") == "true"))
}

func getStatus(ip string, port string, uuid bool) goquery.Data {
	key := fmt.Sprint(ip, ":", port, ":", uuid)
	status, ok := manager.statusCache[key]
	if !ok || time.Now().Sub(status.timestamp) > manager.statusExpire {
		status = Status{goquery.GetStatus(ip, port), time.Now()}
		if uuid {
			players, _ := status.status["players"]
			profiles := goquery.Profiles(players)
			status.status.Put("players", profiles)
		}
		manager.statusCache[key] = status
		fmt.Println("Refreshing status for", key)
	}
	return status.status
}