package kfnetwork

import (
	keyforge "keyforge/game"
	"net"
	"sync"
)

// PlayerClient - This type holds both the keyforge player type along with
// the net.Conn object required for networked communication.
type PlayerClient struct {
	Active      bool
	ID          string
	affectMutex sync.Mutex
	affects     []*PlayerAffect
	Client      net.Conn
	keyforge.Player
}

func NewPlayerClient() *PlayerClient {
	playerClient := new(PlayerClient)
	return playerClient
}

func (p *PlayerClient) Affects() []*PlayerAffect {
	return p.affects
}

func (p *PlayerClient) AddAffect(affect *PlayerAffect) {

	if !p.HasAffect(affect) {
		p.affectMutex.Lock()
		p.affects = append(p.affects, affect)
		p.affectMutex.Unlock()
	}
}

func (p *PlayerClient) RemoveAffect(affect *PlayerAffect) {
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

func (p *PlayerClient) FindAffectByCard(card *keyforge.Card) []*PlayerAffect {
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

func (p *PlayerClient) HasAffect(affect *PlayerAffect) bool {
	for _, a := range p.affects {
		if a == affect {
			return true
		}
	}

	return false
}
