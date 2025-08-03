package playnine_test

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

var testStrategyFirstTwoOpeningFlips playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, 1}
}

var testPlayer = playnine.NewPlayer(testStrategyFirstTwoOpeningFlips, nil)

func TestNewGameStartsOnRoundOne(t *testing.T) {
	g, err := playnine.NewGame(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, g.CurrentRound())
}

func TestNewGameStartsWithFlippedCards(t *testing.T) {
	players := []playnine.Player{}

	_, err := playnine.NewGame(players)
	assert.Nil(t, err)
}

func TestNewGameErrorsWhenOpeningStrategyFlipsOnlyOneCard(t *testing.T) {
	brokenFlipStrat := func() [2]int {
		return [2]int{0, 0}
	}

	player := playnine.NewPlayer(brokenFlipStrat, nil)

	_, err := playnine.NewGame([]playnine.Player{player})

	assert.NotNil(t, err, "Expected error when flip strategy creates duplicate entries")
}

func TestNewPlayerStateIsntFinished(t *testing.T) {
	players := []playnine.Player{testPlayer}
	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	state := g.PlayerStates()[0]
	assert.False(t, state.IsFinished(), "State shouldn't be finished")
}

func TestTwoPlayersHaveDifferentStarts(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	// Technically there's an astronomically small chance of these
	// being the same, but for now that's fine
	states := g.PlayerStates()
	assert.Len(t, states, 2)
	assert.NotElementsMatch(t, states[0].CurrentBoard(), states[1].CurrentBoard())
}
