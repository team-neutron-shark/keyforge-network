package kfnetwork

type Affect interface {
	Duration() uint
	Type() uint
	Card() *Card
	IsPermanent() bool
}

type PlayerAffect struct {
	duration     uint
	affectType   uint
	card         *Card
	permanent    bool
	buffAmount   uint
	debuffAmount uint
}

func NewPlayerAffect() *PlayerAffect {
	playerAffect := new(PlayerAffect)
	return playerAffect
}

func (p *PlayerAffect) Duration() uint {
	return p.duration
}

func (p *PlayerAffect) Type() uint {
	return p.affectType
}

func (p *PlayerAffect) Card() *Card {
	return p.card
}

func (p *PlayerAffect) IsPermanent() bool {
	return p.permanent
}

func (p *PlayerAffect) SetDuration(d uint) {
	p.duration = d
}

func (p *PlayerAffect) SetType(t uint) {
	p.affectType = t
}

func (p *PlayerAffect) SetCard(c *Card) {
	p.card = c
}

func (p *PlayerAffect) SetPermanent(b bool) {
	p.permanent = b
}
