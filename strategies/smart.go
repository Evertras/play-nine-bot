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

// returns -1 if no index found, otherwise returns the index to replace
func smartTryMatch(board playnine.PlayerBoard, consideredCard playnine.Card) int {
	for i, card := range board {
		if !card.FaceUp {
			continue
		}

		if card.Card != consideredCard {
			continue
		}

		var iMatchingCard int

		if i < rowSize {
			iMatchingCard = i + rowSize
		} else {
			iMatchingCard = i - rowSize
		}

		// Skip if already a match
		if board[iMatchingCard].FaceUp && board[iMatchingCard].Card == consideredCard {
			continue
		}

		return iMatchingCard
	}

	return -1
}

// returns -1 if no inex found, otherwise returns the index to replace
func smartTryReplaceHighest(board playnine.PlayerBoard, consideredCard playnine.Card) int {
	for i, card := range board {
		if !card.FaceUp {
			continue
		}

		if card.Card > consideredCard {
			return i
		}
	}

	return -1
}

func SmartDrawOrUseDiscard(cfg SmartConfig) playnine.PlayerStrategyTakeTurnDrawOrUseDiscard {
	return func(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
		state := g.CurrentPlayerState()
		board := state.CurrentBoard()
		availDiscard := g.AvailableDiscard()

		// Check if we can complete any matches, whether the other card is visible or not
		if !cfg.ignoreMatches {
			iMatching := smartTryMatch(board, availDiscard)

			if iMatching >= 0 {
				return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iMatching), nil
			}
		}

		// Check if we can replace any of our face up cards with a lower discard
		iReplace := smartTryReplaceHighest(board, availDiscard)

		if iReplace >= 0 {
			return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iReplace), nil
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
			iMatching := smartTryMatch(board, drawnCard)

			if iMatching >= 0 {
				return playnine.DecisionDrawnReplaceCard, playnine.DecisionCardIndex(iMatching), nil
			}
		}

		// Check if we can replace any of our face up cards with the card we drew, ignore matches
		iReplace := smartTryReplaceHighest(board, drawnCard)
		if iReplace >= 0 {
			return playnine.DecisionDrawnReplaceCard, playnine.DecisionCardIndex(iReplace), nil
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
