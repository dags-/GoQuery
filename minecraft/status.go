package goquery

import (
	"io"
	"sort"
	"encoding/json"
)

type Status map[string]interface{}

func (status Status) SetValue(key string, value string) {
	switch key {
	case "numplayers":
		status["online"] = parseInt(value)
	case "maxplayers":
		status["max"] = parseInt(value)
	case "hostport":
		status["port"] = parseInt(value)
	case "hostip":
		status["ip"] = value
	default:
		status[key] = value
	}
}

func (status Status) SetPlayers(players []string) {
	sort.Strings(players)
	status["players"] = players
}

func (status *Status) ToJson(wr io.Writer, pretty bool) error {
	var prefix, indent = "", ""
	if pretty {
		indent = "    "
	}
	encoder := json.NewEncoder(wr)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(status)
	return err
}
