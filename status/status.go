package status

import (
	"github.com/dags-/goquery/profile"
	"github.com/dags-/goquery/utils"
	"net"
	"strconv"
	"strings"
	"time"
)

type packetHead struct {
	Magic0    byte
	Magic1    byte
	Type      byte
	SessionId int32
}

type queryPacket struct {
	Head    packetHead
	Token   int32
	Padding int32
}

type ServerStatus struct {
	MOTD     string
	GameType string
	GameId   string
	Version  string
	Plugins  string
	Map      string
	Online   int32
	Max      int32
	IP       string
	Port     int32
	Players  profile.Profiles
}

const timeout = 200 * time.Millisecond

func QueryServer(ip string, port string) (ServerStatus, error) {
	conn, connErr := net.Dial("udp4", ip + ":" + port)
	if connErr != nil {
		return ServerStatus{}, connErr
	}

	token, tokenErr := getToken(conn)
	if tokenErr != nil {
		conn.Close()
		return ServerStatus{}, tokenErr
	}

	response, statsErr := getStats(conn, token)
	conn.Close()
	if statsErr != nil {
		return ServerStatus{}, statsErr
	}

	return parseResponse(response), nil
}

func getToken(conn net.Conn) (int32, error) {
	handshake := packetHead{0xFE, 0xFD, 0x09, 1}
	conn.SetWriteDeadline(time.Now().Add(timeout))
	_, writeError := conn.Write(handshake.bytes())
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
	token, err := strconv.Atoi(string(buff[5:length]))
	if err != nil {
		return 0, err
	}

	return int32(token), nil
}

func getStats(conn net.Conn, session int32) (string, error) {
	request := queryPacket{packetHead{0xFE, 0xFD, 0x00, 1}, session, 0}
	_, writeErr := conn.Write(request.bytes())
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

func parseResponse(payload string) ServerStatus {
	serverInfo := ServerStatus{}
	raw := strings.Split(payload, "\x00")

	for i := 0; i + 1 < len(raw); i += 2 {
		key := raw[i]
		value := raw[i + 1]
		if key == "" {
			serverInfo.Players = parsePlayerProfiles(raw, i + 1)
			break
		} else {
			serverInfo.setValue(key, value)
		}
	}

	return serverInfo
}

func parsePlayerProfiles(raw []string, pos int) profile.Profiles {
	players := make([]string, len(raw) - pos)

	var index = 0
	for ; pos + 1 < len(raw); pos++ {
		value := raw[pos]
		if value == "\x01player_" {
			// Full token is '\x00\x01player_\x00', but 'raw' was split on \x00 so
			// next string will be empty. Players array will follow, and ends on first empty element
			pos += 1
		} else if value == "" {
			break
		} else {
			players[index] = value
			index++
		}
	}

	return profile.LookupProfiles(players[0:index])
}

func (serverStatus *ServerStatus) setValue(key string, value string) {
	// Not sure if there's a less pants way of doing this
	switch key {
	case "hostname":
		serverStatus.MOTD = value
	case "gametype":
		serverStatus.GameType = value
	case "game_id":
		serverStatus.GameId = value
	case "version":
		serverStatus.Version = value
	case "plugins":
		serverStatus.Plugins = value
	case "map":
		serverStatus.Map = value
	case "numplayers":
		serverStatus.Online = queryutils.ParseInt(value)
	case "maxplayers":
		serverStatus.Max = queryutils.ParseInt(value)
	case "hostport":
		serverStatus.Port = queryutils.ParseInt(value)
	case "hostip":
		serverStatus.IP = value
	}
}

func (head *packetHead) bytes() []byte {
	return queryutils.ToBytes(head)
}

func (request *queryPacket) bytes() []byte {
	return queryutils.ToBytes(request)
}

func (serverStatus *ServerStatus) ToJson() string {
	return queryutils.ToJson(serverStatus, false)
}

func (serverStatus *ServerStatus) ToPrettyJson() string {
	return queryutils.ToJson(serverStatus, true)
}
