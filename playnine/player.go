package playnine

import (
	"fmt"
	"sync/atomic"
)

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

// Player is a player that has some strategy to play.
type Player struct {
	name string
	id   int

	strategyOpeningFlips             PlayerStrategyOpeningFlips
	strategyTakeTurnDrawOrUseDiscard PlayerStrategyTakeTurnDrawOrUseDiscard
	strategyTakeTurnDrawn            PlayerStrategyTakeTurnDrawn
}

var idCounter atomic.Int32

// PlayerState contains information for a single player in an active game.
type PlayerState struct {
	board PlayerBoard

	player *Player
}

// NewPlayer creates a new player with the given strategies that can be used to
// play the game.
func NewPlayer(
	name string,
	strategyOpeningFlips PlayerStrategyOpeningFlips,
	strategyTakeTurnDrawOrUseDiscard PlayerStrategyTakeTurnDrawOrUseDiscard,
	strategyTakeTurnDrawn PlayerStrategyTakeTurnDrawn,
) Player {
	return Player{
		name: name,
		id:   int(idCounter.Add(1)),

		strategyOpeningFlips:             strategyOpeningFlips,
		strategyTakeTurnDrawOrUseDiscard: strategyTakeTurnDrawOrUseDiscard,
		strategyTakeTurnDrawn:            strategyTakeTurnDrawn,
	}
}

// Name gets the name of the player.
func (p Player) Name() string {
	return p.name
}

// ID returns an atomically incremented ID so that two
// players with the same name can be differentiated.
func (p Player) ID() int {
	return p.id
}

// CurrentBoard returns the player's current board state.
func (p PlayerState) CurrentBoard() PlayerBoard {
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

// startGame draws a starting hand for the player and initializes
// a board, then executes the opening flip strategy to flip exactly two cards.
//
// Returns an error if the deck cannot be drawn from or the opening flip
// strategy doesn't flip exactly two unique cards.
func (p *Player) startGame(d *Deck) (PlayerState, error) {
	if d == nil {
		return PlayerState{}, fmt.Errorf("given deck was nil")
	}

	board := [PlayerBoardSize]PlayerBoardCard{}

	for i := range PlayerBoardSize {
		c, err := d.draw()

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

// ScoreVisible returns the visible score for face up cards.
func (p PlayerBoard) ScoreVisible() int {
	const halfBoard = PlayerBoardSize / 2
	total := 0

	// card value -> how many matched pairs
	matchCount := make(map[int]int)

	for i := range halfBoard {
		upperCard := p[i]
		lowerCard := p[i+halfBoard]

		upperCardValue := int(upperCard.Card)
		lowerCardValue := int(lowerCard.Card)

		if upperCardValue == lowerCardValue && upperCard.FaceUp && lowerCard.FaceUp {
			matchCount[upperCardValue] += 1

			// Still count the hole in ones (-5) in the total below
			if upperCardValue > 0 {
				continue
			}
		}

		if upperCard.FaceUp {
			total += upperCardValue
		}

		if lowerCard.FaceUp {
			total += lowerCardValue
		}
	}

	for _, count := range matchCount {
		// only care about >1
		switch count {
		case 2:
			total -= 10
		case 3:
			total -= 15
		case 4:
			total -= 20
		}
	}

	return total
}

// scoreFinal returns the total player score with their current board. This
// includes face down cards and is only intended to be used to tally final
// scores for the round internally.
func (p PlayerBoard) scoreFinal() int {
	const halfBoard = PlayerBoardSize / 2
	total := 0

	// card value -> how many matched pairs
	matchCount := make(map[int]int)

	for i := range halfBoard {
		upperCard := int(p[i].Card)
		lowerCard := int(p[i+halfBoard].Card)

		if upperCard == lowerCard {
			matchCount[upperCard] += 1

			// Still count the hole in ones (-5) in the total below
			if upperCard > 0 {
				continue
			}
		}

		total += upperCard + lowerCard
	}

	for _, count := range matchCount {
		// only care about >1
		switch count {
		case 2:
			total -= 10
		case 3:
			total -= 15
		case 4:
			total -= 20
		}
	}

	return total
}
