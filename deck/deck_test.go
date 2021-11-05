package main

import (
	"fmt"
	"github.com/maxatome/go-testdeep/td"
	"math/rand"
	"testing"
)

func ExampleCard_String() {
	fmt.Println(Card{Ace, Club})
	fmt.Println(Card{King, Spade})
	fmt.Println(Card{Seven, Diamond})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Clubs
	// King of Spades
	// Seven of Diamonds
	// Joker
}

func TestNew(t *testing.T) {
	deck := New()
	if len(deck) != 52 {
		t.Errorf("expected: 13 * 4 = 52 cards in a deck, got: %d", len(deck))
	}
}

func TestJokers(t *testing.T) {
	expected := 3
	deck := New(Jokers(expected))
	nJokers := 0
	for _, c := range deck {
		if c.Suit == Joker {
			nJokers++
		}
	}
	td.Cmp(t, nJokers, expected)
}

func TestFilter(t *testing.T) {
	ranks := []Rank{King, Queen, Jack}
	deck := New(Filter(WithoutRanks(ranks)))
	if len(deck) != 40 {
		t.Errorf("expected: 10 * 4 = 40 cards in a deck, got: %d", len(deck))
	}
}

func TestShuffle(t *testing.T) {
	// Fix random seed for testing
	randSource = rand.New(rand.NewSource(0))
	deck := Shuffle(New())
	td.Cmp(t, deck[5], Card{Eight, Spade})
	td.Cmp(t, deck[10], Card{Nine, Diamond})
	td.Cmp(t, deck[15], Card{Queen, Spade})
}

func TestDuplicate(t *testing.T) {
	deck := Duplicate(New())
	if len(deck) != 104 {
		t.Errorf("expected: 2 * 4 * 13 = 104 cards in a deck, got: %d", len(deck))
	}
}
