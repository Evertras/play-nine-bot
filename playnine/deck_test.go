package playnine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	expectedDeckSize = 108
	maxCardValue     = 12
	minCardValue     = 0
)

func TestNewDeckHas108Cards(t *testing.T) {
	d := NewDeck()

	assert.Equal(t, expectedDeckSize, d.RemainingCardCount())
}

func TestDeckDrawsValidCardsUntilEmpty(t *testing.T) {
	d := NewDeck()

	for range expectedDeckSize {
		c, err := d.draw()

		assert.Nil(t, err)

		if c == CardHoleInOne {
			continue
		}

		// 0-12
		assert.LessOrEqual(t, c, Card(12))
		assert.GreaterOrEqual(t, c, Card(0))
	}

	_, err := d.draw()

	assert.ErrorIs(t, err, ErrDeckEmpty)
}

func TestDeckContainsExpectedSpreadOfCards(t *testing.T) {
	// Counts per deck
	const (
		expectedHoleInOnes       = 4
		expectedRegularCards     = 8
		expectedRegularCardTypes = 13 // 0-12 inclusive
	)

	d := NewDeck()

	sawCardCount := make(map[Card]int)

	for range expectedDeckSize {
		c, err := d.draw()
		assert.Nil(t, err)

		sawCardCount[c]++
	}

	assert.Equal(t, expectedHoleInOnes, sawCardCount[CardHoleInOne])

	for c := range expectedRegularCardTypes {
		assert.Equal(t, expectedRegularCards, sawCardCount[Card(c)])
	}
}

func TestDrawingFromDeckIsRandomized(t *testing.T) {
	d := NewDeck()

	lastSaw := Card(-1)
	streakCount := 0

	for range expectedDeckSize {
		c, err := d.draw()
		assert.Nil(t, err)

		if lastSaw == c {
			streakCount++

			assert.NotEqual(t, 8, streakCount, fmt.Sprintf("Saw 8 in a row of card %d", c))
		} else {
			lastSaw = c
			streakCount = 0
		}
	}
}
