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
		b.WriteString(game.Players()[i].Name())
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
	players := game.Players()
	numPlayers := len(players)
	roundScores := game.PlayerRoundScores()
	numRounds := len(roundScores)

	// Include a header row
	numRows := numPlayers + 1

	// First column is names, last column is total scores
	numCols := numRounds + 2

	// Some more explicit indices to use
	const (
		iHeaderRow = 0
		iNameCol   = 0
	)
	iTotalScoreCol := numCols - 1

	rows := make([][]string, numRows)

	// Header text is the number of the round
	rows[iHeaderRow] = make([]string, numCols)
	rows[iHeaderRow][iTotalScoreCol] = "Total"

	for i := range numPlayers {
		rows[i+1] = make([]string, numCols)
		rows[i+1][iNameCol] = players[i].Name()
	}

	scoreTotals := make([]int, numPlayers)

	for iRound, playerScores := range roundScores {
		rows[iHeaderRow][iRound+1] = fmt.Sprintf("%3d", iRound+1)

		for iPlayer, score := range playerScores {
			rows[iPlayer+1][iRound+1] = fmt.Sprintf("%d", score)
			scoreTotals[iPlayer] += score
		}

	}

	for iPlayer, total := range scoreTotals {
		rows[iPlayer+1][iTotalScoreCol] = fmt.Sprintf("%d", total)
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
