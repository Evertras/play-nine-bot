package playnine

import (
	"errors"
	"fmt"
)

// TotalRounds is how many rounds are played in total.
const TotalRounds = 9

var ErrGameAlreadyFinished = errors.New("game is already finished")

// Game represents the state of the game and can advance forward through
// each player's turn.
type Game struct {
	deck               *Deck
	discarded          Card
	players            []Player
	playerStates       []PlayerState
	currentPlayerIndex int

	round             int
	playerRoundScores [][]int
}

// NewGame creates a new game with the given players ready to play.
func NewGame(players []Player) (Game, error) {
	g := Game{
		players: players,

		round: 1,
	}

	err := g.dealFreshDeck()

	return g, err
}

// CurrentRound returns the current round the game is on,
// from 1-9 inclusive.
func (g Game) CurrentRound() int {
	return g.round
}

// CurrentPlayerIndex returns the current player's index so a player strategy
// can know who they are. 0-indexed.
func (g Game) CurrentPlayerIndex() int {
	return g.currentPlayerIndex
}

// PlayerStates gets the current round's player states.
func (g Game) PlayerStates() []PlayerState {
	return g.playerStates
}

// CurrentPlayerState returns the current player's state, as a shortcut for
// finding the current player's index.
func (g Game) CurrentPlayerState() PlayerState {
	return g.playerStates[g.CurrentPlayerIndex()]
}

// AvailableDiscard returns the card that's available for choosing to discard.
func (g Game) AvailableDiscard() Card {
	return g.discarded
}

// PlayerRoundScores returns the round-by-round scores of the players, filled
// in for each round completed. The first index is the round, and the second
// index is the player index.
func (g Game) PlayerRoundScores() [][]int {
	return g.playerRoundScores
}

// Finished returns true if the game is complete.
func (g Game) Finished() bool {
	return g.round == TotalRounds+1
}

// TakeTurn takes the turn for the current player and advances to the next player.
func (g *Game) TakeTurn() error {
	if g.Finished() {
		return ErrGameAlreadyFinished
	}

	// If the current player is done then that means we've finished the round,
	// because everyone else got to go once.
	if g.CurrentPlayerState().IsFinished() {
		g.advanceRound()

		return nil
	}

	curPlayer := g.players[g.CurrentPlayerIndex()]

	drawOrUseDiscard, discardIndex, err := curPlayer.strategyTakeTurnDrawOrUseDiscard(*g)

	if err != nil {
		return fmt.Errorf("failed to choose draw/discard: %w", err)
	}

	switch drawOrUseDiscard {
	case DecisionDrawOrUseDiscardDraw:
		// This is a bit nested, but so be it for now...
		decisionDrawn, replaceIndex, err := curPlayer.strategyTakeTurnDrawn(*g)

		if err != nil {
			return fmt.Errorf("failed to choose what to do when drawing: %w", err)
		}

		drawnCard, err := g.deck.draw()

		if err != nil {
			return fmt.Errorf("failed to draw card: %w", err)
		}

		switch decisionDrawn {
		case DecisionDrawnReplaceCard:
			if !replaceIndex.valid() {
				return fmt.Errorf("invalid index to replace: %v", replaceIndex)
			}

			oldCard := g.playerStates[g.currentPlayerIndex].board[replaceIndex].Card

			g.playerStates[g.currentPlayerIndex].board[replaceIndex] = PlayerBoardCard{
				Card:   drawnCard,
				FaceUp: true,
			}

			g.discarded = oldCard

		case DecisionDrawnDiscardAndFlip:
			if !replaceIndex.valid() {
				return fmt.Errorf("invalid index to replace: %v", replaceIndex)
			}

			g.playerStates[g.currentPlayerIndex].board[replaceIndex].FaceUp = true

			g.discarded = drawnCard

		case DecisionDrawnDiscardAndSkip:
			// Enforce correctness
			seenFaceDown := false
			for _, c := range g.playerStates[g.currentPlayerIndex].board {
				if !c.FaceUp {
					if seenFaceDown {
						return fmt.Errorf("can only skip if one card is left face down")
					}

					seenFaceDown = true
				}
			}

		default:
			return fmt.Errorf("unexpected decision for drawn card: %w", err)
		}

	case DecisionDrawOrUseDiscardUseDiscard:
		// TODO: do the thing
		return fmt.Errorf("will flip %v eventually but not implemented yet", discardIndex)

	default:
		return fmt.Errorf("unexpected decision to draw or use discard: %v", drawOrUseDiscard)
	}

	// Advance the turn
	g.currentPlayerIndex++

	if g.currentPlayerIndex >= len(g.players) {
		g.currentPlayerIndex = 0
	}

	return nil
}

func (g *Game) dealFreshDeck() error {
	g.deck = NewDeck()

	g.playerStates = make([]PlayerState, len(g.players))

	for i, p := range g.players {
		state, err := p.startGame(g.deck)

		if err != nil {
			return fmt.Errorf("failed to create new player state for index %d: %w", i, err)
		}

		g.playerStates[i] = state
	}

	var err error
	g.discarded, err = g.deck.draw()
	if err != nil {
		return fmt.Errorf("failed to draw for discard: %w", err)
	}

	return nil
}

func (g *Game) advanceRound() error {
	scores := make([]int, len(g.playerStates))

	for i, p := range g.playerStates {
		scores[i] = p.board.scoreFinal()
	}

	g.playerRoundScores = append(g.playerRoundScores, scores)

	// Do this before incrementing round because round is 1-index
	g.currentPlayerIndex = g.round % len(g.players)

	g.round++

	return g.dealFreshDeck()
}
