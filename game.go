package kfnetwork

import (
	"errors"
	keyforge "keyforge/game"
	"net"
)

const (
	GameStateGameStarted = iota
	GameStateDrawFirstHand
	GameStateMulligan
	GameStateRoundStarted
	GameStateRoundEnded
	GameStateTurnStarted
	GameStateTurnEnded
)

// PlayerClient - This type holds both the keyforge player type along with
// the net.Conn object required for networked communication.
type PlayerClient struct {
	Active bool
	Client net.Conn
	keyforge.Player
}

type Game struct {
	Seed    int64
	Turn    int
	Round   int
	Running bool
	Players []PlayerClient
}

func NewGame() *Game {
	game := new(Game)
	return game
}

func (g *Game) FindActivePlayer() (PlayerClient, error) {
	for _, player := range g.Players {
		if player.Active {
			return player, nil
		}
	}

	return PlayerClient{}, errors.New("no active player found")
}

func (g *Game) AdvanceTurn() {
	if g.Turn == 0 {
		g.Turn = 1
		return
	}

	if g.Turn > len(g.Players) {
		g.Turn = 1
		return
	}

	g.Turn++
}

func (g *Game) AdvanceRound() {
	g.Round++
}
