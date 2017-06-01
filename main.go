package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"github.com/dags-/goquery/server"
)

type Response struct {
	Result string `json:"result"`
	Time   string `json:"time"`
	Data   interface{} `json:"data"`
}

func main() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if scanner.Text() == "stop" {
				fmt.Println("Stopping...")
				os.Exit(0)
				break
			}
		}
	}()

	var port string
	flag.StringVar(&port, "port", "8080", "Query port")
	flag.Parse()
	fmt.Printf("Launcing on port %v\n", port)
	server.StartServer(port)
}