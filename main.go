package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"bufio"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/dags-/goquery/query"
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

	router := mux.NewRouter()
	router.HandleFunc("/status/{ip}", serveMincraftStatus)
	router.HandleFunc("/status/{ip}/{port}", serveMincraftStatus)
	router.HandleFunc("/discord/{id}", serveDiscordStatus)
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

func serveMincraftStatus(wr http.ResponseWriter, rq *http.Request) {
	var status goquery.Status
	var err error
	var pretty bool

	vars := mux.Vars(rq)
	if ip, ok := vars["ip"]; ok {
		port := "25565"
		pretty = rq.FormValue("pretty") == "true"

		if _, ok := vars["port"]; ok {
			port = vars["port"]
		}

		status, err = goquery.GetStatus(ip, port)
	}

	response := wrapResponse(status, err)
	writeResponse(response, wr, pretty)
}

func serveDiscordStatus(wr http.ResponseWriter, rq *http.Request) {
	var data map[string]interface{}
	var err error
	var pretty bool

	vars := mux.Vars(rq)
	if id, ok := vars["id"]; ok {
		url := fmt.Sprintf("http://discordapp.com/api/guilds/%s/widget.json", id)
		pretty = rq.FormValue("pretty") == "true"

		var response *http.Response
		response, err = http.Get(url)

		defer response.Body.Close()

		if err == nil {
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&data)
		}
	}

	response := wrapResponse(data, err)
	writeResponse(response, wr, pretty)
}

func wrapResponse(data interface{}, err error) Response {
	var result = "fail"
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