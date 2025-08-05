package playnine

import (
	"errors"
	"math/rand/v2"
)

// ErrDeckEmpty indicates the deck is empty.
var ErrDeckEmpty error = errors.New("deck is empty")

// Deck contains some number of remaining valid cards.
type Deck struct {
	cards []Card
}

// NewDeck returns a new Deck of 108 cards.
func NewDeck() Deck {
	const (
		numHoleInOnes       = 4
		numRegularCards     = 8
		numRegularCardTypes = 13 // 0-12
	)

	cards := make([]Card, 108)

	i := 0

	for range numHoleInOnes {
		cards[i] = CardHoleInOne
		i++
	}

	for c := range numRegularCardTypes {
		for range numRegularCards {
			cards[i] = Card(c)
			i++
		}
	}

	return Deck{
		cards,
	}
}

// RemainingCardCount returns how many cards are left in the deck.
func (d *Deck) RemainingCardCount() int {
	return len(d.cards)
}

// draw returns a random card from the deck. Returns ErrDeckEmpty
// if no cards remain in the deck.
func (d *Deck) draw() (Card, error) {
	l := len(d.cards)

	if l == 0 {
		return 0, ErrDeckEmpty
	}

	// Pick a card at random, swap it with the last, then discard last
	i := rand.IntN(l)
	c := d.cards[i]
	d.cards[i] = d.cards[l-1]
	d.cards = d.cards[:l-1]

	return c, nil
}
