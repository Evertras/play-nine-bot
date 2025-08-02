package playnine_test

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

const (
	playerBoardSize = 8
)

var testStrategyFirstTwoOpeningFlips playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, 1}
}

func TestPlayerDrawsFromDeckOnCreationAndFlipsTwoCards(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize)

	p, err := playnine.NewPlayerFromDeck(&d, testStrategyFirstTwoOpeningFlips, nil)

	assert.Nil(t, err)

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize-playerBoardSize)

	seenFaceUp := 0

	for _, c := range p.CurrentBoard() {
		if c.FaceUp {
			seenFaceUp++
		}
	}

	assert.Equal(t, seenFaceUp, 2, "Didn't see two cards face up")
}

func TestTwoPlayersHaveDifferentStarts(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize)

	p1, err := playnine.NewPlayerFromDeck(&d, testStrategyFirstTwoOpeningFlips, nil)
	assert.Nil(t, err)

	p2, err := playnine.NewPlayerFromDeck(&d, testStrategyFirstTwoOpeningFlips, nil)
	assert.Nil(t, err)

	// Technically there's an astronomically small chance of these
	// being the same, but for now that's fine
	assert.NotElementsMatch(t, p1.CurrentBoard(), p2.CurrentBoard())
}

func TestNewPlayerIsntFinished(t *testing.T) {
	d := playnine.NewDeck()
	p, err := playnine.NewPlayerFromDeck(&d, testStrategyFirstTwoOpeningFlips, nil)
	assert.Nil(t, err)

	assert.False(t, p.IsFinished())
}

func TestNewPlayerErrorsWhenOpeningFlipsOnlyOneCard(t *testing.T) {
	d := playnine.NewDeck()

	_, err := playnine.NewPlayerFromDeck(&d, func() [2]int {
		return [2]int{0, 0}
	}, nil)

	assert.NotNil(t, err)
}
