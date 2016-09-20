package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/dags-/goquery/test/handler"
	"github.com/gorilla/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/status/{ip}", handler.IpOnly)
	r.HandleFunc("/status/{ip}/{port}", handler.IpAndPort)
	r.HandleFunc("/head/{uuid}", handler.NewHeadServer())
	http.ListenAndServe(":8080", handlers.CORS()(r))
}