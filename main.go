package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	model := Model{
		state: WelcomeState{numberOfWords: NUMBER_OF_WORDS_OPTIONS[0]},
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Chimptype: %s\n", err)
		os.Exit(1)
	}
}
