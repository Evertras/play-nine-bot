package strategies

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

const rowSize = playnine.PlayerBoardSize / 2

type SmartConfig struct {
	// Defaults should be reasonably expected to be smarter
	ignoreMatches bool
}

func SmartDrawOrUseDiscard(cfg SmartConfig) playnine.PlayerStrategyTakeTurnDrawOrUseDiscard {
	return func(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
		state := g.CurrentPlayerState()
		board := state.CurrentBoard()
		availDiscard := g.AvailableDiscard()

		// Check if we can complete any matches, whether the other card is visible or not
		if !cfg.ignoreMatches {
			for i, card := range board {
				if !card.FaceUp {
					continue
				}

				if card.Card != availDiscard {
					continue
				}

				var iMatchingCard int

				if i < rowSize {
					iMatchingCard = i + rowSize
				} else {
					iMatchingCard = i - rowSize
				}

				// Skip if already a match
				if board[iMatchingCard].FaceUp && board[iMatchingCard].Card == availDiscard {
					continue
				}

				return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iMatchingCard), nil
			}
		}

		// Check if we can replace any of our face up cards with a lower discard, ignore matches
		for i, card := range state.CurrentBoard() {
			if !card.FaceUp {
				continue
			}

			if card.Card > availDiscard {
				return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(i), nil
			}
		}

		// If the discarded card isn't lower than anything we have, then try our luck with draw
		return playnine.DecisionDrawOrUseDiscardDraw, 0, nil
	}
}

func SmartDrawn(cfg SmartConfig) playnine.PlayerStrategyTakeTurnDrawn {
	return func(g playnine.Game, drawnCard playnine.Card) (playnine.DecisionDrawn, playnine.DecisionCardIndex, error) {
		state := g.CurrentPlayerState()
		board := state.CurrentBoard()

		// Check if we can complete any matches, whether the other card is visible or not
		if !cfg.ignoreMatches {
			for i, card := range board {
				if !card.FaceUp {
					continue
				}

				if card.Card != drawnCard {
					continue
				}

				var iMatchingCard int

				if i < rowSize {
					iMatchingCard = i + rowSize
				} else {
					iMatchingCard = i - rowSize
				}

				// Skip if already a match
				if board[iMatchingCard].FaceUp && board[iMatchingCard].Card == drawnCard {
					continue
				}

				return playnine.DecisionDrawnDiscardAndFlip, playnine.DecisionCardIndex(iMatchingCard), nil
			}
		}

		// Check if we can replace any of our face up cards with the card we drew, ignore matches
		for i, card := range state.CurrentBoard() {
			if !card.FaceUp {
				continue
			}

			if card.Card > drawnCard {
				return playnine.DecisionDrawnDiscardAndFlip, playnine.DecisionCardIndex(i), nil
			}
		}

		// Flip whatever the first face down card is
		for i, card := range state.CurrentBoard() {
			if !card.FaceUp {
				return playnine.DecisionDrawnDiscardAndFlip, playnine.DecisionCardIndex(i), nil
			}
		}

		return 0, 0, fmt.Errorf("didn't find a face down card to flip")
	}
}
