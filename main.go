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
	text = printAs(120)

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

		if char != 127 {
			userInput[pos] = rune(char)
			pos++
		} else {
			pos = max(pos-1, 0)
		}

		fmt.Print(carriageReturn)
		fmt.Print(deleteTillNewLine)
		// width here comes out to be the number of chars
		// since we are only dealing with 1 byte chars
		width, _, err := term.GetSize(int(fd))
		if err != nil {
			log.Fatalf("reading terminal size: %v", err)
			return
		}
		lineLen := 0

		for i, char := range text {
			fmt.Print(resetColor)
			log.Printf("i = %d and lineLen = %d and width = %d", i, lineLen, width)

			if lineLen >= width {
				fmt.Print(carriageNewLine)
				fmt.Print(deleteTillNewLine)
				lineLen = 0
			}

			// mark the chars as green which are still not written
			if i >= pos {
				fmt.Print(cyanColor)
				fmt.Printf("%c", char)
				lineLen++
				fmt.Print(resetColor)
				continue
			}
			if userInput[i] == char {
				fmt.Print(greenColor)
				fmt.Printf("%c", char)
				lineLen++
				fmt.Print(resetColor)
				continue
			}
			fmt.Print(redColor)
			fmt.Printf("%c", char)
			lineLen++
			fmt.Print(resetColor)
		}
		// Fill blank spaces for shorter input
		for lineLen < width {
			fmt.Print(" ")
			lineLen++
		}
	}
	fmt.Print(resetColor)
	fmt.Print(carriageNewLine)
	diff := time.Since(start).Seconds()
	stats := GetStats([]rune(text), userInput, diff)
	fmt.Printf("result: %s", stats)
	fmt.Print(carriageNewLine)
}

// printAs will print 'a's
func printAs(width int) string {
	var s strings.Builder
	for range width {
		s.WriteByte('a')
	}
	return s.String()
}
