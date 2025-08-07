package strategies

import (
	"testing"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/stretchr/testify/assert"
)

func TestSmartPlacesMatches(t *testing.T) {
	cases := []struct {
		name string

		// -1 will indicate face down for a shortcut
		board [8]int

		consideredCard     playnine.Card
		expectedMatchIndex int
	}{
		{
			name: "match bottom row 12",
			board: [8]int{
				7, 12, -1, -1,
				-1, -1, -1, 5,
			},
			consideredCard:     12,
			expectedMatchIndex: 5,
		},
		{
			name: "match top row 5",
			board: [8]int{
				7, 12, -1, -1,
				-1, -1, -1, 5,
			},
			consideredCard:     5,
			expectedMatchIndex: 3,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := SmartConfig{}

			board := playnine.PlayerBoard{}

			for i, card := range c.board {
				board[i] = playnine.PlayerBoardCard{
					Card:   playnine.Card(card),
					FaceUp: card != -1,
				}
			}

			result := cfg.tryMatch(playnine.PlayerBoard(board), playnine.Card(c.consideredCard))

			assert.Equal(t, c.expectedMatchIndex, result)
		})
	}
}
