package main

import (
	"net/url"
	"strings"
	"net/http"
	"fmt"
	"github.com/dags-/goquery"
	"os"
	"time"
)

func main() {
	StartServer();
}

func StartServer()  {
	root, _ := os.Getwd()
	fetcher := goquery.NewHeadFetcher(root, "../heads", time.Duration(12 * time.Hour), ".png")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ip, port := parseArgs(*request.URL)
		if ip != "" && port != "" && ip != "favicon.ico" {
			status := goquery.GetStatus(ip, port)
			profiles := goquery.Profiles(status.Players)

			go fetchHeads(fetcher, profiles)

			status.Players = profiles
			fmt.Fprintf(writer, goquery.ToJson(status, true))
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

func fetchHeads(fetcher goquery.HeadFetcher, profiles []goquery.Profile)  {
	sessions := goquery.Sessions(profiles...)
	fetcher.FetchHeads(sessions...)
}