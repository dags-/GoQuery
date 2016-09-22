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

	handler.SetWhitelist(ipWhitelist)

	r := mux.NewRouter()
	r.HandleFunc("/status/{ip}", handler.IpOnly)
	r.HandleFunc("/status/{ip}/{port}", handler.IpAndPort)
	r.HandleFunc("/head/{uuid}", handler.NewHeadServer(scale))

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
