package main

import (
	"fmt"
	"strings"

	"github.com/evertras/play-nine-bot/playnine"
)

func stringifyBoardState(state playnine.PlayerBoard) string {
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

func (m model) View() string {
	var b strings.Builder

	newline := func() {
		b.WriteRune('\n')
	}

	fmt.Fprintf(&b, "Discard: %2d", m.game.AvailableDiscard())

	for i, state := range m.game.PlayerStates() {
		newline()
		newline()
		curBoard := state.CurrentBoard()
		if m.game.CurrentPlayerIndex() == i {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}
		fmt.Fprintf(&b, "Player #%d", i+1)
		newline()
		b.WriteString(stringifyBoardState(curBoard))
		newline()
		fmt.Fprintf(&b, "Score: %d", curBoard.ScoreVisible())
	}

	return b.String()
}
