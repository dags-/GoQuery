package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/dags-/goquery/handler"
	"flag"
	"fmt"
	"bufio"
	"os"
)

func main() {
	go handleStop()

	var ipWhitelist string
	var port string
	var scale int

	flag.StringVar(&port, "port", "8080", "Query port")
	flag.IntVar(&scale, "scale", 8, "PNG scale")
	flag.StringVar(&ipWhitelist, "whitelist", "", "Whitelisted server IPs")
	flag.Parse()

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(pageNotFound)
	r.HandleFunc("/head/{uuid}", handler.NewHeadServer(scale))
	r.HandleFunc("/status/{ip}", handler.NewIPOnlyHandler(ipWhitelist))
	r.HandleFunc("/status/{ip}/{port}", handler.NewIpAndPortHandler(ipWhitelist))
	fmt.Println("Running on port", port)
	http.ListenAndServe(":" + port, handlers.CORS()(r))
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

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "not_fount.html")
}
