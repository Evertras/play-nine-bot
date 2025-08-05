package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/evertras/play-nine-bot/playnine"
	"github.com/evertras/play-nine-bot/strategies"
)

type model struct {
	game playnine.Game
}

func newGame() playnine.Game {
	makePlayer := func() playnine.Player {
		return playnine.NewPlayer(
			strategies.OpeningFlipsOppositeCorners,
			strategies.FastestDrawOrUseDiscard,
			strategies.FastestDrawn,
		)
	}

	players := []playnine.Player{}

	for range 4 {
		players = append(players, makePlayer())
	}

	game, err := playnine.NewGame(players)

	if err != nil {
		panic(err)
	}

	return game
}

func initialModel() model {
	return model{
		game: newGame(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			// This isn't really in the spirit of bubble tea's immutability, but close
			// enough for now... there's so many underlying slices to copy that it's
			// not worth messing with it atm.
			m.game.TakeTurn()

		case "r":
			// Restart the game
			m.game = newGame()
		}

	}

	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
