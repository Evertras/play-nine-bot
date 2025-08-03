package strategies

import "github.com/evertras/play-nine-bot/playnine"

// OpeningFlipsFirstVertical flips the first two vertical cards.
var OpeningFlipsFirstVertical playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, 4}
}

// OpeningFlipsOppositeCorners flips the top left and bottom right cards.
var OpeningFlipsOppositeCorners playnine.PlayerStrategyOpeningFlips = func() [2]int {
	return [2]int{0, playnine.PlayerBoardSize - 1}
}
