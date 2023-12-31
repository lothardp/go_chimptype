package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	model := Model{
		state: WelcomeState{numberOfWords: NUMBER_OF_WORDS_OPTIONS[0]},
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Chimptype: %s\n", err)
		os.Exit(1)
	}
}
