package main

import (
	"fmt"
	"math/rand"
	"time"
)

//go:generate stringer -type=Suit,Rank

type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Heart
	Club
	Joker
)

var suits = [...]Suit{Spade, Diamond, Heart, Club}

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

func New(options ...func([]Card) []Card) []Card {
	var cards []Card
	for _, s := range suits {
		for r := Ace; r <= King; r++ {
			cards = append(cards, Card{r, s})
		}
	}
	// Apply Functional Options
	for _, opt := range options {
		cards = opt(cards)
	}
	return cards
}

func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{Suit: Joker})
		}
		return cards
	}
}

func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

func WithoutRanks(ranks []Rank) func(Card) bool {
	return func(c Card) bool {
		for _, r := range ranks {
			if c.Rank == r {
				return false
			}
		}
		return true
	}
}

var randSource = rand.New(rand.NewSource(time.Now().Unix()))

func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	for i, j := range randSource.Perm(len(cards)) {
		ret[i] = cards[j]
	}
	return ret
}

func Duplicate(deck []Card) []Card {
	return append(deck, deck...)
}
