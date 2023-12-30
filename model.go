package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

// Holds the state of the program following the ELM architecture of bubbletea
type Model struct {
	state ModelState
}

type ModelState interface{}

// States of the whole program
type WelcomeState struct{}

type TestRunningState struct {
	testState TestState
}

type TestFinishedState struct {
	testResult TestResult
}

func (c Model) Init() tea.Cmd {
	return nil
}

// Returns the string to be displayed based on the current state
func (c Model) View() string {
	switch state := (c.state).(type) {
	case WelcomeState:
		return c.viewWelcome()

	case TestRunningState:
		return c.viewTestRunning(state.testState)

	case TestFinishedState:
		return c.viewTestFinished(state.testResult)

	default:
		panic("unknown state")
	}
}

// Handles messages and updates the state
func (c Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch (c.state).(type) {
	case WelcomeState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				c.startNewTest()
				return c, nil

			case tea.KeyCtrlC, tea.KeyEsc:
				return c, tea.Quit
			}
		}

	case TestRunningState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				c.interruptTest()
				return c, nil

			default:
				c.passMsgToTestState(msg)
				return c, nil
			}
		}

	case TestFinishedState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyEsc, tea.KeyCtrlC:
				c.state = WelcomeState{}
				return c, nil
			}
		}

	default:
		panic("unknown state")
	}

	return c, nil
}

// Translates tea.KeyMsg to Key for the internal test_state
func translateKey(msg tea.KeyMsg) Key {
	var key Key

	switch msg.Type {
	case tea.KeyBackspace:
		key.char = 127
		key.keyType = Backspace

	case tea.KeySpace:
		key.char = 32
		key.keyType = Space

	default:
		// TODO: Maybe only accept chars/numbers/punctuation?
		if len(msg.Runes) == 0 {
			// TODO: handle this better
			key.char = 'X'
		} else {
			key.char = msg.Runes[0] //TODO: check if this could be a problem
		}
		key.keyType = Char
	}

	return key
}

func (c *Model) startNewTest() {
	wordList := c.generateWordList(25)
	c.state = TestRunningState{
		testState: NewTestState(wordList),
	}
}

func (c *Model) interruptTest() {
	// TODO: do something else?
	// maybe save the half test result?
	c.state = WelcomeState{}
}

func (c *Model) passMsgToTestState(msg tea.KeyMsg) {
	if state, ok := c.state.(TestRunningState); ok {
		key := translateKey(msg)
		state.testState.HandleKey(key)
		c.state = state

		if state.testState.finished {
			c.state = TestFinishedState{
				testResult: TestResult{state.testState},
			}
		}
	} else {
		panic("test is not running")
	}
}

func (c *Model) viewWelcome() string {
	return "Welcome!"
}

func (c *Model) viewTestRunning(testState TestState) string {
	s := ""
	for _, word := range testState.wordList {
		for _, char := range word {
			s += fmt.Sprintf("%c", char.char)
		}
		s += fmt.Sprintf(" ")
	}
	s += fmt.Sprintln()

	for _, char := range testState.typedChars {
		if char.keyType == Space {
			s += fmt.Sprintf(" ")
		} else {
			s += fmt.Sprintf("%c", char.char)
		}
	}
	s += fmt.Sprintln()
	return s
}

func (c *Model) viewTestFinished(testResult TestResult) string {
	return "Well done!"
}

func (c *Model) generateWordList(wordCount int) []string {
	return []string{"no", "yes", "fish", "tree", "road",
		"music", "stone", "bird", "book", "light", "glass",
		"flower", "table", "phone", "house", "fish", "house",
		"phone", "music", "stone", "tree", "river", "green",
		"flower", "glass"}
}
