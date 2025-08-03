package playnine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreCalculations(t *testing.T) {
	testCases := []struct {
		name       string
		cardValues []int
		score      int
	}{
		{
			name: "all different",
			cardValues: []int{
				1, 2, 3, 4,
				5, 6, 7, 8,
			},
			score: 36,
		},
		{
			name: "one match",
			cardValues: []int{
				1, 2, 3, 4,
				1, 6, 7, 8,
			},
			score: 30,
		},
		{
			name: "two different matches",
			cardValues: []int{
				1, 2, 3, 4,
				1, 6, 7, 4,
			},
			score: 18,
		},
		{
			name: "two identical matches for negative score",
			cardValues: []int{
				1, 2, 4, 4,
				2, 1, 4, 4,
			},
			score: -4,
		},
		{
			name: "matched hole in ones don't negate them",
			cardValues: []int{
				1, 2, -5, 4,
				2, 1, -5, 5,
			},
			score: 5,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			board := PlayerBoard{}

			for i, val := range c.cardValues {
				board[i].Card = Card(val)
			}

			calculatedScore := board.scoreFinal()

			assert.Equal(t, c.score, calculatedScore)
		})
	}
}
