package main

import (
	"fmt"
	"github.com/dags-/goquery/status"
)

func main() {
	serverStatus, err := status.QueryServer("mc.ardacraft.me", "25565")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(serverStatus.ToPrettyJson())
}
