package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

func main() {
	fd := os.Stdout.Fd()
	logFile, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	oldState, err := term.MakeRaw(int(fd))
	if err != nil {
		panic(err)
	}
	defer func() {
		err := term.Restore(int(fd), oldState)
		fmt.Printf("\n error in restoring: %v", err)
	}()

	terminal := term.NewTerminal(os.Stdout, GetRandomTextWithGreeting())
	terminal.AutoCompleteCallback = autoComplete
	log.Printf("Ready for the Race!!!!\n")
	line, err := terminal.ReadLine()
	if err != nil {
		fmt.Printf("error in reading the line: %v", err)
	}
	fmt.Println("LINE", line)
}

func autoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	log.Printf("key : %v,", key)
	if key == 10 {
		return "FUCK THIS SHIT", pos + 1, false
	}
	return line, pos, false
}
