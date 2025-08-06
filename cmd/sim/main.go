package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/evertras/play-nine-bot/strategies"
)

type result struct {
	player     playnine.Player
	finalScore int
}

func main() {
	const numPlayers = 4

	players := make([]playnine.Player, numPlayers)

	for i := range numPlayers {
		players[i] = playnine.NewPlayer(
			fmt.Sprintf("Fast #%d", i+1),
			strategies.OpeningFlipsOppositeCorners,
			strategies.FastestDrawOrUseDiscard,
			strategies.FastestDrawn,
		)
	}

	finalScores, err := runGame(players)

	if err != nil {
		fmt.Println("Failed to run game:", err)
		os.Exit(1)
	}

	slices.SortFunc(finalScores, func(a, b result) int {
		if a.finalScore < b.finalScore {
			return -1
		}

		return 1
	})

	for _, finalScore := range finalScores {
		fmt.Println(finalScore.player.Name(), " ", finalScore.finalScore)
	}
}
