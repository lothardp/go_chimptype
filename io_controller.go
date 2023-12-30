package main

import (
	"fmt"
	"os"
)

type IOController struct {
	// Fields
}

func NewIOController() IOController {
	return IOController{
		// Initialization
	}
}

func (io *IOController) ReadOneChar() Key {
	var key Key
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)
	key.char = rune(b[0])
	key.keyType = Char

	if key.char == 32 {
		key.keyType = Space
	} else  if key.char == 13 {
		key.keyType = Enter
	} else if key.char == 127 {
		key.keyType = Backspace
	} else if key.char == 27 {
		key.keyType = Esc
	}

	return key
}

func (io *IOController) ClearScreenRaw() {
	fmt.Println("\r")
}

func (io *IOController) ClearScreen() {
	// Clear screen logic
	fmt.Println("\r")
}

func (io *IOController) Flush() {
	// Flush logic
}

func (c *IOController) drawWelcome() {
	fmt.Println("Welcome")
}

func (c *IOController) drawTestRunning(testState TestState) {
	for _, word := range testState.wordList {
		for _, char := range word {
			fmt.Printf("%c", char.char)
		}
		fmt.Printf(" ")
	}
	fmt.Println("\r")

	for _, char := range testState.typedChars {
		if char.keyType == Space {
			fmt.Printf(" ")
		} else {
			fmt.Printf("%c", char.char)
		}
	}
	fmt.Println()
}

func (c *IOController) drawTestFinished(testResult TestResult) {
	fmt.Println("Test finished")
}
