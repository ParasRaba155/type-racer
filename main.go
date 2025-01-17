package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	greetingMessage = "Ready For the Race!!!!!"

	// ascii escape chars for colors
	resetColor = "\033[0m"
	redColor   = "\033[31m"
	greenColor = "\033[32m"
	cyanColor  = "\033[36m"
	grayColor  = "\033[37m"

	carriageReturn  = "\r" // takes the cursor to the very beginning
	carriageNewLine = "\r\n"

	// special chars
	hideCursor        = "\033[?25l"
	showCursor        = "\033[?25h"
	deleteTillNewLine = "\033[K"
	underLineText     = "\033[4m"

	// backspace char
	backSpaceChar = 127

	// delete till
	clearScreen      = "\033[2J"
	moveCursorToHome = "\033[H"
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
		fmt.Print(showCursor)
		log.Printf("\n error in restoring: %v", err)
	}()

	text := GetRandomText()
	fmt.Print(clearScreen)

	fmt.Print(grayColor)
	fmt.Print(greetingMessage, carriageNewLine)
	fmt.Print(resetColor)
	fmt.Print(cyanColor)
	fmt.Print(text)
	fmt.Print(resetColor)
	fmt.Print(carriageReturn)
	fmt.Print(hideCursor)

	userInput := make([]rune, len(text))
	pos := 0
	start := time.Now()

	for pos < len(text) {
		buf := make([]byte, 4)
		n, err := os.Stdin.Read(buf)
		if n != 1 || err != nil {
			log.Printf("reading buffer: %d read: %v", n, err)
			return
		}

		char := buf[0]

		if char != backSpaceChar {
			userInput[pos] = rune(char)
			pos++
		} else {
			pos = max(pos-1, 0)
		}

		width, _, err := term.GetSize(int(fd))
		if err != nil {
			log.Printf("reading the terminal size: %v", err)
		}
		display(text, userInput, pos, width)
	}
	fmt.Print(resetColor)
	fmt.Print(carriageNewLine)
	diff := time.Since(start).Seconds()
	stats := GetStats([]rune(text), userInput, diff)
	fmt.Printf("result: %s", stats)
	fmt.Print(carriageNewLine)
}

// printABCD will print the char A to Z, and will do it repatatively until the width is reached
func printABCD(width int) string {
	var s strings.Builder
	repeat := (width / 26) + 1
	counter := 0
	for range repeat {
		for i := 65; i <= 90; i++ {
			if counter > width {
				return s.String()
			}
			counter++
			s.WriteByte(byte(i))
		}
	}
	return s.String()
}

// display will pretty print the text according to the userInput
func display(text string, userInput []rune, pos, width int) {
	fmt.Print(clearScreen)
	fmt.Print(moveCursorToHome)

	text = getWrappedText(text, width)
	userInputIdx := 0

	for i, char := range text {
		fmt.Print(resetColor)
		if text[i] == '\n' {
			fmt.Print(carriageNewLine)
			continue
		}

		// mark the chars as cyan which are still not written
		if userInputIdx >= pos {
			// the current char should show an underline underneath (for virtual cursor)
			if userInputIdx == pos {
				fmt.Print(underLineText)
			}
			fmt.Print(cyanColor)
			fmt.Printf("%c", char)
			fmt.Print(resetColor)
			userInputIdx++
			continue
		}

		if userInput[userInputIdx] == char {
			fmt.Print(greenColor)
			fmt.Printf("%c", char)
			fmt.Print(resetColor)
			userInputIdx++
			continue
		}
		fmt.Print(redColor)
		fmt.Printf("%c", char)
		fmt.Print(resetColor)
		userInputIdx++
	}
}

func getWrappedText(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if len(text) <= width {
		return text
	}

	var s strings.Builder
	remainingWidth := width
	for i := range text {
		if remainingWidth == 0 {
			s.WriteByte('\n')
			remainingWidth = width
		}
		s.WriteByte(text[i])
		remainingWidth--
	}
	return s.String()
}
