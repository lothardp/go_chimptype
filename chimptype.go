package main

import (
// "time"
)

type ChimpType struct {
	ioController IOController
	state        State
}

type TestResult struct {
	finalState TestState
}

type State interface {
	// Interface methods
}

type WelcomeState struct{}
type TestRunningState struct {
	testState TestState
}
type TestFinishedState struct {
	testResult TestResult
}
type ExitState struct{}

func NewChimpType() ChimpType {
	return ChimpType{
		ioController: NewIOController(),
		state:        WelcomeState{},
	}
}

func (c *ChimpType) Start() {
	c.ioController.ClearScreenRaw()
	c.Draw()
	c.MainLoop()
}

func (c *ChimpType) MainLoop() {
	for {
		key := c.ioController.ReadOneChar()
		c.HandleKey(key)
		c.Draw()

		if _, ok := c.state.(ExitState); ok {
			break
		}
	}
}

func (c *ChimpType) HandleKey(key Key) {
	println(key.char)
	println(key.keyType)
	switch state := (c.state).(type) {
	case WelcomeState:
		if key.keyType == Enter {
			wordList := c.generateWordList(25)
			c.state = TestRunningState{
				testState: NewTestState(wordList),
			}
		} else if key.keyType == Esc {
			c.state = ExitState{}
		}

	case TestRunningState:
		if key.keyType == Esc {
			c.state = WelcomeState{}
		} else {
			state.testState.HandleKey(key)
			c.state = state
			if state.testState.finished {
				c.state = TestFinishedState{
					testResult: TestResult{state.testState},
				}
			}
		}

	case TestFinishedState:
		c.state = WelcomeState{}
	}
}

func (c *ChimpType) Draw() {
	c.ioController.ClearScreen()

	switch state := (c.state).(type) {
	case WelcomeState:
		c.drawWelcome()

	case TestRunningState:
		c.drawTestRunning(state.testState)

	case TestFinishedState:
		c.drawTestFinished(state.testResult)
	}
}

func (c *ChimpType) drawWelcome() {
	c.ioController.drawWelcome()
}

func (c *ChimpType) drawTestRunning(testState TestState) {
	c.ioController.drawTestRunning(testState)

}

func (c *ChimpType) drawTestFinished(testResult TestResult) {
	c.ioController.drawTestFinished(testResult)
}

func (c *ChimpType) generateWordList(wordCount int) []string {
	return []string{"no", "yes", "fish", "tree", "road",
		"music", "stone", "bird", "book", "light", "glass",
		"flower", "table", "phone", "house", "fish", "house",
		"phone", "music", "stone", "tree", "river", "green",
		"flower", "glass"}
}
