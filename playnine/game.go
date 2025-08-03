package playnine

import "fmt"

// TotalRounds is how many rounds are played in total.
const TotalRounds = 9

// Game represents the state of the game and can advance forward through
// each player's turn.
type Game struct {
	deck            *Deck
	players         []Player
	playerStates    []PlayerState
	playerTurnIndex int

	round int
}

// NewGame creates a new game with the given players ready to play.
func NewGame(players []Player) (Game, error) {
	d := NewDeck()

	playerStates := make([]PlayerState, len(players))

	for i, p := range players {
		state, err := p.StartGame(&d)

		if err != nil {
			return Game{}, fmt.Errorf("failed to create new player state for index %d: %w", i, err)
		}

		playerStates[i] = state
	}

	return Game{
		deck:            &d,
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

// PlayerStates gets the current round's player states.
func (g Game) PlayerStates() []PlayerState {
	return g.playerStates
}
