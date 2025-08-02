package playnine_test

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

const (
	playerBoardSize = 8
)

func TestPlayerDrawsFromDeckOnCreation(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize)

	_, err := playnine.NewPlayerFromDeck(&d)

	assert.Nil(t, err)

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize-playerBoardSize)
}

func TestTwoPlayersHaveDifferentStarts(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize)

	p1, err := playnine.NewPlayerFromDeck(&d)
	assert.Nil(t, err)

	p2, err := playnine.NewPlayerFromDeck(&d)
	assert.Nil(t, err)

	// Technically there's an astronomically small chance of these
	// being the same, but for now that's fine
	assert.NotElementsMatch(t, p1.CurrentBoard(), p2.CurrentBoard())
}
