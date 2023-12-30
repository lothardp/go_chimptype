package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	MAX_TEST_WIDTH = 60
	MIN_TEST_WIDTH = 10
)

const (
	cursorBackgroundColor = lipgloss.Color("250")
	cursorForegroundColor = lipgloss.Color("0")
	untypedCharColor      = lipgloss.Color("250")
	rightCharColor        = lipgloss.Color("250")
	wrongCharColor        = lipgloss.Color("1")
	extraCharColor        = lipgloss.Color("1")
)

func (c *Model) viewWelcome() string {
	title := lipgloss.NewStyle().Render("Welcome to Chimptype!\n")
	message := lipgloss.NewStyle().Render("Press enter key to start the test\nPress ESC or ctrl-c to quit")

	view := lipgloss.JoinVertical(lipgloss.Center, title, message)

	return lipgloss.Place(c.window.width, c.window.height, lipgloss.Center, lipgloss.Center, view)
}

func (c *Model) viewTestRunning(testState TestState) string {
	s, rawS := getTestString(testState)

	s = wrap(rawS, s, getTestWidth(c.window.width))

	test := lipgloss.NewStyle().Render(s)

	return lipgloss.Place(c.window.width, c.window.height, lipgloss.Center, lipgloss.Center, test)
}

func (c *Model) viewTestFinished(testResult TestResult) string {
	return "Well done!"
}

// Returns the string to be displayed and the raw string
func getTestString(testState TestState) (view string, rawString string) {
	words := translateWordList(testState.wordList)
	typedWords := translateTypedChars(testState.typedChars)

	wordIndex := 0
	for wordIndex < len(words) {
		word := words[wordIndex]

		typedWord := ""
		if wordIndex < len(typedWords) {
			typedWord = typedWords[wordIndex]
		}

		isCurrentWord := false
		if wordIndex == len(typedWords)-1 {
			isCurrentWord = true
		}

		w, rw := getWordViewString(word, typedWord, isCurrentWord)
		view += w
		rawString += rw

		if wordIndex < len(words)-1 {
			if isCurrentWord && len(typedWord) >= len(word) {
				view += renderCursor(" ")
			} else {
				view += renderUntypedChar(" ")
			}
			rawString += " "
		}

		wordIndex++
	}

	return view, rawString
}

// Returns the string to be displayed for a single word and the raw string
func getWordViewString(word string, typedWord string, isCurrentWord bool) (view string, rawWord string) {
	charIndex := 0

loop:
	for {
		var char, typedChar string

		if charIndex < len(word) {
			char = string(word[charIndex])
		}

		if charIndex < len(typedWord) {
			typedChar = string(typedWord[charIndex])
		}

		switch {
		// Word ended
		case char == "" && typedChar == "":
			break loop

		// Still chars left to type
		case char != "" && typedChar == "":
			if isCurrentWord && charIndex == len(typedWord) {
				view += renderCursor(char)
			} else {
				view += renderUntypedChar(char)
			}
			rawWord += char

		// Extra char typed
		case char == "" && typedChar != "":
			view += renderExtraChar(typedChar)
			rawWord += typedChar

		// Right char typed
		case char != "" && typedChar != "" && char == typedChar:
			view += renderRightChar(typedChar)
			rawWord += typedChar

		// Wrong char typed
		case char != "" && typedChar != "" && char != typedChar:
			view += renderWrongChar(typedChar)
			rawWord += typedChar

		default:
			panic("Unreachable!")
		}

		charIndex++
	}

	return view, rawWord
}

func renderExtraChar(char string) string {
	return lipgloss.NewStyle().
		Foreground(extraCharColor).
		Render(char)
}

func renderWrongChar(char string) string {
	return lipgloss.NewStyle().
		Foreground(wrongCharColor).
		Render(char)
}

func renderRightChar(char string) string {
	return lipgloss.NewStyle().
		Foreground(rightCharColor).
		Render(char)
}

func renderUntypedChar(char string) string {
	return lipgloss.NewStyle().
		Foreground(untypedCharColor).
		Render(char)
}

func renderCursor(char string) string {
	return lipgloss.NewStyle().
		Background(cursorBackgroundColor).
		Foreground(cursorForegroundColor).
		Render(char)
}

func translateWordList(wordList [][]Key) []string {
	words := []string{}
	for _, keys := range wordList {
		word := ""
		for _, key := range keys {
			word += string(key.char)
		}
		words = append(words, word)
	}

	return words
}

func translateTypedChars(typedChars []Key) []string {
	words := []string{""}
	for _, key := range typedChars {
		if key.keyType == Space {
			words = append(words, "")
			continue
		}
		words[len(words)-1] += string(key.char)
	}

	return words
}

func getTestWidth(w int) int {
	switch {
	case w < MIN_TEST_WIDTH:
		return MIN_TEST_WIDTH
	case w > MAX_TEST_WIDTH:
		return MAX_TEST_WIDTH
	}
	return w
}

// Wraps a string to a given width, and returns the wrapped string but with
// ANSI escape codes
// TODO: This is a bit hacky, maybe there is a better way to do this, i couldnt
// wrap the string with colors
func wrap(rawString string, viewString string, testWidth int) string {
	words := strings.Split(rawString, " ")
	wrappedRawString := []string{""}

	for _, word := range words {
		currentLineIndex := len(wrappedRawString) - 1
		currentLine := wrappedRawString[currentLineIndex]

		if len(currentLine)+len(word) <= testWidth {
			wrappedRawString[currentLineIndex] += word + " "
		} else {
			wrappedRawString[currentLineIndex] = strings.TrimRight(currentLine, " ")
			wrappedRawString = append(wrappedRawString, word+" ")
		}
	}

	viewWords := strings.Split(viewString, " ")

	finalViewString := ""
	i := 0
	for _, line := range wrappedRawString {
		nWords := strings.Split(line, " ")
		for range nWords {
			if i < len(viewWords) {
				finalViewString += viewWords[i] + " "
				i++
			}
		}
		finalViewString = strings.TrimRight(finalViewString, " ")
		finalViewString += "\n"
	}

	return finalViewString
}