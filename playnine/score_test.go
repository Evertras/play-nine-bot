package playnine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoreCalculations(t *testing.T) {
	testCases := []struct {
		name         string
		cardValues   []int
		cardFaceUp   []bool
		scoreFinal   int
		scoreVisible int
	}{
		{
			name: "all different",
			cardValues: []int{
				1, 2, 3, 4,
				5, 6, 7, 8,
			},
			cardFaceUp: []bool{
				true, true, true, true,
				false, false, false, false,
			},
			scoreFinal:   36,
			scoreVisible: 10,
		},
		{
			name: "one match visible",
			cardValues: []int{
				1, 2, 3, 4,
				1, 6, 7, 8,
			},
			cardFaceUp: []bool{
				true, false, false, true,
				true, false, false, true,
			},
			scoreFinal:   30,
			scoreVisible: 12,
		},
		{
			name: "one match hidden",
			cardValues: []int{
				1, 2, 3, 4,
				1, 6, 7, 8,
			},
			cardFaceUp: []bool{
				true, true, true, true,
				false, true, true, true,
			},
			scoreFinal:   30,
			scoreVisible: 31,
		},
		{
			name: "two different matches visible",
			cardValues: []int{
				1, 2, 3, 4,
				1, 6, 7, 4,
			},
			cardFaceUp: []bool{
				true, true, false, true,
				true, false, true, true,
			},
			scoreFinal:   18,
			scoreVisible: 9,
		},
		{
			name: "two identical matches for negative score",
			cardValues: []int{
				1, 2, 4, 4,
				2, 1, 4, 4,
			},
			cardFaceUp: []bool{
				true, true, true, true,
				false, false, false, true,
			},
			scoreFinal:   -4,
			scoreVisible: 7,
		},
		{
			name: "matched hole in ones don't negate them",
			cardValues: []int{
				1, 2, -5, 4,
				2, 1, -5, 5,
			},
			cardFaceUp: []bool{
				true, true, true, true,
				true, true, true, true,
			},
			scoreFinal:   5,
			scoreVisible: 5,
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			board := PlayerBoard{}

			for i, val := range c.cardValues {
				board[i].Card = Card(val)
			}

			for i, faceUp := range c.cardFaceUp {
				board[i].FaceUp = faceUp
			}

			calculatedScoreFinal := board.scoreFinal()
			assert.Equal(t, c.scoreFinal, calculatedScoreFinal)

			calculatedScoreVisible := board.ScoreVisible()
			assert.Equal(t, c.scoreVisible, calculatedScoreVisible)
		})
	}
}
