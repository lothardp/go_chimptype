package main

import (
	"errors"
	"io"
)

type KeyType int

const (
	Char KeyType = iota
	Backspace
	Space
	Enter
	Esc
)

type Key struct {
	keyType KeyType
	char    rune
}

type TestState struct {
	wordList    [][]Key
	wordIndex   int
	typedChars  []Key
	rawCharList []Key
	finished    bool
}

type TestResult struct {
	finalState TestState
}

func NewTestState(wordList []string) TestState {
	return TestState{
		wordList:    wordsToKeyLists(wordList),
		wordIndex:   0,
		typedChars:  []Key{},
		rawCharList: []Key{},
		finished:    false,
	}
}

func (ts *TestState) HandleKey(key Key) error {
	ts.rawCharList = append(ts.rawCharList, key)

	switch key.keyType {
	case Char:
		ts.handleChar(key)

	case Space:
		ts.handleSpace()

	case Backspace:
		ts.handleBackspace()

	default:
		return errors.New("Invalid key type")
	}

	return nil
}

func (ts *TestState) handleChar(key Key) {
	// TODO: Add logic to finish test without space
	ts.typedChars = append(ts.typedChars, key)
}

func (ts *TestState) handleSpace() {
	ts.typedChars = append(ts.typedChars, Key{keyType: Space})
	if ts.wordIndex == len(ts.wordList)-1 {
		ts.finished = true
	} else {
		ts.wordIndex++
	}
}

func (ts *TestState) handleBackspace() {
	if ts.wordIndex == 0 && len(ts.typedChars) == 0 {
		return
	}

	charIndex := ts.charIndex()
	if charIndex == 0 {
		ts.wordIndex--
	}
	ts.typedChars = ts.typedChars[:len(ts.typedChars)-1]
}

// Returns the index in the current word of the char the cursor is on
func (ts *TestState) charIndex() int {
	charIndex := len(ts.typedChars) - 1
	for i := len(ts.typedChars) - 1; i >= 0; i-- {
		key := ts.typedChars[i]
		if key.keyType == Space {
			charIndex = len(ts.typedChars) - i - 1
			break
		}
	}

	return charIndex
}

func (ts *TestState) Draw(w io.Writer) {
	// Drawing logic
}

func wordsToKeyLists(wordList []string) [][]Key {
	keyLists := [][]Key{}
	for _, word := range wordList {
		keyLists = append(keyLists, wordToKeyList(word))
	}

	return keyLists
}

func wordToKeyList(word string) []Key {
	keyList := []Key{}
	for _, char := range word {
		keyList = append(keyList, Key{Char, char})
	}

	return keyList
}
