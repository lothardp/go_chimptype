package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	MAX_TEST_WIDTH = 60
	MIN_TEST_WIDTH = 10
)

const (
	defaultBackgroundColor = lipgloss.Color("94")
	defaultForegroundColor = lipgloss.Color("126")

	cursorBackgroundColor = lipgloss.Color("250")
	cursorForegroundColor = lipgloss.Color("0")
	untypedCharColor      = lipgloss.Color("250")
	rightCharColor        = lipgloss.Color("244")
	wrongCharColor        = lipgloss.Color("1")
	missedCharColor       = lipgloss.Color("1")
	extraCharColor        = lipgloss.Color("1")
)

func (c *Model) viewWelcome(state WelcomeState) string {
	title := lipgloss.NewStyle().Render("Welcome to Chimptype!\n")

	wordsMenu := viewNumberOfWordsMenu(state.numberOfWords)

	footer := lipgloss.NewStyle().Render("Press ESC or ctrl-c to quit")

	view := lipgloss.JoinVertical(lipgloss.Center, title, wordsMenu, footer)

	return lipgloss.Place(c.window.width, c.window.height, lipgloss.Center, lipgloss.Center, view)
}

func (c *Model) viewTestRunning(state TestRunningState) string {
	testState := state.testState

	durationString := durationView(state.duration.Seconds())

	s, _ := getTestString(testState)

	testStyle := lipgloss.NewStyle().Width(getTestWidth(c.window.width)).Align(lipgloss.Center)

	s = testStyle.Render(s)

	s = strings.ReplaceAll(s, "\n", "\x1b[40m\n")
	s = strings.ReplaceAll(s, "  ", " \x1b[40m ")

	// str := lipgloss.NewStyle().Render("h")
	// // Print the string as i
	// for _, char := range str {
	// 	fmt.Printf("(%x %c)", char, char)
	// }
	// fmt.Println()

	// s = wrap(rawS, s, getTestWidth(c.window.width))

	test := lipgloss.JoinVertical(lipgloss.Center, lipgloss.NewStyle().Render(durationString), s)

	test = lipgloss.Place(c.window.width, c.window.height, lipgloss.Center, lipgloss.Center, test)

	return test
}

func durationView(seconds float64) string {
	s := fmt.Sprintf("%.1fs\n", seconds)

	return lipgloss.NewStyle().Render(s)
}

func (c *Model) viewTestFinished(testResult TestResult) string {
	s := "Test finished!\n\nYour Results:\n\n"
	s += viewTestResults(testResult)

	return lipgloss.Place(c.window.width, c.window.height, lipgloss.Center, lipgloss.Center, s)
}

func viewTestResults(testResult TestResult) string {
	s := ""
	s += fmt.Sprintf("Net WPM: %.1f\n", testResult.netWPM)
	s += fmt.Sprintf("Raw WPM: %.1f\n", testResult.rawWPM)
	s += fmt.Sprintf("Accuracy: %.1f%%\n", testResult.accuracy)
	s += fmt.Sprintf("Time: %.1f seconds\n", testResult.duration.Seconds())
	s += fmt.Sprintf("Errors: %d\n", testResult.errors)
	s += fmt.Sprintf("Correct: %d\n", testResult.correct)
	s += fmt.Sprintf("Missed: %d\n", testResult.missed)
	s += fmt.Sprintf("Extra: %d\n", testResult.extra)

	return lipgloss.NewStyle().Render(s)
}

func viewNumberOfWordsMenu(numberOfWords int) string {
	s := "Select the number of words:\n"
	for _, n := range NUMBER_OF_WORDS_OPTIONS {
		selected := " "
		if n == numberOfWords {
			selected = "x"
		}
		s += fmt.Sprintf("[%s] %d words\n", selected, n)
	}

	return lipgloss.NewStyle().Render(s)
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

		wordStatus := "passed"
		if wordIndex == len(typedWords)-1 {
			wordStatus = "current"
		} else if wordIndex > len(typedWords)-1 {
			wordStatus = "next"
		}

		w, rw := getWordViewString(word, typedWord, wordStatus)
		view += w
		rawString += rw

		if wordIndex < len(words)-1 {
			isCurrentWord := wordStatus == "current"
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
func getWordViewString(word string, typedWord string, wordStatus string) (view string, rawWord string) {
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
			if wordStatus == "current" && charIndex == len(typedWord) {
				view += renderCursor(char)
			} else if wordStatus == "passed" {
				view += renderMissedChar(char)
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

func renderDefaultStyle(s string) string {
	return lipgloss.NewStyle().
		Background(defaultBackgroundColor).
		Foreground(defaultForegroundColor).
		Render(s)
}

func renderMissedChar(char string) string {
	return lipgloss.NewStyle().
		Foreground(missedCharColor).
		Render(char)
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
		Render(char) + "\x1b[40m"
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
		finalViewString += renderDefaultStyle("\n")
	}

	return finalViewString
}
