package playnine_test

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

const (
	playerBoardSize = 8
)

func TestPlayerDrawsFromDeckOnCreationAndFlipsTwoCards(t *testing.T) {
	d := playnine.NewDeck()

	assert.Equal(t, d.RemainingCardCount(), expectedDeckSize)

	p, err := testPlayer.StartGame(&d)

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
