package strategies

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

type SmartConfig struct {
	// Defaults should be reasonably expected to be smarter

	// IgnoreMatches will ignore trying to complete matches.
	IgnoreMatches bool
}

// returns -1 if no index found, otherwise returns the index to replace
func (cfg SmartConfig) tryMatch(board playnine.PlayerBoard, consideredCard playnine.Card) int {
	for i, card := range board {
		if !card.FaceUp {
			continue
		}

		if card.Card != consideredCard {
			continue
		}

		iMatchingCard := matchIndex(i)

		// Skip if already a match
		if board[iMatchingCard].FaceUp && board[iMatchingCard].Card == consideredCard {
			continue
		}

		return iMatchingCard
	}

	return -1
}

// returns -1 if no index found, otherwise returns the index to replace
func (cfg SmartConfig) tryReplaceHighest(board playnine.PlayerBoard, consideredCard playnine.Card) int {
	highestCardWithoutMatch := consideredCard
	iHighestCardWithoutMatch := -1
	for i, card := range board {
		if !card.FaceUp {
			continue
		}

		// Skip if there's a match
		if !cfg.IgnoreMatches && hasVisibleMatchAtIndex(board, i) {
			continue
		}

		if card.Card > highestCardWithoutMatch {
			iHighestCardWithoutMatch = i
			highestCardWithoutMatch = card.Card
		}
	}

	return iHighestCardWithoutMatch
}

func SmartDrawOrUseDiscard(cfg SmartConfig) playnine.PlayerStrategyTakeTurnDrawOrUseDiscard {
	return func(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
		state := g.CurrentPlayerState()
		board := state.CurrentBoard()
		availDiscard := g.AvailableDiscard()

		// Check if we can complete any matches, whether the other card is visible or not
		if !cfg.IgnoreMatches {
			iMatching := cfg.tryMatch(board, availDiscard)

			if iMatching >= 0 {
				return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iMatching), nil
			}
		}

		// Check if we can replace any of our face up cards with a lower discard
		iReplace := cfg.tryReplaceHighest(board, availDiscard)

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
		if !cfg.IgnoreMatches {
			iMatching := cfg.tryMatch(board, drawnCard)

			if iMatching >= 0 {
				return playnine.DecisionDrawnReplaceCard, playnine.DecisionCardIndex(iMatching), nil
			}
		}

		// Check if we can replace any of our face up cards with the card we drew, ignore matches
		iReplace := cfg.tryReplaceHighest(board, drawnCard)
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
