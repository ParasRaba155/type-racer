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
	terminal.AutoCompleteCallback = autoComplete
	terminal.SetPrompt(GetRandomText())
	fmt.Printf("Ready for the Race!!!!")
	line, err := terminal.ReadLine()
	if err != nil {
		fmt.Printf("error in reading the line: %v", err)
	}
	fmt.Println("LINE", line)
}

func autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	fmt.Println("autoComplete")
	if key == ' ' {
		return line, pos, true
	}
	if key == rune(10) {
		fmt.Println("FUCK THIS SHIT")
		return line, pos + 1, false
	}
	return line, pos, false
}
