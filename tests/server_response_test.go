package tests

import (
	"testing"

	kf "github.com/team-neutron-shark/keyforge-network"
)

func TestServerResponseHandleVersionRequest(t *testing.T) {
	server := kf.NewServer(":4321")
	connection := NewMockNetworkConnection()

	versionRequestPacket := kf.VersionPacket{}
	versionRequestPacket.Type = kf.PacketTypeVersionRequest
	versionRequestPacket.Version = kf.ProtocolVersion

	server.HandleVersionRequest(connection, versionRequestPacket)

	versionResponsePacket, e := kf.ReadPacket(connection)

	if e != nil {
		t.Error(e.Error())
	}

	if versionResponsePacket.(kf.VersionPacket).Version != kf.ProtocolVersion {
		t.Error("protocol version mismatch")
	}

	server.Stop()
}
func TestServerResponseHandleGlobalChatRequest(t *testing.T) {
	server := kf.NewServer(":4321")
	player := kf.NewPlayer()
	connection := NewMockNetworkConnection()

	player.Name = "testing"
	player.ID = kf.GenerateUUID()
	player.Client = connection
	server.AddPlayer(player)

	request := kf.GlobalChatRequestPacket{}
	request.Type = kf.PacketTypeGlobalChatRequest
	request.Message = kf.GenerateUUID()

	server.HandleGlobalChatRequest(connection, request)

	response, e := kf.ReadPacket(connection)

	if e != nil {
		t.Error(e.Error())
	}

	if response.(kf.GlobalChatResponsePacket).Type != request.Type {
		t.Error("type mismatch")
	}

	if response.(kf.GlobalChatResponsePacket).Name != player.Name {
		t.Error("name mismatch")
	}

	if response.(kf.GlobalChatResponsePacket).Message != request.Message {
		t.Error("message mismatch")
	}

	server.Stop()
}

func TestServerResponseHandlePlayerListRequest(t *testing.T) {
	server := kf.NewServer(":4321")
	player := kf.NewPlayer()
	connection := NewMockNetworkConnection()

	player.Name = "testing"
	player.ID = kf.GenerateUUID()
	player.Client = connection
	server.AddPlayer(player)

	request := kf.PlayerListRequestPacket{}
	request.Type = kf.PacketTypePlayerListRequest

	server.HandlePlayerListRequest(connection, request)

	response, e := kf.ReadPacket(connection)

	if e != nil {
		t.Error(e.Error())
	}

	if response.(kf.PlayerListResponsePacket).Type != request.Type {
		t.Error("type mismatch")
	}

	if len(response.(kf.PlayerListResponsePacket).Players) != 1 {
		t.Error("wrong number of players in Players field")
	}

	server.Stop()
}
