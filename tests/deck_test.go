package tests

import (
	"testing"

	kf "github.com/team-neutron-shark/keyforge-network"
)

var testLocation = "test_data/test_deck.json"

func TestLoadDeckFromFile(t *testing.T) {
	deck, e := kf.LoadDeckFromFile(testLocation)

	if e != nil {
		t.Error(e.Error())
	}

	if len(deck.Cards) != 36 {
		t.Errorf("Deck only contains %d cards!", len(deck.Cards))
	}
}
