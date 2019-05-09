package kfnetwork

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
)

type Packet interface {
	GetHeader() PacketHeader
}

type PacketHeader struct {
	Sequence uint16 `json:"sequence"`
	Type     uint16 `json:"type"`
	Length   uint16 `json:"-"`
}

type VersionPacket struct {
	PacketHeader
	Version float32 `json:"version"`
}

type ExitPacket struct {
	PacketHeader
}

type ErrorPacket struct {
	PacketHeader
	Message string `json:"message"`
}

type LoginRequestPacket struct {
	PacketHeader
	Name  string `json:"name"`
	ID    string `json:"id"`
	Token string `json:"token"`
}

type LoginResponsePacket struct {
	PacketHeader
}

type PlayerListRequestPacket struct {
	PacketHeader
}

type PlayerListResponsePacket struct {
	PacketHeader
	PlayerList
}

type CreateLobbyRequestPacket struct {
	PacketHeader
	Name string `json:"name"`
}

type CreateLobbyResponsePacket struct {
	PacketHeader
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LobbyListRequestPacket struct {
	PacketHeader
}

type LobbyListResponsePacket struct {
	PacketHeader
	LobbyList
}

type JoinLobbyRequestPacket struct {
	PacketHeader
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JoinLobbyResponsePacket struct {
	PacketHeader
	ID      string `json:"id"`
	Name    string `json:"name"`
	Success bool   `json:"success"`
}

type LeaveLobbyRequestPacket struct {
	PacketHeader
}

type LeaveLobbyResponsePacket struct {
	PacketHeader
}

type LobbyChatRequestPacket struct {
	PacketHeader
	Message string `json:"message"`
}

type LobbyChatResponsePacket struct {
	PacketHeader
	Name    string `json:"name"`
	Message string `json:"message"`
}

type GlobalChatRequestPacket struct {
	PacketHeader
	Message string `json:"message"`
}

type GlobalChatResponsePacket struct {
	PacketHeader
	Name    string `json:"name"`
	Message string `json:"message"`
}
type LobbyBanRequestPacket struct {
	PacketHeader
	Target string `json:"target"`
}

type LobbyBanResponsePacket struct {
	PacketHeader
	Target  string `json:"target"`
	Success bool   `json:"success"`
}

type LobbyKickRequestPacket struct {
	PacketHeader
	Target string `json:"target"`
}

type LobbyKickResponsePacket struct {
	PacketHeader
	Target  string `json:"target"`
	Success bool   `json:"success"`
}

type UpdateGameStatePacket struct {
	PacketHeader
}

type CardPileRequestPacket struct {
	PacketHeader
	Pile uint8 `json:"pile"`
}

type CardPileResponsePacket struct {
	PacketHeader
	Cards []Card `json:"cards"`
}

type DrawCardRequestPacket struct {
	PacketHeader
}

type DrawCardResponsePacket struct {
	PacketHeader
	Card Card `json:"card"`
}

type PlayCardRequestPacket struct {
	PacketHeader
	Pile  uint8  `json:"pile"`
	ID    string `json:"id"`
	Index uint8  `json:"index"`
}

type PlayCardResponsePacket struct {
	PacketHeader
	Pile   uint8  `json:"pile"`
	ID     string `json:"id"`
	Index  uint8  `json:"index"`
	Played bool   `json:"played"`
}

type DiscardCardRequestPacket struct {
	PacketHeader
	Pile  uint8  `json:"pile"`
	ID    string `json:"id"`
	Index uint8  `json:"index"`
}

type DiscardCardResponsePacket struct {
	PacketHeader
	Pile   uint8  `json:"pile"`
	ID     string `json:"id"`
	Index  uint8  `json:"index"`
	Played bool   `json:"played"`
}

func (p PacketHeader) GetHeader() PacketHeader {
	return p
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
	jsonPayload, e := GetPacketPayload(packet)

	if e != nil {
		return e
	}

	header.Length = uint16(len(jsonPayload))

	payload = append(payload, header.GetBytes()...)
	payload = append(payload, jsonPayload...)

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
	case PacketTypeLoginRequest:
		packet := LoginRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeLoginResponse:
		packet := LoginResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeGlobalChatRequest:
		packet := GlobalChatRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeGlobalChatResponse:
		packet := GlobalChatResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypePlayerListRequest:
		packet := PlayerListRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypePlayerListResponse:
		packet := PlayerListResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeLobbyListRequest:
		packet := LobbyListRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeLobbyListResponse:
		packet := LobbyListResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeCreateLobbyRequest:
		packet := CreateLobbyRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeCreateLobbyResponse:
		packet := CreateLobbyResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeJoinLobbyRequest:
		packet := JoinLobbyRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeJoinLobbyResponse:
		packet := JoinLobbyResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeLeaveLobbyRequest:
		packet := LeaveLobbyRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeLeaveLobbyResponse:
		packet := LeaveLobbyResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeBanLobbyRequest:
		packet := LobbyBanRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeBanLobbyResponse:
		packet := LobbyBanResponsePacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeKickLobbyRequest:
		packet := LobbyKickRequestPacket{}
		e := json.Unmarshal(payload, &packet)
		return packet, e
	case PacketTypeKickLobbyResponse:
		packet := LobbyKickResponsePacket{}
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

func GetPacketPayload(packet Packet) ([]byte, error) {
	return json.Marshal(packet)
}
