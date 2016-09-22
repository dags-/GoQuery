package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/dags-/goquery/handler"
	"flag"
	"fmt"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Query port")
	flag.Parse()
	fmt.Println("Running on port", port)

	r := mux.NewRouter()
	r.HandleFunc("/status/{ip}", handler.IpOnly)
	r.HandleFunc("/status/{ip}/{port}", handler.IpAndPort)
	r.HandleFunc("/head/{uuid}", handler.NewHeadServer())
	http.ListenAndServe(":" + port, handlers.CORS()(r))
}
