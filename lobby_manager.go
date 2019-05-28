package kfnetwork

import (
	"errors"
	"sync"
)

var lobbyManagerOnce sync.Once
var lobbyManagerSingleton *LobbyManager

type LobbyManager struct {
	lobbies []*Lobby
}

// Logger - Function used to access the LogManager singleton pointer.
func Lobbies() *LobbyManager {
	lobbyManagerOnce.Do(func() {
		lobbyManagerSingleton = new(LobbyManager)
	})

	return lobbyManagerSingleton
}

func (l *LobbyManager) AddLobby(creator *Player, name string) *Lobby {
	lobby := NewLobby()
	lobby.SetID(GenerateUUID())
	lobby.AddPlayer(creator)
	lobby.SetHost(creator)
	lobby.SetName(name)

	l.lobbies = append(l.lobbies, lobby)

	return lobby
}

func (l *LobbyManager) RemoveLobby(lobby *Lobby) {
	lobbies := []*Lobby{}

	for _, l := range l.lobbies {
		if l != lobby {
			lobbies = append(lobbies, l)
		}
	}

	l.lobbies = lobbies
}

func (l *LobbyManager) FindLobbyByID(id string) (*Lobby, error) {
	for _, lobby := range l.lobbies {
		if lobby.ID() == id {
			return lobby, nil
		}
	}

	return &Lobby{}, errors.New("no lobby found with the given ID")
}

func (l *LobbyManager) FindLobbyByName(name string) (*Lobby, error) {
	for _, lobby := range l.lobbies {
		if lobby.name == name {
			return lobby, nil
		}
	}

	return &Lobby{}, errors.New("no lobby found with the given ID")
}

func (l *LobbyManager) FindLobbyByPlayer(player *Player) (*Lobby, error) {
	for _, lobby := range l.lobbies {
		if lobby.PlayerExists(player) {
			return lobby, nil
		}
	}

	return &Lobby{}, errors.New("no lobby found with the given player")
}

func (l *LobbyManager) GetLobbies() []*Lobby {
	return l.lobbies
}
