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
		numRounds  = 100000
	)

	players := make([]playnine.Player, 0, numPlayers)

	/*
		players = append(players, playnine.NewPlayer(
			"Fast",
			strategies.OpeningFlipsOppositeCorners,
			strategies.FastestDrawOrUseDiscard,
			strategies.FastestDrawn,
		))
	*/

	players = append(players,
		playnine.NewPlayer(
			"Replacer",
			strategies.OpeningFlipsOppositeCorners,
			strategies.ReplaceHigherDrawOrUseDiscard,
			strategies.ReplaceHigherDrawn,
		),
	)

	smartDefaultCfg := strategies.SmartConfig{}

	players = append(players, smartDefaultCfg.NewPlayer("Smart Default"))

	smartCfgLab := strategies.SmartConfig{
		ReplaceDiffThreshold: 2,
	}

	players = append(players, smartCfgLab.NewPlayer("Smart Lab"))

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
