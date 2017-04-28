package discord

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"time"
)

type Status map[string]interface{}

const api = "http://discordapp.com/api/guilds/%s/widget.json";
var empty = Status{}

func GetStatus(serverId string) (Status, error) {
	var status Status

	url := fmt.Sprintf(api, serverId)
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := client.Get(url)

	if err == nil {
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&status)
	}

	if response != nil {
		response.Body.Close()
	}

	if message, ok := status["message"]; ok {
		err = errors.New(fmt.Sprint(message))
		status = empty
	}

	return status, err
}