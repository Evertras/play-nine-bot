package strategies

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

var FastestDrawOrUseDiscard playnine.PlayerStrategyTakeTurnDrawOrUseDiscard = func(g playnine.Game) (playnine.DecisionDrawOrUseDiscard, playnine.DecisionCardIndex, error) {
	return playnine.DecisionDrawOrUseDiscardDraw, 0, nil
}

var FastestDrawn playnine.PlayerStrategyTakeTurnDrawn = func(g playnine.Game, _ playnine.Card) (playnine.DecisionDrawn, playnine.DecisionCardIndex, error) {
	// The strategy is dumb and simple: ignore everything and just flip cards
	state := g.CurrentPlayerState()

	for index, card := range state.CurrentBoard() {
		if !card.FaceUp {
			return playnine.DecisionDrawnDiscardAndFlip, playnine.DecisionCardIndex(index), nil
		}
	}

	return 0, 0, fmt.Errorf("didn't find a face down card to flip")
}
