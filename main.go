package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

const (
	greetingMessage = "Ready For the Race!!!!!"

	// ascii escape chars for colors
	grayColor  = "\033[37m"
	resetColor = "\033[0m"
	redColor   = "\033[31m"
	greenColor = "\033[32m"

	// carriage chars
	carriageReturn  = "\r"
	carriageNewLine = "\r\n"

	// special chars
	hideCursor        = "\033[?25l"
	showCursor        = "\033[?25h"
	deleteTillNewLine = "\033[K"
)

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
		fmt.Printf(showCursor)
		log.Printf("\n error in restoring: %v", err)
	}()

	text := GetRandomText()
	text = abcd()

	fmt.Print(grayColor)
	fmt.Print(greetingMessage, carriageNewLine)
	fmt.Print(resetColor)
	fmt.Print(grayColor)
	fmt.Print(text)
	fmt.Print(resetColor)
	fmt.Print(carriageReturn)
	fmt.Print(hideCursor)

	userInput := make([]rune, len(text))
	pos := 0

	for pos < len(text) {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if n != 1 || err != nil {
			log.Printf("reading buffer: %d: %v", n, err)
			return
		}

		char := buf[0]

		if char != 127 {
			userInput[pos] = rune(char)
			pos++
		} else {
			pos = max(pos-1, 0)
		}

		fmt.Print(carriageReturn)
		for i, char := range text {
			if i > len(userInput) {
				fmt.Print(grayColor)
				fmt.Printf("%c", char)
				fmt.Print(resetColor)
				continue
			}
			if userInput[i] == char {
				fmt.Print(greenColor)
				fmt.Printf("%c", char)
				fmt.Print(resetColor)
				continue
			}
			fmt.Print(redColor)
			fmt.Printf("%c", char)
			fmt.Print(resetColor)
		}
	}
	fmt.Print(resetColor)
}

func abcd() string {
	var c [26]byte
	var i byte
	for i = 65; i <= 90; i++ {
		c[i-65] = i
	}
	return string(c[:])
}
