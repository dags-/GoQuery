package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"bufio"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dags-/goquery/discord"
	"github.com/dags-/goquery/minecraft"
)

type Response struct {
	Result string `json:"result"`
	Time   string `json:"time"`
	Data   interface{} `json:"data"`
}

func main() {
	go handleStop()

	var port string
	flag.StringVar(&port, "port", "8080", "Query port")
	flag.Parse()

	notFound, readErr := ioutil.ReadFile("notfound.html")

	router := mux.NewRouter()
	router.HandleFunc("/discord/{id}", discordHandler)
	router.HandleFunc("/minecraft/{ip}", minecraftHandler)
	router.HandleFunc("/minecraft/{ip}/{port}", minecraftHandler)

	if readErr == nil {
		router.NotFoundHandler = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
			wr.Write(notFound)
		})
	} else {
		fmt.Println(readErr)
	}

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

func minecraftHandler(wr http.ResponseWriter, rq *http.Request) {
	var status goquery.Status
	var err error

	vars := mux.Vars(rq)
	if ip, ok := vars["ip"]; ok {
		port := "25565"

		if _, ok := vars["port"]; ok {
			port = vars["port"]
		}

		status, err = goquery.GetStatus(ip, port)
	}

	response := wrapResponse(status, err)
	pretty := rq.FormValue("pretty") == "true"

	writeResponse(response, wr, pretty)
}

func discordHandler(wr http.ResponseWriter, rq *http.Request) {
	var data discord.Status
	var err error

	vars := mux.Vars(rq)
	if id, ok := vars["id"]; ok {
		data, err = discord.GetStatus(id)
	}

	response := wrapResponse(data, err)
	pretty := rq.FormValue("pretty") == "true"

	writeResponse(response, wr, pretty)
}

func wrapResponse(data interface{}, err error) Response {
	var result = fmt.Sprint(err)
	var timestamp = time.Now().Format(time.RFC822)

	if err == nil {
		result = "success"
	}

	return Response{Result: result, Time: timestamp, Data: data}
}

func writeResponse(resp Response, wr http.ResponseWriter, pretty bool) error {
	var prefix, indent = "", ""
	if pretty {
		indent = "  "
	}

	wr.WriteHeader(http.StatusOK)
	wr.Header().Set("Cache-Control", "max-age=60")
	wr.Header().Set("Content-Type", "application/json; charset=UTF-8")

	encoder := json.NewEncoder(wr)
	encoder.SetIndent(prefix, indent)

	return encoder.Encode(resp)
}