package playnine_test

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

var testStrategyFirstTwoOpeningFlips playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, 1}
}

var testPlayer = playnine.NewPlayer(testStrategyFirstTwoOpeningFlips, nil, nil)

func TestNewGameStartsOnRoundOne(t *testing.T) {
	g, err := playnine.NewGame(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, g.CurrentRound())
}

func TestNewGameStartsWithFlippedCards(t *testing.T) {
	const numPlayers = 4

	players := make([]playnine.Player, numPlayers)

	for i := range numPlayers {
		players[i] = playnine.NewPlayer(testStrategyFirstTwoOpeningFlips, nil, nil)
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	states := g.PlayerStates()
	assert.Len(t, states, numPlayers, "Unexpected number of player states returned")

	for i, state := range states {
		numFlipped := 0
		for _, card := range state.CurrentBoard() {
			if card.FaceUp {
				numFlipped++
			}
		}

		assert.Equal(t, 2, numFlipped, "Player #%d had unexpected number of flipped cards", i+1)
	}
}

func TestNewGameHasRandomizedDiscard(t *testing.T) {
	// Max runs, not total runs
	const runCount = 10000

	players := []playnine.Player{testPlayer}

	sawNonZeroCount := 0

	// The chances of getting a 0 card this many times in a row is hilariously small,
	// but technically possible... this is fine for now
	for range runCount {
		g, err := playnine.NewGame(players)
		assert.Nil(t, err)

		card := g.AvailableDiscard()

		if card != 0 {
			return
		}
	}

	t.Errorf("didn't see a non-zero discard after %d runs, unlikely to be randomized", sawNonZeroCount)
}

func TestNewGameErrorsWhenOpeningStrategyFlipsOnlyOneCard(t *testing.T) {
	brokenFlipStrat := func() [2]int {
		return [2]int{0, 0}
	}

	player := playnine.NewPlayer(brokenFlipStrat, nil, nil)

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

func TestGetCurrentPlayerState(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	currentState := g.CurrentPlayerState()
	currentStateByIndex := g.PlayerStates()[g.CurrentPlayerIndex()]

	assert.Equal(t, currentStateByIndex, currentState)
}
