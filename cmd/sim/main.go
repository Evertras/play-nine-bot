package main

import (
	"fmt"
	"os"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/evertras/play-nine-bot/strategies"
)

type result struct {
	player     playnine.Player
	finalScore int
}

func main() {
	const (
		numPlayers = 4
		numRounds  = 10000
	)

	players := make([]playnine.Player, numPlayers)

	for i := range numPlayers {
		players[i] = playnine.NewPlayer(
			fmt.Sprintf("Fast #%d", i+1),
			strategies.OpeningFlipsOppositeCorners,
			strategies.FastestDrawOrUseDiscard,
			strategies.FastestDrawn,
		)
	}

	finalScores, err := runMany(players, numRounds)

	if err != nil {
		fmt.Println("Failed to run game:", err)
		os.Exit(1)
	}

	for _, finalScore := range finalScores {
		fmt.Printf("%s - %.1f - %.1f", finalScore.player.Name(), finalScore.avgScore, finalScore.avgPlacement)
		fmt.Println()
	}
}
