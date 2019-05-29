package kfnetwork

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

var playerManagerOnce sync.Once
var playerManagerSingleton *PlayerManager

type PlayerManager struct {
	players []*Player
}

func Players() *PlayerManager {
	playerManagerOnce.Do(func() {
		playerManagerSingleton = new(PlayerManager)
	})

	return playerManagerSingleton
}

func (p *PlayerManager) AddPlayer(player *Player) {
	logEntry := fmt.Sprintf("Player %s logged in.", player.Name)
	Logger().Log(logEntry)

	if !p.PlayerExists(player) {
		player.Lock()
		defer player.Unlock()

		p.players = append(p.players, player)
	}
}

func (p *PlayerManager) RemovePlayer(player *Player) {
	players := []*Player{}

	for _, p := range p.players {
		p.Lock()
		defer p.Unlock()

		if p != player {
			players = append(players, p)
		}
	}

	p.players = players
}

func (p *PlayerManager) PlayerExists(client *Player) bool {
	for _, player := range p.players {
		player.Lock()
		defer player.Unlock()

		if player == client {
			return true
		}
	}

	return false
}

func (p *PlayerManager) FindPlayerByID(id string) (*Player, error) {
	for _, player := range p.players {
		player.Lock()
		defer player.Unlock()

		if player.ID == id {
			return player, nil
		}
	}

	return &Player{}, errors.New("no such player found")
}

func (p *PlayerManager) FindPlayerByConnection(connection net.Conn) (*Player, error) {
	for _, player := range p.players {
		player.Lock()
		defer player.Unlock()

		if player.Client == connection {
			return player, nil
		}
	}

	return &Player{}, errors.New("could not find any players with the given connection")
}

// PlayerHasLobby - This function is used to determine whether or not a player
// is in a lobby.
func (p *PlayerManager) PlayerHasLobby(player *Player) bool {
	for _, lobby := range Lobbies().GetLobbies() {
		if lobby.PlayerExists(player) {
			return true
		}
	}

	return false
}

func (p *PlayerManager) GetPlayers() []*Player {
	return p.players
}
