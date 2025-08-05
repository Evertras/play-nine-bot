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
			b.WriteString(fmt.Sprintf("%2d ", c.Card))
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

	for _, state := range m.game.PlayerStates() {
		b.WriteRune('\n')
		b.WriteString(stringifyBoardState(state.CurrentBoard()))
		b.WriteRune('\n')
	}

	return b.String()
}
