package playnine

import "fmt"

// TotalRounds is how many rounds are played in total.
const TotalRounds = 9

// Game represents the state of the game and can advance forward through
// each player's turn.
type Game struct {
	deck            *Deck
	discarded       Card
	players         []Player
	playerStates    []PlayerState
	playerTurnIndex int
	finalTurn       bool

	round int
}

// NewGame creates a new game with the given players ready to play.
func NewGame(players []Player) (Game, error) {
	d := NewDeck()

	playerStates := make([]PlayerState, len(players))

	for i, p := range players {
		state, err := p.startGame(&d)

		if err != nil {
			return Game{}, fmt.Errorf("failed to create new player state for index %d: %w", i, err)
		}

		playerStates[i] = state
	}

	discarded, err := d.draw()
	if err != nil {
		return Game{}, fmt.Errorf("failed to draw for discard: %w", err)
	}

	return Game{
		deck:            &d,
		discarded:       discarded,
		players:         players,
		playerTurnIndex: 0,
		playerStates:    playerStates,

		round: 1,
	}, nil
}

// CurrentRound returns the current round the game is on,
// from 1-9 inclusive.
func (g Game) CurrentRound() int {
	return g.round
}

// CurrentPlayerIndex returns the current player's index so a player strategy
// can know who they are. 0-indexed.
func (g Game) CurrentPlayerIndex() int {
	return g.playerTurnIndex
}

// PlayerStates gets the current round's player states.
func (g Game) PlayerStates() []PlayerState {
	return g.playerStates
}

// AvailableDiscard returns the card that's available for choosing to discard.
func (g Game) AvailableDiscard() Card {
	return g.discarded
}
