package main

import (
	"lothardp/go_chimptype/words"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

var NUMBER_OF_WORDS_OPTIONS = [5]int{5, 10, 25, 50, 100}

// Holds the state of the program following the ELM architecture of bubbletea
type Model struct {
	state  ModelState
	window Window
}

type Window struct {
	height int
	width  int
}

type ModelState interface{}

// States of the whole program
type WelcomeState struct {
	numberOfWords int
}

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
		return c.viewWelcome(state)

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
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		c.resizeWindow(msg)
		return c, nil
	}

	switch state := (c.state).(type) {
	case WelcomeState:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				c.startNewTest(state.numberOfWords)
				return c, nil

			case tea.KeyCtrlC, tea.KeyEsc:
				return c, tea.Quit

			case tea.KeyUp, tea.KeyDown, tea.KeyRunes:
				c.handleWelcomeStateKey(state, msg)
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
				c.state = WelcomeState{numberOfWords: NUMBER_OF_WORDS_OPTIONS[0]}
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
		// i think i should check for case tea.KeyRunes not default
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

func (c *Model) startNewTest(numWords int) {
	wordList := c.generateWordList(numWords)
	c.state = TestRunningState{
		testState: NewTestState(wordList),
	}
}

func (c *Model) interruptTest() {
	// TODO: do something else?
	// maybe save the half test result?
	c.state = WelcomeState{numberOfWords: NUMBER_OF_WORDS_OPTIONS[0]}
}

func (c *Model) passMsgToTestState(msg tea.KeyMsg) {
	if state, ok := c.state.(TestRunningState); ok {
		key := translateKey(msg)
		state.testState.HandleKey(key)
		c.state = state

		if state.testState.finished {
			c.state = TestFinishedState{
				testResult: newTestResult(state.testState),
			}
		}
	} else {
		panic("test is not running")
	}
}

func (c *Model) generateWordList(wordCount int) []string {
	list := make([]string, wordCount)

	for i := range list {
		list[i] = words.WORDS[rand.Intn(len(words.WORDS))]
	}

	return list
}

// Window resize msg is sent once at the start of the program
func (c *Model) resizeWindow(msg tea.WindowSizeMsg) {
	c.window.height = msg.Height
	c.window.width = msg.Width
}

func (c *Model) handleWelcomeStateKey(state WelcomeState, msg tea.KeyMsg) {
	upMove := msg.Type == tea.KeyUp || (msg.Type == tea.KeyRunes && len(msg.Runes) > 0 && msg.Runes[0] == 'k')
	downMove := msg.Type == tea.KeyDown || (msg.Type == tea.KeyRunes && len(msg.Runes) > 0 && msg.Runes[0] == 'j')

	currentIndex := findIndex(NUMBER_OF_WORDS_OPTIONS[:], state.numberOfWords)

	if upMove {
		if currentIndex > 0 {
			currentIndex--
		}
	} else if downMove {
		if currentIndex < len(NUMBER_OF_WORDS_OPTIONS)-1 {
			currentIndex++
		}
	}

	state.numberOfWords = NUMBER_OF_WORDS_OPTIONS[currentIndex]
	c.state = state
}

func findIndex(NUMBER_OF_WORDS_OPTIONS []int, i int) int {
	for index, value := range NUMBER_OF_WORDS_OPTIONS {
		if value == i {
			return index
		}
	}
	panic("Number of words not found in options")
}
