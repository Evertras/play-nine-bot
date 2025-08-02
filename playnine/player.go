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

// Player contains information for a single player in the game
type Player struct {
	board [PlayerBoardSize]PlayerBoardCard
}

// NewPlayerFromDeck draws a starting hand for the player and initializes
// a board where all cards are still face down.
func NewPlayerFromDeck(d *Deck) (Player, error) {
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

	return Player{
		board,
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
