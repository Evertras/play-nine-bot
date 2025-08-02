package playnine

import "fmt"

// PlayerBoardSize is how many cards are on the player's board,
// which are then divided into two rows.
const PlayerBoardSize = 8

// PlayerBoardCard represents a card on the player's board, and
// whether it's flipped up so that they know what it is or not.
type PlayerBoardCard struct {
	Card   Card
	FaceUp bool
}

// PlayerStrategyOpeningFlips must return two unique indices of the first two
// cards to flip.
type PlayerStrategyOpeningFlips func() [2]int

// PlayerStrategyTakeTurn takes in the state of the game and takes a turn,
// mutating the underlying game state (for now).
type PlayerStrategyTakeTurn func(d *Deck)

// Player contains information for a single player in the game
type Player struct {
	board [PlayerBoardSize]PlayerBoardCard

	strategyOpeningFlips PlayerStrategyOpeningFlips
	strategyTakeTurn     PlayerStrategyTakeTurn
}

// NewPlayerFromDeck draws a starting hand for the player and initializes
// a board, then executes the opening flip strategy to flip exactly two cards.
//
// Returns an error if the deck cannot be drawn from or the opening flip
// strategy doesn't flip exactly two unique cards.
func NewPlayerFromDeck(d *Deck, strategyOpeningFlips PlayerStrategyOpeningFlips, strategyTakeTurn PlayerStrategyTakeTurn) (Player, error) {
	if d == nil {
		return Player{}, fmt.Errorf("given deck was nil")
	}

	board := [PlayerBoardSize]PlayerBoardCard{}

	for i := range PlayerBoardSize {
		c, err := d.Draw()

		if err != nil {
			return Player{}, fmt.Errorf("tried to draw starting board for player: %w", err)
		}

		board[i] = PlayerBoardCard{
			Card: c,
		}
	}

	flipIndices := strategyOpeningFlips()

	if flipIndices[0] == flipIndices[1] {
		return Player{}, fmt.Errorf("opening flip strategy failed to produce unique card indices to flip, both returned values were %d", flipIndices[0])
	}

	for _, i := range flipIndices {
		if i < 0 || i >= PlayerBoardSize {
			return Player{}, fmt.Errorf("index to flip was out of range, should be between 0 and 7 but got %d", i)
		}

		board[i].FaceUp = true
	}

	return Player{
		board:                board,
		strategyOpeningFlips: strategyOpeningFlips,
		strategyTakeTurn:     strategyTakeTurn,
	}, nil
}

// CurrentBoard returns the player's current board state.
func (p Player) CurrentBoard() [PlayerBoardSize]PlayerBoardCard {
	return p.board
}

// IsFinished returns true if all the player's cards are face up.
func (p Player) IsFinished() bool {
	// TODO: optimize this with better internal state tracking
	for _, c := range p.board {
		if !c.FaceUp {
			return false
		}
	}

	return true
}
