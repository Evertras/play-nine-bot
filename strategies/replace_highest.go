package strategies

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

var ReplaceHigherDrawOrUseDiscard playnine.PlayerStrategyTakeTurnDrawOrUseDiscard = func(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
	state := g.CurrentPlayerState()
	availDiscard := g.AvailableDiscard()

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

var ReplaceHigherDrawn playnine.PlayerStrategyTakeTurnDrawn = func(g playnine.Game, drawnCard playnine.Card) (playnine.DecisionDrawn, playnine.DecisionCardIndex, error) {
	state := g.CurrentPlayerState()

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
