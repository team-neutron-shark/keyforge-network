package tests

import (
	"testing"

	kf "github.com/team-neutron-shark/keyforge-network"
)

func TestPlayerAddAffect(t *testing.T) {
	player := kf.NewPlayer()

	// Simulate a one-turn inability to forge a la Miasma.
	affect := kf.NewPlayerAffect()
	affect.SetType(kf.PlayerAffectCannotForge)
	affect.SetPermanent(false)
	affect.SetDuration(1)

	player.AddAffect(affect)

	affects := player.Affects()

	if len(affects) != 1 {
		t.Error("player should have one affect")
	}

	if affects[0].Type() != kf.PlayerAffectCannotForge {
		t.Error("affect type does not match the added affect")
	}

	if affects[0].Duration() != 1 {
		t.Error("duration does not match the added affect")
	}

	if affects[0].IsPermanent() {
		t.Error("permanent flag does not match the added affect")
	}
}

func TestPlayerAddSameAffectTwice(t *testing.T) {
	player := kf.NewPlayer()

	// Simulate a one-turn inability to forge a la Miasma.
	miasmaAffect := kf.NewPlayerAffect()
	miasmaAffect.SetType(kf.PlayerAffectCannotForge)
	miasmaAffect.SetPermanent(false)
	miasmaAffect.SetDuration(1)

	player.AddAffect(miasmaAffect)
	player.AddAffect(miasmaAffect)

	affects := player.Affects()

	if len(affects) != 1 {
		t.Error("player should have one affect")
	}
}

func TestPlayerRemoveAffect(t *testing.T) {
	player := kf.NewPlayer()

	// Simulate a one-turn inability to forge a la Miasma.
	miasmaAffect := kf.NewPlayerAffect()
	miasmaAffect.SetType(kf.PlayerAffectCannotForge)
	miasmaAffect.SetPermanent(false)
	miasmaAffect.SetDuration(1)

	player.AddAffect(miasmaAffect)

	player.RemoveAffect(miasmaAffect)

	affects := player.Affects()

	if len(affects) != 0 {
		t.Error("player should have no affects")
	}
}
