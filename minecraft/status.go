package goquery

import (
	"time"
	"encoding/json"
	"io"
	"sort"
)

const timeout = time.Duration(1 * time.Second)

type Status struct {
	Ip       string `json:"ip"`
	Port     int32 `json:"port"`
	Motd     string `json:"motd"`
	GameType string `json:"gametype"`
	GameId   string `json:"gameid"`
	Version  string `json:"version"`
	Plugins  string `json:"plugins"`
	Map      string `json:"map"`
	Online   int32 `json:"online"`
	Max      int32 `json:"max"`
	Players  []string `json:"players"`
}

func (status *Status) SetValue(key string, value string) {
	switch key {
	case "hostname":
		status.Motd = value
	case "gametype":
		status.GameType = value
	case "game_id":
		status.GameId = value
	case "version":
		status.Version = value
	case "plugins":
		status.Plugins = value
	case "map":
		status.Map = value
	case "numplayers":
		status.Online = parseInt(value)
	case "maxplayers":
		status.Max = parseInt(value)
	case "hostport":
		status.Port = parseInt(value)
	case "hostip":
		status.Ip = value
	}
}

func (status *Status) SetPlayers(players []string) {
	sort.Strings(players)
	status.Players = players
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
