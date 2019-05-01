package kfnetwork

import ()

const (
	GameStateGameStarted = iota
	GameStateDrawFirstHand
	GameStateMulligan
	GameStateRoundStarted
	GameStateRoundEnded
	GameStateTurnStarted
	GameStateTurnEnded
)

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
