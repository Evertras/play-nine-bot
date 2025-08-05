package playnine

// PlayerStrategyOpeningFlips must return two unique indices of the first two
// cards to flip.
type PlayerStrategyOpeningFlips func() [2]int

/*
A player strategy must make the following decisions:

First, draw or take from discard?

If draw, then decide:

- Replace either a face down or face up card on their board
- Discard the drawn card and flip up a facedown card
- OR, if there is only one facedown card, optionally skip the turn

If take from discard:

- Replace either a face down or face up card on their board
*/

// DecisionDrawOrUseDiscard represents the first decision to be made,
// whether to draw a new card or use the available discarded card.
type DecisionDrawOrUseDiscard int8

// DecisionDrawn represents a decision after drawing a card.
type DecisionDrawn int8

const (
	DecisionDrawOrUseDiscardDraw DecisionDrawOrUseDiscard = iota
	DecisionDrawOrUseDiscardUseDiscard

	DecisionDrawnReplaceCard DecisionDrawn = iota
	DecisionDrawnDiscardAndFlip
	DecisionDrawnDiscardAndSkip
)

// CardIndex is the 0 based index of which card to affect.
type DecisionCardIndex uint8

// valid is a quick check of whether the index is valid to avoid panics later,
// so we can safeguard against strategies trying to do wild things.
func (i DecisionCardIndex) valid() bool {
	return i < PlayerBoardSize
}

// PlayerStrategyTakeTurn takes in the state of the game and returns what
// decision to make when it comes to drawing or discarding. If using discard, it
// must also return the card index of which card to replace/flip.
type PlayerStrategyTakeTurnDrawOrUseDiscard func(g Game) (DecisionDrawOrUseDiscard, DecisionCardIndex, error)

// PlayerStrategyTakeTurnDrawn takes in the state of the game and returns what
// decision to make when the prior decision was to draw a card. Must return the
// chosen card index if replacing or discarding/flipping.
type PlayerStrategyTakeTurnDrawn func(g Game) (DecisionDrawn, DecisionCardIndex, error)
