package goquery

import (
	"bytes"
	"encoding/binary"
)

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

func NewHandshake() basePacket {
	return basePacket{0xFE, 0xFD, 0x09, 1}
}

func NewStatusQuery(token int32) queryPacket {
	return queryPacket{basePacket{0xFE, 0xFD, 0x00, 1}, token, 0}
}

func (packet basePacket) ToBytes() []byte {
	return packetToBytes(packet)
}

func (packet queryPacket) ToBytes() []byte {
	return packetToBytes(packet)
}

func packetToBytes(packet interface{}) []byte {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, packet)
	return buffer.Bytes()
}