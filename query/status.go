package goquery

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"
)

const timeout = time.Duration(200 * time.Millisecond)

func GetStatus(ip string, port string) Data {
	conn, connErr := net.Dial("udp", ip + ":" + port)
	if connErr != nil {
		return Data{}
	}

	token, tokenErr := token(conn)
	if tokenErr != nil && token != 0 {
		conn.Close()
		return Data{}
	}

	resp, statsErr := stats(conn, token)
	conn.Close()
	if statsErr != nil {
		return Data{}
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

func response(payload string) Data {
	status := Data{}
	raw := strings.Split(payload, "\x00")

	for i := 0; i + 1 < len(raw); i += 2 {
		key := raw[i]
		value := raw[i + 1]
		if key == "" {
			status.Put("players", players(raw, i + 1))
			break
		} else {
			status.Put(parseKeyValue(key, value))
		}
	}

	return status
}

func players(raw []string, pos int) []string {
	var start = pos
	for ; pos < len(raw); pos++ {
		value := raw[pos]
		if value == "\x01player_" {
			pos += 1
			start = pos + 1
		} else if value == "" {
			break
		}
	}
	return raw[start:pos]
}

func parseKeyValue(key string, value string) (string, interface{}) {
	switch key {
	case "hostname":
		return "motd", value
	case "gametype":
		return "game_type", value
	case "game_id":
		return "game_id", value
	case "version":
		return "version", value
	case "plugins":
		return "plugins", value
	case "map":
		return "map", value
	case "numplayers":
		return "online", parseInt(value)
	case "maxplayers":
		return "max", parseInt(value)
	case "hostport":
		return "port", parseInt(value)
	case "hostip":
		return "ip", value
	}
	return "", nil
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
