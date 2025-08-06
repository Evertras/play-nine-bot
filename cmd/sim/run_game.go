package main

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

// runGame runs a full game and returns an ordered result set
// with the first entry being first place and so on.
func runGame(players []playnine.Player) ([]result, error) {
	numPlayers := len(players)

	g, err := playnine.NewGame(players)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	for !g.Finished() {
		err := g.TakeTurn()
		if err != nil {
			return nil, fmt.Errorf("failed to take turn: %w", err)
		}
	}

	finalScores := make([]result, numPlayers)

	for i := range numPlayers {
		finalScores[i].player = players[i]
	}

	for _, roundScores := range g.PlayerRoundScores() {
		for i, score := range roundScores {
			finalScores[i].finalScore += score
		}
	}

	return finalScores, nil
}
