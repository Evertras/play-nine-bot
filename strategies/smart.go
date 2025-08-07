package strategies

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

type SmartConfig struct {
	// Defaults should be reasonably expected to be smarter

	// IgnoreMatches will ignore trying to complete matches.
	//
	// This is a terrible idea.
	IgnoreMatches bool

	// FlipFirstVertical will flip the first vertical pair on the opener
	// rather than the opposite corners.
	//
	// In testing, this is a small but noticeable disadvantage of a few points per round.
	FlipFirstVertical bool

	// ReplaceDiffThreshold is how big the difference must be before
	// replacing vs flipping is considered
	ReplaceDiffThreshold int
}

func (cfg SmartConfig) NewPlayer(name string) playnine.Player {
	return playnine.NewPlayer(name, cfg.OpeningFlips, cfg.DrawOrUseDiscard, cfg.Drawn)
}

func (cfg SmartConfig) OpeningFlips() [2]int {
	if cfg.FlipFirstVertical {
		return OpeningFlipsFirstVertical()
	}

	return OpeningFlipsOppositeCorners()
}

func (cfg SmartConfig) DrawOrUseDiscard(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
	state := g.CurrentPlayerState()
	board := state.CurrentBoard()
	availDiscard := g.AvailableDiscard()

	// Check if we can complete any matches, whether the other card is visible or not
	if !cfg.IgnoreMatches {
		if iMatching := cfg.tryMatch(board, availDiscard); iMatching >= 0 {
			return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iMatching), nil
		}
	}

	// Check if we can replace any of our face up cards with a lower discard
	if iReplace := cfg.tryReplaceHighest(board, availDiscard); iReplace >= 0 {
		return playnine.DecisionDrawOrUseDiscardUseDiscard, playnine.DecisionCardIndex(iReplace), nil
	}

	// If the discarded card isn't lower than anything we have, then try our luck with draw
	return playnine.DecisionDrawOrUseDiscardDraw, 0, nil
}

func (cfg SmartConfig) Drawn(g playnine.Game, drawnCard playnine.Card) (playnine.DecisionDrawn, playnine.DecisionCardIndex, error) {
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
	if iReplace := cfg.tryReplaceHighest(board, drawnCard); iReplace >= 0 {
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
	highestCardWithoutMatch := consideredCard + playnine.Card(cfg.ReplaceDiffThreshold-1)
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
