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
	request.Header.Set("Connection", "close")

	resp, err := http.DefaultClient.Do(request)
	request.Close = true
	request.Body.Close()

	if err == nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&status)
		resp.Body.Close()
	}

	if message, ok := status["message"]; ok {
		err = errors.New(fmt.Sprint(message))
		status = empty
	}

	return status, err
}