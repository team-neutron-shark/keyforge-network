package kfnetwork

import (
	"errors"
)

const (
	GameStateGameStarted uint = iota
	GameStateDrawFirstHand
	GameStateMulligan
	GameStateRoundStarted
	GameStateRoundEnded
	GameStateTurnStarted
	GameStateTurnEnded
)

const (
	PlayerAffectAdditionalForgeCost uint = iota
	PlayerAffectCannotForge
	PlayerAffectCreaturesCannotFight
	PlayerAffectPlayCardLimit
	PlayerAffectDrawReduction
	PlayerAffectDrawBonus
	PlayerAffectPlayCreatureLimit
	PlayerAffectCannotPlayBrobnar
	PlayerAffectCannotPlayDis
	PlayerAffectCannotPlayLogos
	PlayerAffectCannotPlayMars
	PlayerAffectCannotPlaySanctum
	PlayerAffectCannotPlayShadows
	PlayerAffectCannotPlayUntamed
	PlayerAffectMustPlayBrobnar
	PlayerAffectMustPlayDis
	PlayerAffectMustPlayLogos
	PlayerAffectMustPlayMars
	PlayerAffectMustPlaySanctum
	PlayerAffectMustPlayShadows
	PlayerAffectMustPlayUntamed
	PlayerAffectPlayArtifactAmberPenalty
)

const (
	UpgradeAffectPower uint = iota
	UpgradeAffectGainArmor
	UpgradeAffectGainTaunt
	UpgradeAffectGainSteal
	UpgradeAffectGainCapture
	UpgradeAffectGainAssault
	UpgradeAffectGainElusive
	UpgradeAffectGainSkirmish
	UpgradeAffectGainHazardous
	UpgradeAffectGainAllHouses
	UpgradeAffectGainActiveHouse
	UpgradeAffectStunExhaust
	UpgradeAffectReadySameCreature
	UpgradeAffectArchiveOnDeath
	UpgradeAffectSwapBattleLine
	UpgradeAffectIncreaseOpponentForgeCost
	UpgradeAffectReapDamage
	UpgradeAffectDestroyedDamageEachCreature
	UpgradeAffectDestroyMostPowerfulCreature
)

type PlayerListEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlayerList struct {
	Count   uint              `json:"count"`
	Players []PlayerListEntry `json:"players"`
}

type Game struct {
	Seed    int64
	Turn    int
	Round   int
	Running bool
	Players []*Player
}

func NewGame() *Game {
	game := new(Game)
	return game
}

func (g *Game) Start() {

}

func (g *Game) FindActivePlayer() (*Player, error) {
	for _, player := range g.Players {
		if player.Active {
			return player, nil
		}
	}

	return &Player{}, errors.New("no active player found")
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
