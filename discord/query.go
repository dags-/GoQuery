package discord

import (
	"fmt"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type Status map[string]interface{}

const api = "http://discordapp.com/api/guilds/%s/widget.json";

func GetStatus(serverId string) (Status, error) {
	var status Status

	url := fmt.Sprintf(api, serverId)
	_, data, err := fasthttp.Get(nil, url)

	if err == nil {
		err = json.Unmarshal(data, &status)
	}

	return status, err
}