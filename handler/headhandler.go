package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dags-/goquery/query"
	"os"
	"time"
	"fmt"
)

var fetcher goquery.HeadFetcher

func NewHeadServer(scale int) func(w http.ResponseWriter, r *http.Request) {
	root, _ := os.Getwd()
	fetcher = goquery.NewHeadFetcher(root, "/heads", time.Duration(12 * time.Hour), ".png", scale)
	fmt.Println("Using PNG scale", scale)
	return head
}

func head(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	path := fetcher.Fetch(uuid)
	http.ServeFile(w, r, path)
}
