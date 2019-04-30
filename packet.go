package kfnetwork

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type Packet interface {
	GetHeader() PacketHeader
	GetPayload() ([]byte, error)
}

type PacketHeader struct {
	Sequence uint16
	Type     uint16
	Length   uint16 `json:"-"`
}

type VersionPacket struct {
	PacketHeader
	Version float32
}

type ExitPacket struct {
	PacketHeader
}

type ErrorPacket struct {
	PacketHeader
	Message string
}

type LoginPacket struct {
	PacketHeader
	Name string
	ID   string
}

type UpdateGameStatePacket struct {
	PacketHeader
}

func (p VersionPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p VersionPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p ExitPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p ExitPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p ErrorPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p ErrorPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p LoginPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p LoginPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p UpdateGameStatePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p UpdateGameStatePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p *PacketHeader) GetBytes() []byte {
	bytes := []byte{}

	//sequence := make([]byte, 2)
	packetType := make([]byte, 2)
	length := make([]byte, 2)

	//binary.LittleEndian.PutUint16(sequence, p.Sequence)
	binary.LittleEndian.PutUint16(packetType, p.Type)
	binary.LittleEndian.PutUint16(length, p.Length)

	//bytes = append(bytes, sequence...)
	bytes = append(bytes, packetType...)
	bytes = append(bytes, length...)

	return bytes
}

func ReadPacketHeader(client net.Conn) (PacketHeader, error) {
	header := PacketHeader{}
	//sequence := make([]byte, 2)
	packetType := make([]byte, 2)
	length := make([]byte, 2)

	//_, e := client.Read(sequence)

	//if e != nil {
	//	return header, e
	//}

	_, e := client.Read(packetType)

	if e != nil {
		return header, e
	}

	_, e = client.Read(length)

	if e != nil {
		return header, e
	}

	//header.Sequence = binary.LittleEndian.Uint16(sequence)
	header.Type = binary.LittleEndian.Uint16(packetType)
	header.Length = binary.LittleEndian.Uint16(length)

	return header, nil
}

func ReadPacket(client net.Conn) (Packet, error) {
	var packet Packet

	header, e := ReadPacketHeader(client)

	if e != nil {
		return packet, e
	}

	jsonBytes := make([]byte, header.Length)

	size, e := client.Read(jsonBytes)

	if size != len(jsonBytes) {
		e := errors.New("packet payload size does not match read bytes")
		return packet, e
	}

	if e != nil {
		return packet, e
	}

	return ParsePacket(header, jsonBytes)
}

func WritePacket(client net.Conn, packet Packet) error {
	payload := []byte{}
	header := packet.GetHeader()
	jsonPayload, e := packet.GetPayload()

	if e != nil {
		return e
	}

	header.Length = uint16(len(jsonPayload))
	actualPacket, e := ParsePacket(header, jsonPayload)

	if e != nil {
		return e
	}

	bytes, e := actualPacket.GetPayload()

	if e != nil {
		return e
	}

	payload = append(payload, header.GetBytes()...)
	payload = append(payload, bytes...)

	client.Write(payload)
	return nil
}

func ParsePacket(header PacketHeader, payload []byte) (Packet, error) {
	packet, e := RenderPacket(header, payload)

	if e != nil {
		return packet, e
	}

	return packet, nil
}

func RenderPacket(header PacketHeader, payload []byte) (Packet, error) {
	var packet Packet

	switch header.Type {
	case PacketTypeExit:
		packet := ExitPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeError:
		packet := ErrorPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeVersion:
		packet := VersionPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeUpdateGameState:
		packet := UpdateGameStatePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	default:
		return packet, errors.New("unknown packet type")
	}
}
