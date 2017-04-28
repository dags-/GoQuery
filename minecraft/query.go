package goquery

import (
	"strconv"
	"net"
	"time"
	"strings"
	"fmt"
	"net/http"
)

func GetStatus(ip string, port string) (Status, error) {
	client := http.Transport{
		Dial: (&net.Dialer{Timeout: 5 * time.Second, }).Dial,
	}

	defer client.CloseIdleConnections()

	conn, connErr := client.Dial("udp", ip + ":" + port)

	if conn != nil {
		defer conn.Close()
	}

	if connErr != nil {
		fmt.Println(ip, port)
		fmt.Println(connErr)
		return Status{}, connErr
	}

	token, tokenErr := getToken(conn)
	if tokenErr != nil && token != 0 {
		return Status{}, tokenErr
	}

	resp, statsErr := getStats(conn, token)
	if statsErr != nil {
		return Status{}, statsErr
	}

	return parseResponse(resp), nil
}

func getToken(conn net.Conn) (int32, error) {
	conn.SetWriteDeadline(time.Now().Add(timeout))

	_, writeError := conn.Write(NewHandshake().ToBytes())
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

func getStats(conn net.Conn, token int32) (string, error) {
	_, writeErr := conn.Write(NewStatusQuery(token).ToBytes())
	if writeErr != nil {
		return "", writeErr
	}

	buff := make([]byte, 8196)
	conn.SetReadDeadline(time.Now().Add(timeout))
	count, readErr := conn.Read(buff)
	if readErr != nil || count == 0 {
		return "", readErr
	}

	return string(buff[16:count]), nil
}

func parseResponse(response string) Status {
	status := Status{}
	raw := strings.Split(response, "\x00")

	for i := 0; i + 1 < len(raw); i += 2 {
		key := raw[i]
		value := raw[i + 1]
		if key == "" {
			players := parsePlayers(raw, i + 1)
			status.SetPlayers(players)
			break
		} else {
			status.SetValue(key, value)
		}
	}

	return status
}

func parsePlayers(raw []string, pos int) []string {
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

func parseInt(value string) int32 {
	if num, err := strconv.Atoi(value); err == nil {
		return int32(num)
	}
	return 0
}