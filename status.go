package goquery

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"
)

type Status struct {
	MOTD       string `json:"motd"`
	GameType   string `json:"game_type"`
	Game_Id    string `json:"game_id"`
	Version    string `json:"version"`
	Plugins    string `json:"plugins,omitempty"`
	Map        string `json:"map"`
	NumPlayers int32 `json:"online"`
	MaxPlayers int32 `json:"max"`
	IP         string `json:"ip"`
	Port       int32 `json:"port"`
	Players    interface{} `json:"players"`
}

const timeout = time.Duration(200 * time.Millisecond)

func GetStatus(ip string, port string) Status {
	conn, connErr := net.Dial("udp", ip + ":" + port)
	if connErr != nil {
		return Status{}
	}

	token, tokenErr := token(conn)
	if tokenErr != nil && token != 0 {
		conn.Close()
		return Status{}
	}

	resp, statsErr := stats(conn, token)
	conn.Close()
	if statsErr != nil {
		return Status{}
	}

	return response(resp)
}

func token(conn net.Conn) (int32, error) {
	conn.SetWriteDeadline(time.Now().Add(timeout))

	_, writeError := conn.Write(toBytes(handshake()))
	if writeError != nil {
		return 0, writeError
	}

	buff := make([]byte, 32)
	conn.SetReadDeadline(time.Now().Add(timeout))
	count, readErr := conn.Read(buff)
	if readErr != nil || count == 0 {
		return 0, readErr
	}

	length := 5
	for ; (length < count) && (buff[length] != 0); length++ {}

	return parseInt(string(buff[5:length])), nil
}

func stats(conn net.Conn, token int32) (string, error) {
	_, writeErr := conn.Write(toBytes(statusRequest(token)))
	if writeErr != nil {
		return "", writeErr
	}

	buff := make([]byte, 8192)
	conn.SetReadDeadline(time.Now().Add(timeout))
	count, readErr := conn.Read(buff)
	if readErr != nil || count == 0 {
		return "", readErr
	}

	return string(buff[16:count]), nil
}

func response(payload string) Status {
	status := Status{}
	raw := strings.Split(payload, "\x00")

	for i := 0; i + 1 < len(raw); i += 2 {
		key := raw[i]
		value := raw[i + 1]
		if key == "" {
			status.Players = players(raw, i + 1)
			break
		} else {
			status.setValue(key, value)
		}
	}

	return status
}

func (serverStatus *Status) setValue(key string, value string) {
	switch key {
	case "hostname":
		serverStatus.MOTD = value
	case "gametype":
		serverStatus.GameType = value
	case "game_id":
		serverStatus.Game_Id = value
	case "version":
		serverStatus.Version = value
	case "plugins":
		serverStatus.Plugins = value
	case "map":
		serverStatus.Map = value
	case "numplayers":
		serverStatus.NumPlayers = parseInt(value)
	case "maxplayers":
		serverStatus.MaxPlayers = parseInt(value)
	case "hostport":
		serverStatus.Port = parseInt(value)
	case "hostip":
		serverStatus.IP = value
	}
}

func players(raw []string, pos int) []string {
	var start = pos
	for ; pos < len(raw); pos++ {
		value := raw[pos]
		if value == "\x01player_" {
			pos += 2
			start = pos
		} else if value == "" {
			break
		}
	}
	return raw[start:pos]
}

func parseInt(value string) int32 {
	if num, err := strconv.Atoi(value); err == nil {
		return int32(num)
	}
	return 0
}

func toBytes(packet interface{}) []byte {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, packet)
	return buffer.Bytes()
}
