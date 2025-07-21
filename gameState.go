package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// GameState holds the dynamic state of the current game
type GameState struct {
	mu         sync.Mutex
	TargetText string
	UserInput  []rune
	Position   int
	StartTime  time.Time
}

// NewGameState initializes a new game state.
func NewGameState(targetText string) *GameState {
	return &GameState{
		TargetText: targetText,
		UserInput:  make([]rune, len(targetText)),
		Position:   0,
		StartTime:  time.Now(),
	}
}

// ProcessInput handles a single character input from the user.
func (gs *GameState) ProcessInput(char byte) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	if char == backSpaceChar {
		gs.Position = max(gs.Position-1, 0)
		return
	}
	// Ensure we don't write beyond allocated buffer
	if gs.Position >= len(gs.TargetText) {
		log.Print("GameState Position has exceeded TargetText")
		gs.Position++
		return
	}
	gs.UserInput[gs.Position] = rune(char)
	gs.Position++
}

func (gs *GameState) Render(width int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	wrappedText := getWrappedText(gs.TargetText, width)
	linesToClear := strings.Count(wrappedText, "\n") + 1
	for i := range linesToClear {
		fmt.Print(deleteTillNewLine)
		if i < linesToClear-1 {
			fmt.Print(moveCursorUpOneLine)
		}
	}
	fmt.Print(carriageReturn)

	userInputIdx := 0

	for _, char := range wrappedText {
		fmt.Print(resetColor)
		if char == '\n' {
			fmt.Print(carriageNewLine)
			continue
		}

		toPrint := fmt.Sprintf("%c", char)
		// mark the chars as cyan which are still not written
		if userInputIdx >= gs.Position {
			// the current char should show an underline underneath (for virtual cursor)
			if userInputIdx == gs.Position {
				fmt.Print(underLineText)
			}
			printToTerminal(toPrint, cyanColor)
			userInputIdx++
			continue
		}

		toColor := greenColor
		// the typed char is correct
		if gs.UserInput[userInputIdx] != char {
			toColor = redColor
		}
		// the type char is incorrect
		printToTerminal(toPrint, toColor)
		userInputIdx++
	}
}

func (gs *GameState) RunGameLoop(fd uintptr) {
	reader := os.Stdin
	for gs.Position < len(gs.TargetText) {
		var charArr [4]byte
		n, err := reader.Read(charArr[:])
		if err != nil || n != 1 {
			log.Printf("reading buffer: %v %d", err, n)
			return
		}

		gs.ProcessInput(charArr[0])

		width, _, err := term.GetSize(int(fd))
		if err != nil {
			log.Printf("reading terminal size: %v", err)
		}
		gs.Render(width)
	}
}

func (gs *GameState) Reset() {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.Position = 0
	gs.StartTime = time.Now()
	clear(gs.UserInput)

	clearAndResetCursor()
	printToTerminal(greetingMessage+carriageNewLine, grayColor)
	printToTerminal(gs.TargetText, cyanColor)
	fmt.Print(carriageReturn)
}

func (gs *GameState) ShowGameResult() {
	fmt.Print(resetColor)
	fmt.Print(carriageNewLine)
	duration := time.Since(gs.StartTime).Seconds()
	stats := GetStats([]rune(gs.TargetText), gs.UserInput, duration)
	fmt.Printf("result: %s", stats)
	fmt.Print(carriageNewLine)
}
