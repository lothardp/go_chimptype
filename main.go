package main

import (
	"golang.org/x/term"
	"os"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	chimptypeInstance := NewChimpType()
	chimptypeInstance.Start()
}
