package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/play-nine-bot/playnine"
)

func renderPlayerBoards(game playnine.Game) string {

	stringifyBoardState := func(state playnine.PlayerBoard) string {
		var b strings.Builder

		for i, c := range state {
			if c.FaceUp {
				fmt.Fprintf(&b, "%2d ", c.Card)
			} else {
				b.WriteString(" - ")
			}

			if i == playnine.PlayerBoardSize/2-1 {
				b.WriteRune('\n')
			}
		}

		return b.String()
	}

	states := game.PlayerStates()

	boards := make([]string, len(states))

	for i, state := range states {
		var b strings.Builder

		style := lipgloss.NewStyle().PaddingLeft(1)

		if i != len(states)-1 {
			style = style.PaddingRight(1)
		}

		curBoard := state.CurrentBoard()

		if game.CurrentPlayerIndex() == i {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}
		fmt.Fprintf(&b, "Player #%d", i+1)
		b.WriteRune('\n')
		b.WriteString(stringifyBoardState(curBoard))
		b.WriteRune('\n')
		fmt.Fprintf(&b, "Score: %d", curBoard.ScoreVisible())

		boards[i] = style.Render(b.String())
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, boards...)
}

func (m model) View() string {
	var b strings.Builder

	playerBoards := renderPlayerBoards(m.game)

	border := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

	b.WriteString(border.Render(playerBoards))

	fmt.Fprintf(&b, "\n\nDiscard: %2d", m.game.AvailableDiscard())

	return b.String()
}
