package server

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
	Length   uint16
}

type CommonPacket struct {
	PacketHeader
}
type VersionPacket struct {
	CommonPacket
	Version float32
}

type ExitPacket struct {
	CommonPacket
}

type ErrorPacket struct {
	CommonPacket
	Message string
}

func (p VersionPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p VersionPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p *PacketHeader) GetBytes() []byte {
	bytes := []byte{}

	sequence := make([]byte, 2)
	packetType := make([]byte, 2)
	length := make([]byte, 2)

	binary.LittleEndian.PutUint16(sequence, p.Sequence)
	binary.LittleEndian.PutUint16(packetType, p.Type)
	binary.LittleEndian.PutUint16(length, p.Length)

	bytes = append(bytes, sequence...)
	bytes = append(bytes, packetType...)
	bytes = append(bytes, length...)

	return bytes
}

func ReadPacketHeader(client net.Conn) (PacketHeader, error) {
	header := PacketHeader{}
	sequence := make([]byte, 2)
	packetType := make([]byte, 2)
	length := make([]byte, 2)

	_, e := client.Read(sequence)

	if e != nil {
		return header, e
	}

	_, e = client.Read(packetType)

	if e != nil {
		return header, e
	}

	_, e = client.Read(length)

	if e != nil {
		return header, e
	}

	header.Sequence = binary.LittleEndian.Uint16(sequence)
	header.Type = binary.LittleEndian.Uint16(packetType)
	header.Length = binary.LittleEndian.Uint16(length)
	fmt.Println(header)
	return header, nil
}

func ReadPacket(client net.Conn) (Packet, error) {
	var packet Packet

	header, e := ReadPacketHeader(client)

	fmt.Println("Length:", header.Length)
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
	var packet Packet

	switch header.Type {
	case PacketTypeVersion:
		packet, e := RenderVersionPacket(payload)

		if e != nil {
			return packet, e
		}

		return packet, nil
	}

	return packet, nil
}

func RenderVersionPacket(payload []byte) (VersionPacket, error) {
	packet := VersionPacket{}

	e := json.Unmarshal(payload, &packet)

	if e != nil {
		return packet, e
	}

	return packet, nil
}

func RenderExitPacket(payload []byte) (ExitPacket, error) {
	packet := ExitPacket{}

	e := json.Unmarshal(payload, &packet)

	if e != nil {
		return packet, e
	}

	return packet, nil
}
