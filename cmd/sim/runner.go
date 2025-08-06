package main

import (
	"fmt"

	"github.com/evertras/play-nine-bot/playnine"
)

type longtermResult struct {
	player playnine.Player

	// placementCounts is 0 indexed so that [0] is 1st place finishes, etc
	placementCounts []int

	// avgPlacement is the 1-indexed result of where this player usually finished
	avgPlacement float32

	// avgScore is the average final score the player had
	avgScore float32
}

// runMany runs a given number of simulations with the players and
// returns the number of placement finishes for each
func runMany(players []playnine.Player, numRounds int) ([]longtermResult, error) {
	numPlayers := len(players)

	type totalResults struct {
		placementCounts []int
		totalScore      int
		totalPlacement  int
	}

	// id -> result
	longtermResults := make(map[int]totalResults)

	for _, p := range players {
		longtermResults[p.ID()] = totalResults{
			placementCounts: make([]int, numPlayers),
		}
	}

	type gameResult struct {
		result []result
		err    error
	}

	resultChan := make(chan gameResult)

	for range numRounds {
		go func(players []playnine.Player) {
			result, err := runGame(players)

			resultChan <- gameResult{
				result: result,
				err:    err,
			}
		}(players)
	}

	for range numRounds {
		gameResult := <-resultChan

		// Technically leaks the rest of the goroutines but this is a one shot run anyway
		// so they'll be destroyed by the process quitting soon after
		if gameResult.err != nil {
			return nil, fmt.Errorf("failed to run game: %w", gameResult.err)
		}

		for iPlace, r := range gameResult.result {
			id := r.player.ID()

			result := longtermResults[id]
			result.placementCounts[iPlace] += 1
			result.totalScore += r.finalScore
			result.totalPlacement += iPlace + 1

			longtermResults[id] = result
		}
	}

	returnedResults := make([]longtermResult, numPlayers)

	for iPlayer, p := range players {
		playerResult := longtermResults[p.ID()]

		returnedResults[iPlayer] = longtermResult{
			player:          p,
			placementCounts: playerResult.placementCounts,
			avgPlacement:    float32(playerResult.totalPlacement) / float32(numRounds),
			avgScore:        float32(playerResult.totalScore) / float32(numRounds),
		}
	}

	return returnedResults, nil
}
