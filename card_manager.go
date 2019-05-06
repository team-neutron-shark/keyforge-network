package kfnetwork

import (
	"bytes"
	"encoding/json"
	"errors"
	keyforge "keyforge/game"
	"os"
)

type Card struct {
	keyforge.Card
}

type CardQuery struct {
	id        string
	number    int
	expansion int
}

func NewCardQuery() *CardQuery {
	cardQuery := new(CardQuery)
	return cardQuery
}

type CardManager struct {
	cards []Card
}

func NewCardManager() *CardManager {
	cardManager := new(CardManager)
	return cardManager
}

func (q *CardQuery) ID() string {
	return q.id
}

func (q *CardQuery) SetID(id string) {
	q.id = id
}

func (q *CardQuery) Number() int {
	return q.number
}

func (q *CardQuery) SetNumber(number int) {
	q.number = number
}

func (q *CardQuery) Expansion() int {
	return q.expansion
}

func (q *CardQuery) SetExpansion(expansion int) {
	q.expansion = expansion
}

func (c *CardManager) LoadFromFile(location string) error {
	cards := []Card{}
	buffer := bytes.Buffer{}

	file, e := os.Open(location)

	if e != nil {
		return e
	}

	buffer.ReadFrom(file)

	if e != nil {
		return e
	}

	e = json.Unmarshal(buffer.Bytes(), &cards)

	if e != nil {
		return e
	}

	return nil
}

func (c *CardManager) WriteToFile(location string) error {
	bytes, e := json.MarshalIndent(c.cards, "", "    ")

	if e != nil {
		return e
	}

	file, e := os.Create(location)

	if e != nil {
		return e
	}

	file.WriteString(string(bytes))
	file.Close()

	return nil
}

func (c *CardManager) QueryCard(query *CardQuery) (Card, error) {
	for _, card := range c.cards {
		if query.ID() == card.ID && query.Expansion() == card.Expansion && query.Number() == card.CardNumber {
			return card, nil
		}
	}

	for _, card := range c.cards {
		if query.Expansion() == card.Expansion && query.Number() == card.CardNumber {
			return card, nil
		}
	}

	return Card{}, errors.New("no card found with the given query")
}

func (c *CardManager) CardExists(query *CardQuery) bool {
	for _, card := range c.cards {
		if query.ID() == card.ID && query.Expansion() == card.Expansion && query.Number() == card.CardNumber {
			return true
		}
	}

	for _, card := range c.cards {
		if query.Expansion() == card.Expansion && query.Number() == card.CardNumber {
			return true
		}
	}

	return false
}
