package kfnetwork

import (
	keyforge "keyforge/game"
	"net"
	"sync"
)

type Player struct {
	Active      bool
	ID          string
	affectMutex sync.Mutex
	affects     []*PlayerAffect
	Client      net.Conn
	Name        string
	Game        *Game
	Debug       bool
	PlayerDeck  Deck
	HandPile    []Card
	DrawPile    []Card
	DiscardPile []Card
	ArchivePile []Card
	PurgePile   []Card
	Artifacts   []Card
	Creatures   []Card
	FirstTurn   bool
	Amber       int
	Keys        int
	Chains      int
}

// NewPlayer - Returns a pointer to a new player object.
func NewPlayer() *Player {
	player := new(Player)
	player.DrawPile = make([]Card, 0)
	player.HandPile = make([]Card, 0)
	player.ArchivePile = make([]Card, 0)
	player.DiscardPile = make([]Card, 0)

	return player
}

func (p *Player) Affects() []*PlayerAffect {
	return p.affects
}

func (p *Player) AddAffect(affect *PlayerAffect) {

	if !p.HasAffect(affect) {
		p.affectMutex.Lock()
		p.affects = append(p.affects, affect)
		p.affectMutex.Unlock()
	}
}

func (p *Player) RemoveAffect(affect *PlayerAffect) {
	returnAffects := []*PlayerAffect{}

	p.affectMutex.Lock()
	defer p.affectMutex.Unlock()

	for _, a := range p.affects {
		if a != affect {
			returnAffects = append(returnAffects, a)
		}
	}

	p.affects = returnAffects
}

func (p *Player) FindAffectByCard(card *keyforge.Card) []*PlayerAffect {
	foundAffects := []*PlayerAffect{}

	for _, affect := range p.affects {
		if affect.Card() == card {
			p.affectMutex.Lock()
			foundAffects = append(foundAffects, affect)
			p.affectMutex.Unlock()
		}
	}

	return foundAffects
}

func (p *Player) HasAffect(affect *PlayerAffect) bool {
	for _, a := range p.affects {
		if a == affect {
			return true
		}
	}

	return false
}
