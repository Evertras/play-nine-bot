package playnine_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/evertras/play-nine-bot/playnine"
)

var testStrategyFirstTwoOpeningFlips playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, 1}
}

var testStrategyDrawOrDiscard playnine.PlayerStrategyTakeTurnDrawOrUseDiscard = func(playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
	return playnine.DecisionDrawOrUseDiscardDraw, 0, nil
}

var testStrategyDrawn playnine.PlayerStrategyTakeTurnDrawn = func(g playnine.Game, _ playnine.Card) (playnine.DecisionDrawn, playnine.DecisionCardIndex, error) {
	// The strategy is dumb and simple: ignore everything and just flip cards
	state := g.CurrentPlayerState()

	for index, card := range state.CurrentBoard() {
		if !card.FaceUp {
			return playnine.DecisionDrawnDiscardAndFlip, playnine.DecisionCardIndex(index), nil
		}
	}

	return 0, 0, fmt.Errorf("didn't find a face down card to flip")
}

var testPlayer = playnine.NewPlayer(
	"test",
	testStrategyFirstTwoOpeningFlips,
	testStrategyDrawOrDiscard,
	testStrategyDrawn,
)

func TestNewGameStartsOnRoundOne(t *testing.T) {
	g, err := playnine.NewGame(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1, g.CurrentRound())
}

func TestNewGameStartsWithFlippedCards(t *testing.T) {
	const numPlayers = 4

	players := make([]playnine.Player, numPlayers)

	for i := range numPlayers {
		players[i] = playnine.NewPlayer("test", testStrategyFirstTwoOpeningFlips, nil, nil)
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
	const maxRunCount = 10000

	players := []playnine.Player{testPlayer}

	sawNonZeroCount := 0

	// The chances of getting a 0 card this many times in a row is hilariously small,
	// but technically possible... this is fine for now
	for range maxRunCount {
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

	player := playnine.NewPlayer("broken flipper", brokenFlipStrat, nil, nil)

	_, err := playnine.NewGame([]playnine.Player{player})

	assert.NotNil(t, err, "Expected error when flip strategy creates duplicate entries")
}

func TestNewPlayerStateIsntFinished(t *testing.T) {
	players := []playnine.Player{testPlayer}
	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	assert.Len(t, g.PlayerStates(), 1)

	if t.Failed() {
		return
	}

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
	if t.Failed() {
		return
	}
	assert.NotElementsMatch(t, states[0].CurrentBoard(), states[1].CurrentBoard())
}

func TestGetCurrentPlayerState(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)

	assert.Len(t, g.PlayerStates(), 2)

	if t.Failed() {
		return
	}

	currentState := g.CurrentPlayerState()
	currentStateByIndex := g.PlayerStates()[g.CurrentPlayerIndex()]

	assert.Equal(t, currentStateByIndex, currentState)
}

func TestTakingTurnAdvancesToNextPlayer(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)
	assert.Equal(t, 0, g.CurrentPlayerIndex(), "Should start at player index 0")

	turn := func() {
		err := g.TakeTurn()
		assert.Nil(t, err, "Failed to take turn")
	}

	turn()
	assert.Equal(t, 1, g.CurrentPlayerIndex(), "Should be at player index 1 after the first turn")

	turn()
	assert.Equal(t, 2, g.CurrentPlayerIndex(), "Should be at player index 2 after the second turn")

	turn()
	assert.Equal(t, 0, g.CurrentPlayerIndex(), "Should be at player index 0 again after the third turn")
}

func TestTakingTurnAppliesTurnStrategy(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err)
	assert.Equal(t, 0, g.CurrentPlayerIndex(), "Should start at player index 0")

	turn := func() {
		err := g.TakeTurn()
		assert.Nil(t, err, "Failed to take turn")
	}

	oldState := playnine.PlayerBoard{}

	for i, s := range g.PlayerStates()[0].CurrentBoard() {
		oldState[i] = s
	}

	turn()

	assert.NotElementsMatch(t, oldState, g.PlayerStates()[0].CurrentBoard(), "State didn't update, but should've")
}

func TestRoundAdvancesWhenAllPlayersFinish(t *testing.T) {
	players := []playnine.Player{
		testPlayer,
		testPlayer,
	}

	g, err := playnine.NewGame(players)
	assert.Nil(t, err, "Failed to create new game")

	const expectedTurnsToCompleteRound = 2*6 + 1

	assert.Equal(t, 1, g.CurrentRound(), "Should start on the first round")

	for i := range expectedTurnsToCompleteRound {
		err = g.TakeTurn()
		assert.Nil(t, err, "Failed to take turn on turn %d", i+1)
	}

	assert.Equal(t, 2, g.CurrentRound(), "Should be on the second round")

	scores := g.PlayerRoundScores()

	assert.Len(t, scores, 1, "Should have 1 round of scores")

	if t.Failed() {
		return
	}

	roundOneScores := scores[0]

	assert.Len(t, roundOneScores, len(players), "Should have all players' scores in the round")

	for i, s := range g.PlayerStates() {
		sawFaceUp := 0

		for _, c := range s.CurrentBoard() {
			if c.FaceUp {
				sawFaceUp++
			}
		}

		assert.Equal(t, 2, sawFaceUp, "Player %d didn't have exactly two face up cards for the start of the next round", i+1)
	}

	assert.Equal(t, 1, g.CurrentPlayerIndex(), "Should let the next player deal in the next round")
}
