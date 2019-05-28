package tests

import (
	"testing"

	kf "github.com/team-neutron-shark/keyforge-network"
)

func TestServerAddLobby(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()
	name := "test lobby"

	server.AddLobby(player, name)

	if len(server.Lobbies) != 1 {
		t.Error("failed to add lobby to lobby array")
	}
}

func TestServerRemoveLobby(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()
	name := "test lobby"

	server.AddLobby(player, name)

	if len(server.Lobbies) != 1 {
		t.Error("failed to add lobby to lobby array")
	}

	lobby, e := server.FindLobbyByName(name)

	if e != nil {
		t.Error(e.Error())
	}

	server.RemoveLobby(lobby)

	if len(server.Lobbies) != 0 {
		t.Error("failed to remove lobby from lobby array")
	}
}

func TestServerFindLobbyByName(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()
	name := "test lobby"

	server.AddLobby(player, name)

	if len(server.Lobbies) != 1 {
		t.Error("failed to add lobby to lobby array")
	}

	lobby, e := server.FindLobbyByName(name)

	if e != nil {
		t.Error(e.Error())
	}

	if lobby.Name() != name {
		t.Error("lobby names do not match")
	}
}

func TestServerFindLobbyByID(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()
	name := "test lobby"

	server.AddLobby(player, name)

	if len(server.Lobbies) != 1 {
		t.Error("failed to add lobby to lobby array")
	}

	lobby, e := server.FindLobbyByName(name)

	if e != nil {
		t.Error(e.Error())
	}

	if lobby.Name() != name {
		t.Error("lobby names do not match")
	}

	newLobby, e := server.FindLobbyByID(lobby.ID())

	if e != nil {
		t.Error(e.Error())
	}

	if newLobby.Name() != name {
		t.Error("lobby names do not match")
	}

	if newLobby.ID() != lobby.ID() {
		t.Error("lobby IDs do not match")
	}
}

func TestServerAddPlayer(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()

	server.AddPlayer(player)

	if len(server.Clients) != 1 {
		t.Error("server does not have the correct number of players")
	}
}

func TestServerRemovePlayer(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()

	server.AddPlayer(player)

	if len(server.Clients) != 1 {
		t.Error("server did not add the correct number of players")
	}

	server.RemovePlayer(player)

	if len(server.Clients) != 0 {
		t.Error("server did not remove the correct number of players")
	}
}

func TestServerFindPlayerByConnection(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()
	connection := MockNetworkConnection{}

	player.Client = connection
	server.AddPlayer(player)

	searchPlayer, e := server.FindPlayerByConnection(connection)

	if e != nil {
		t.Error(e.Error())
	}

	if searchPlayer.Client != connection {
		t.Error("player connections do not match")
	}

	if searchPlayer != player {
		t.Error("player returned does not match the player searched")
	}
}

func TestServerPlayerExists(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()

	server.AddPlayer(player)

	if !server.PlayerExists(player) {
		t.Error("player does not exist")
	}

}

func TestServerPlayerHasLobby(t *testing.T) {
	server := kf.NewServer(":4321")
	defer server.Stop()
	player := kf.NewPlayer()

	server.AddPlayer(player)
	server.AddLobby(player, "testing")

	if !server.PlayerHasLobby(player) {
		t.Error("player does not have a lobby")
	}
}
