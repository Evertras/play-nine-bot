package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

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

	border := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder())

	return border.Render(lipgloss.JoinHorizontal(lipgloss.Bottom, boards...))
}

func renderPlayerScoreTable(game playnine.Game) string {
	roundScores := game.PlayerRoundScores()

	rows := make([][]string, len(roundScores))

	for round, playerScores := range roundScores {
		row := make([]string, len(playerScores)+1)

		row[0] = fmt.Sprintf("%d ", round+1)
		for i, score := range playerScores {
			row[i+1] = fmt.Sprintf("%d", score)
		}

		rows[round] = row
	}

	t := table.New().Border(lipgloss.NormalBorder()).Rows(rows...)

	return t.String()
}

func (m model) View() string {
	var b strings.Builder

	playerBoards := renderPlayerBoards(m.game)
	playerScores := renderPlayerScoreTable(m.game)

	b.WriteString(lipgloss.JoinVertical(lipgloss.Left, playerBoards, playerScores))

	fmt.Fprintf(&b, "\n\nDiscard: %2d", m.game.AvailableDiscard())

	return b.String()
}
