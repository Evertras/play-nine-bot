package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
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

func initialModel() model {
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

	return model{
		game: game,
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
