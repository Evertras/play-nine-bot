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

// PlayerBoard is represented by a 1D array, where 0-3 are the
// top row and 4-7 are the bottom row.
//
// TODO: Make this clearer/easier for strategies to work with, very leaky as-is...
type PlayerBoard [PlayerBoardSize]PlayerBoardCard

// PlayerStrategyOpeningFlips must return two unique indices of the first two
// cards to flip.
type PlayerStrategyOpeningFlips func() [2]int

// PlayerStrategyTakeTurn takes in the state of the game and takes a turn,
// mutating the underlying game state (for now).
type PlayerStrategyTakeTurn func(d *Deck)

// Player is a player that has some strategy to play.
type Player struct {
	strategyOpeningFlips PlayerStrategyOpeningFlips
	strategyTakeTurn     PlayerStrategyTakeTurn
}

// PlayerState contains information for a single player in an active game.
type PlayerState struct {
	board PlayerBoard

	player *Player
}

// NewPlayer creates a new player with the given strategies that can be used to
// play the game.
func NewPlayer(strategyOpeningFlips PlayerStrategyOpeningFlips, strategyTakeTurn PlayerStrategyTakeTurn) Player {
	return Player{
		strategyOpeningFlips: strategyOpeningFlips,
		strategyTakeTurn:     strategyTakeTurn,
	}
}

// StartGame draws a starting hand for the player and initializes
// a board, then executes the opening flip strategy to flip exactly two cards.
//
// Returns an error if the deck cannot be drawn from or the opening flip
// strategy doesn't flip exactly two unique cards.
func (p *Player) StartGame(d *Deck) (PlayerState, error) {
	if d == nil {
		return PlayerState{}, fmt.Errorf("given deck was nil")
	}

	board := [PlayerBoardSize]PlayerBoardCard{}

	for i := range PlayerBoardSize {
		c, err := d.Draw()

		if err != nil {
			return PlayerState{}, fmt.Errorf("tried to draw starting board for player: %w", err)
		}

		board[i] = PlayerBoardCard{
			Card: c,
		}
	}

	flipIndices := p.strategyOpeningFlips()

	if flipIndices[0] == flipIndices[1] {
		return PlayerState{}, fmt.Errorf("opening flip strategy failed to produce unique card indices to flip, both returned values were %d", flipIndices[0])
	}

	for _, i := range flipIndices {
		if i < 0 || i >= PlayerBoardSize {
			return PlayerState{}, fmt.Errorf("index to flip was out of range, should be between 0 and 7 but got %d", i)
		}

		board[i].FaceUp = true
	}

	return PlayerState{
		board:  board,
		player: p,
	}, nil
}

// CurrentBoard returns the player's current board state.
func (p PlayerState) CurrentBoard() [PlayerBoardSize]PlayerBoardCard {
	return p.board
}

// IsFinished returns true if all the player's cards are face up.
func (p PlayerState) IsFinished() bool {
	// TODO: optimize this with better internal state tracking
	for _, c := range p.board {
		if !c.FaceUp {
			return false
		}
	}

	return true
}

// scoreFinal returns the total player score with their current board. This
// includes face down cards and is only intended to be used to tally final
// scores for the round internally.
func (p PlayerState) scoreFinal() int {
	return 0
}
