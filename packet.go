package kfnetwork

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	keyforge "keyforge/game"
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

type LoginRequestPacket struct {
	PacketHeader
	Name string
	ID   string
}

type LoginResponsePacket struct {
	PacketHeader
}

type UpdateGameStatePacket struct {
	PacketHeader
}

type CardPileRequestPacket struct {
	PacketHeader
	Pile uint8
}

type CardPileResponsePacket struct {
	PacketHeader
	Cards []keyforge.Card
}

type DrawCardRequestPacket struct {
	PacketHeader
}

type DrawCardResponsePacket struct {
	PacketHeader
	Card keyforge.Card
}

type PlayCardRequestPacket struct {
	PacketHeader
	Pile  uint8
	ID    string
	Index uint8
}

type PlayCardResponsePacket struct {
	PacketHeader
	PlayCardRequestPacket
	Played bool
}

type DiscardCardRequestPacket struct {
	PacketHeader
	Pile  uint8
	ID    string
	Index uint8
}

type DiscardCardResponsePacket struct {
	PacketHeader
	DiscardCardRequestPacket
	Played bool
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

func (p LoginRequestPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p LoginRequestPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p LoginResponsePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p LoginResponsePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p UpdateGameStatePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p UpdateGameStatePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p CardPileRequestPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p CardPileRequestPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p DrawCardRequestPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p DrawCardRequestPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p DrawCardResponsePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p DrawCardResponsePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p PlayCardRequestPacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p PlayCardRequestPacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p PlayCardResponsePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p PlayCardResponsePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p DiscardCardResponsePacket) GetPayload() ([]byte, error) {
	return json.Marshal(p)
}

func (p DiscardCardResponsePacket) GetHeader() PacketHeader {
	return p.PacketHeader
}

func (p *PacketHeader) GetBytes() []byte {
	bytes := []byte{}

	packetType := make([]byte, 2)
	length := make([]byte, 2)

	binary.LittleEndian.PutUint16(packetType, p.Type)
	binary.LittleEndian.PutUint16(length, p.Length)

	bytes = append(bytes, packetType...)
	bytes = append(bytes, length...)

	return bytes
}

func ReadPacketHeader(client net.Conn) (PacketHeader, error) {
	header := PacketHeader{}
	packetType := make([]byte, 2)
	length := make([]byte, 2)

	_, e := client.Read(packetType)

	if e != nil {
		return header, e
	}

	_, e = client.Read(length)

	if e != nil {
		return header, e
	}

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

// RenderPacket - a giant case switch used to output the correct packet type
// when packets are read off of the wire.
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
	case PacketTypeVersionRequest:
		packet := VersionPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeVersionResponse:
		packet := VersionPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeUpdateGameState:
		packet := UpdateGameStatePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeCardPileRequest:
		packet := CardPileRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeCardPileResponse:
		packet := CardPileRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeDrawCardRequest:
		packet := DrawCardRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeDrawCardResponse:
		packet := DrawCardResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypePlayCardRequest:
		packet := PlayCardRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypePlayCardResponse:
		packet := PlayCardResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	default:
		return packet, errors.New("unknown packet type")
	}
}
