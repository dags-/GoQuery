package main

import (
	"fmt"
	"github.com/dags-/goquery/status"
	"net/http"
	"strings"
	"time"
	"net/url"
)

var servers = make(map[string]status.CachedServerQuery)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ip, port := parseArgs(*request.URL)
		if ip != "" && port != "" {
			serverStatus, err := getStatus(ip, port)

			if err != nil {
				fmt.Fprintf(writer, "Error: %q", err.Error())
			} else {
				fmt.Fprintf(writer, serverStatus.ToPrettyJson())
			}
		}
	})
	http.ListenAndServe(":8080", nil)
}

func parseArgs(url url.URL) (string, string) {
	path := strings.Trim(url.Path, "/")
	args := strings.Split(path, "/")
	if len(args) > 0 {
		ip := args[0]
		port := "25565"
		if (len(args) > 1) {
			port = args[1]
		}
		return ip, port
	}
	return "", ""
}

func getStatus(ip string, port string) (status.ServerStatus, error) {
	key := ip + ":" + port
	serverQuery := servers[key]

	if !serverQuery.IsPresent() || !serverQuery.Matches(ip, port) {
		fmt.Println("Caching query for ", key)
		serverQuery = status.NewCachedServerQuery(ip, port).ExpireAfter(15 * time.Second)
	}

	serverStatus, err := serverQuery.Poll()
	servers[key] = serverQuery
	return serverStatus, err
}