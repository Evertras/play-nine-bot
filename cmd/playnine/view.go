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

	for _, state := range m.game.PlayerStates() {
		newline()
		b.WriteString(stringifyBoardState(state.CurrentBoard()))
		newline()
		fmt.Fprintf(&b, "Visible score: %d", state.CurrentBoard().ScoreVisible())
		newline()
	}

	return b.String()
}
