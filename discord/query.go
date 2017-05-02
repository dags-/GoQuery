package discord

import (
	"fmt"
	"bytes"
	"errors"
	"net/http"
	"encoding/json"
)

type Status map[string]interface{}

const api = "http://discordapp.com/api/guilds/%s/widget.json";
var empty = Status{}

func GetStatus(serverId string) (Status, error) {
	var status Status

	url := fmt.Sprintf(api, serverId)
	request, err := http.NewRequest("GET", url, &bytes.Buffer{})

	if err == nil {
		request.Header.Set("Connection", "close")

		response, err := http.DefaultClient.Do(request)
		if err == nil {
			decoder := json.NewDecoder(response.Body)
			err = decoder.Decode(&status)

			response.Close = true
			response.Body.Close()
		}

		if message, ok := status["message"]; ok {
			err = errors.New(fmt.Sprint(message))
			status = empty
		}

		request.Close = true
		request.Body.Close()
	}

	return status, err
}