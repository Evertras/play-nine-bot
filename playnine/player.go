package playnine

import "fmt"

const PlayerBoardSize = 8

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

// CurrentBoard returns the player's current board state
func (p Player) CurrentBoard() [PlayerBoardSize]PlayerBoardCard {
	return p.board
}
