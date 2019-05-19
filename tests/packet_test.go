package tests

import (
	"fmt"
	"testing"

	kfnetwork "github.com/team-neutron-shark/keyforge-network"
)

func TestReadWritePacket(t *testing.T) {
	testConnection := NewMockNetworkConnection()

	packet := kfnetwork.VersionPacket{}
	packet.Sequence = 12345
	packet.Type = kfnetwork.PacketTypeVersionRequest
	packet.Version = 1.23

	e := kfnetwork.WritePacket(testConnection, packet)

	if e != nil {
		errorMessage := fmt.Sprintf("error writing packet - %s", e.Error())
		t.Error(errorMessage)
	}

	packetResult, e := kfnetwork.ReadPacket(testConnection)

	if e != nil {
		errorMessage := fmt.Sprintf("error reading packet - %s", e.Error())
		t.Error(errorMessage)
	}

	if packetResult.(kfnetwork.VersionPacket).Sequence != 12345 {
		t.Error("Packet sequence not read correctly.")
	}

	if packetResult.(kfnetwork.VersionPacket).Type != kfnetwork.PacketTypeVersionRequest {
		t.Error("Packet type not read correctly.")
	}

	if packetResult.(kfnetwork.VersionPacket).Version != 1.23 {
		t.Error("Packet version not read correctly")
	}
}
