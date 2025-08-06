package strategies

import "github.com/evertras/play-nine-bot/playnine"

const rowSize = playnine.PlayerBoardSize / 2

// matchIndex returns the index of the card that would be
// considered a match for the given card index
func matchIndex(i int) int {
	if i < rowSize {
		return i + rowSize
	}

	return i - rowSize
}

// hasVisibleMatchAtIndex returns true if the card at the index
// has a matching card above/below it for a match
func hasVisibleMatchAtIndex(board playnine.PlayerBoard, i int) bool {
	card := board[i]
	iCheck := matchIndex(i)

	return board[iCheck].FaceUp && board[iCheck].Card == card.Card
}
