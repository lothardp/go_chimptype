package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"fmt"
	"os"
)

func main() {
	model := Model{
		state: WelcomeState{},
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Chimptype: %s\n", err)
		os.Exit(1)
	}
}
