package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

type KeyType int

const (
	Char KeyType = iota
	Backspace
	Space
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
	startTime   time.Time
	finishTime  time.Time
}

type TestResult struct {
	finalState TestState
	duration   time.Duration
	netWPM     float64
	rawWPM     float64
	accuracy   float64
	correct    int
	errors     int
	missed     int
	extra      int
}

func NewTestState(wordList []string) TestState {
	return TestState{
		wordList:    stringsToKeyLists(wordList),
		wordIndex:   0,
		typedChars:  []Key{},
		rawCharList: []Key{},
		finished:    false,
	}
}

func newTestResult(testState TestState) TestResult {
	duration := testState.finishTime.Sub(testState.startTime)
	netWPM, rawWPM := calculateWPM(testState, duration)
	accuracy := calculateAccuracy(testState)

	return TestResult{
		finalState: testState,
		duration:   duration,
		netWPM:     netWPM,
		rawWPM:     rawWPM,
		accuracy:   accuracy,
	}
}

// TODO: There is different formulas for calculating calculate wpm, choose best one
func calculateWPM(state TestState, duration time.Duration) (float64, float64) {
	words := keyListsToString(state.wordList)
	typedWords := strings.Split(keyListToString(state.typedChars), " ")

	if len(words) != len(typedWords) {
		log.Print("Error: word count mismatch")
		return 0, 0
	}

	typedWordsNum := float64(len(state.typedChars)) / 5.0
	wrongWordsNum := float64(wrongTypedWords(words, typedWords))

	rawWPM := typedWordsNum / duration.Minutes()

	netWPM := rawWPM - wrongWordsNum/duration.Minutes()

	return netWPM, rawWPM
}

func wrongTypedWords(words []string, typedWords []string) int {
	wrongWords := 0
	for i := range words {
		if words[i] != typedWords[i] {
			wrongWords++
		}
	}
	return wrongWords
}

func calculateAccuracy(state TestState) float64 {
	return 0
}

func (ts *TestState) HandleKey(key Key) error {
	if ts.startTime.IsZero() {
		ts.startTime = time.Now()
	}

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
	ts.typedChars = append(ts.typedChars, key)

	// TODO: Add logic to finish test without space
	inLastWord := ts.wordIndex == len(ts.wordList)-1

	if inLastWord {
		typedWords := strings.Split(keyListToString(ts.typedChars), " ")
		lastTypedWord := typedWords[len(typedWords)-1]

		if lastTypedWord == keyListToString(ts.wordList[ts.wordIndex]) {
			ts.finishTest()
		}
	}
}

func (ts *TestState) handleSpace() {
	ts.typedChars = append(ts.typedChars, Key{Space, ' '})
	if ts.wordIndex == len(ts.wordList)-1 {
		ts.finishTest()
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

func stringsToKeyLists(wordList []string) [][]Key {
	keyLists := [][]Key{}
	for _, word := range wordList {
		keyLists = append(keyLists, stringToKeyList(word))
	}

	return keyLists
}

func stringToKeyList(word string) []Key {
	keyList := []Key{}
	for _, char := range word {
		keyList = append(keyList, Key{Char, char})
	}

	return keyList
}

func keyListsToString(keyLists [][]Key) []string {
	strs := []string{}
	for _, keyList := range keyLists {
		strs = append(strs, keyListToString(keyList))
	}

	return strs
}

func keyListToString(keyList []Key) string {
	str := ""
	for _, key := range keyList {
		str += string(key.char)
	}

	return str
}

func (ts *TestState) finishTest() {
	ts.finished = true
	ts.finishTime = time.Now()
}
