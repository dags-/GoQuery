package goquery

type basePacket struct {
	Magic0    byte
	Magic1    byte
	Type      byte
	SessionId int32
}

type queryPacket struct {
	basePacket
	Token   int32
	Padding int32
}

func handshake() basePacket {
	return basePacket{0xFE, 0xFD, 0x09, 1}
}

func statusRequest(token int32) queryPacket {
	return queryPacket{basePacket{0xFE, 0xFD, 0x00, 1}, token, 0}
}
