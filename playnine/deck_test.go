package playnine_test

import (
	"fmt"
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

const (
	expectedDeckSize = 108
	maxCardValue     = 12
	minCardValue     = 0
)

func TestNewDeckHas108Cards(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, expectedDeckSize, d.RemainingCardCount())
}

func TestDeckDrawsValidCardsUntilEmpty(t *testing.T) {
	d := playnine.NewDeck()

	for range expectedDeckSize {
		c, err := d.Draw()

		assert.Nil(t, err)

		if c == playnine.CardHoleInOne {
			continue
		}

		// 0-12
		assert.LessOrEqual(t, c, playnine.Card(12))
		assert.GreaterOrEqual(t, c, playnine.Card(0))
	}

	_, err := d.Draw()

	assert.ErrorIs(t, err, playnine.ErrDeckEmpty)
}

func TestDeckContainsExpectedSpreadOfCards(t *testing.T) {
	// Counts per deck
	const (
		expectedHoleInOnes       = 4
		expectedRegularCards     = 8
		expectedRegularCardTypes = 13 // 0-12 inclusive
	)

	d := playnine.NewDeck()

	sawCardCount := make(map[playnine.Card]int)

	for range expectedDeckSize {
		c, err := d.Draw()
		assert.Nil(t, err)

		sawCardCount[c]++
	}

	assert.Equal(t, expectedHoleInOnes, sawCardCount[playnine.CardHoleInOne])

	for c := range expectedRegularCardTypes {
		assert.Equal(t, expectedRegularCards, sawCardCount[playnine.Card(c)])
	}
}

func TestDrawingFromDeckIsRandomized(t *testing.T) {
	d := playnine.NewDeck()

	lastSaw := playnine.Card(-1)
	streakCount := 0

	for range expectedDeckSize {
		c, err := d.Draw()
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
