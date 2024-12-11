package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Printf("\n error in restoring: %v", err)
	}()

	terminal := term.NewTerminal(os.Stdout, "Ready for the race??")
	terminal.SetPrompt(GetRandomText())
	terminal.ReadLine()
}
