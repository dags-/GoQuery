package goquery

import (
	"fmt"
	"net"
	"time"
	"strings"
	"strconv"
	"github.com/pkg/errors"
)

const timeout = time.Duration(1 * time.Second)

func GetStatus(ip string, port string) (Status, error) {
	var status Status
	var token int32

	conn, err := net.Dial("udp", fmt.Sprint(ip, ":", port))
	if err == nil {
		token, err = getToken(conn)
	}

	if token == 0 {
		err = errors.New("Handshake with server failed")
	} else {
		resp, err := getStats(conn, token)
		if err == nil {
			status = parseResponse(resp)
		}
	}

	if conn != nil {
		conn.Close()
	}

	return status, err
}

func getToken(conn net.Conn) (int32, error) {
	var token int32

	conn.SetWriteDeadline(time.Now().Add(timeout))

	_, err := conn.Write(NewHandshake().ToBytes())
	if err == nil {
		buff := make([]byte, 32)

		conn.SetReadDeadline(time.Now().Add(timeout))
		count, err := conn.Read(buff)

		if err == nil {
			length := 5
			for ; (length < count) && (buff[length] != 0); length++ {}
			token = parseInt(string(buff[5:length]))
		}
	}

	return token, err
}

func getStats(conn net.Conn, token int32) (string, error) {
	var stats string
	var err error

	conn.SetWriteDeadline(time.Now().Add(timeout))
	_, err = conn.Write(NewStatusQuery(token).ToBytes())

	if err == nil {
		buff := make([]byte, 8196)
		conn.SetReadDeadline(time.Now().Add(timeout))
		count, err := conn.Read(buff)
		if err == nil {
			stats = string(buff[16:count])
		}
	}

	return stats, err
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