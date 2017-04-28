package discord

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
)

type Status map[string]interface{}

const api = "http://discordapp.com/api/guilds/%s/widget.json";
var empty = Status{}

func GetStatus(serverId string) (Status, error) {
	var status Status

	url := fmt.Sprintf(api, serverId)
	response, err := http.Get(url)

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