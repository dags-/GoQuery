package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dags-/goquery/query"
)

func main() {
	go handleStop()

	var port string
	flag.StringVar(&port, "port", "8080", "Query port")
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/status/{ip}", serveStatus)
	router.HandleFunc("/status/{ip}/{port}", serveStatus)
	router.NotFoundHandler = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		http.ServeFile(wr, rq, "notfound.html")
	})

	fmt.Println("Launching on port", port)
	err := http.ListenAndServe(":" + port, handlers.CORS()(router))
	fmt.Println(err)
}

func handleStop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if text == "stop\n" {
			fmt.Println()
			fmt.Println("Stopping...")
			os.Exit(0)
			break
		}
	}
}

func serveStatus(wr http.ResponseWriter, rq *http.Request) {
	vars := mux.Vars(rq)
	if ip, ok := vars["ip"]; ok {
		port := "25565"
		pretty := rq.FormValue("pretty") == "true"
		if _, ok := vars["port"]; !ok {
			port = vars["port"]
		}
		wr.Header().Set("Content-Type", "application/json; charset=UTF-8")
		wr.WriteHeader(http.StatusOK)
		status := goquery.GetStatus(ip, port)
		status.ToJson(wr, pretty)
	}
}