package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

const greetingMessage = "Ready For the Race!!!!!"

func main() {
	fd := os.Stdin.Fd()
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

	text := GetRandomText()

	fmt.Println(greetingMessage)
	fmt.Println(text)
	userInput := make([]rune, len(text))
	pos := 0

	for pos < len(text) {
		buf := make([]byte, 1)
		os.Stdin.Read(buf)

		char := buf[0]

		if char == 127 {
			fmt.Print("\b \b")
			continue
		}

		if char != text[pos] {
			pos++
			fmt.Printf("diff: %s, %s\n", string(char), string(text[pos]))
			continue
		}
        pos++
		userInput[pos] = rune(char)

	}
}
