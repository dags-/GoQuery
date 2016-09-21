package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/dags-/goquery/query"
	"os"
	"time"
	"io/ioutil"
)

var fetcher goquery.HeadFetcher

func NewHeadServer() func(w http.ResponseWriter, r *http.Request) {
	root, _ := os.Getwd()
	fetcher = goquery.NewHeadFetcher(root, "../heads", time.Duration(12 * time.Hour), ".png", 64)
	return head
}

func head(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	path := fetcher.Fetch(uuid)
	data, _ := ioutil.ReadFile(path)
	w.Write(data)
}
