package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// GameState holds the dynamic state of the current game
type GameState struct {
	mu          sync.Mutex
	TargetText  string
	UserInput   []rune
	Position    int
	StartTime   time.Time
	Width       int
	wrappedText string
}

// NewGameState initializes a new game state.
func NewGameState(targetText string, width int) *GameState {
	return &GameState{
		TargetText:  targetText,
		UserInput:   make([]rune, len(targetText)),
		Position:    0,
		StartTime:   time.Now(),
		Width:       width,
		wrappedText: getWrappedText(targetText, width),
	}
}

// ProcessInput handles a single character input from the user.
func (gs *GameState) ProcessInput(char byte) {
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

func (gs *GameState) Render() {
	linesToClear := strings.Count(gs.wrappedText, "\n") + 1
	var b strings.Builder

	for i := range linesToClear {
		b.WriteString(deleteTillNewLine)
		if i < linesToClear-1 {
			b.WriteString(moveCursorUpOneLine)
		}
	}
	b.WriteString(carriageReturn)
	fmt.Print(b.String())
	b.Reset()

	userInputIdx := 0

	for _, char := range gs.wrappedText {
		b.WriteString(resetColor)
		if char == '\n' {
			b.WriteString(carriageNewLine)
			continue
		}

		// mark the chars as cyan which are still not written
		if userInputIdx >= gs.Position {
			// the current char should show an underline underneath (for virtual cursor)
			if userInputIdx == gs.Position {
				b.WriteString(underLineText)
			}
			b.WriteString(cyanColor)
			b.WriteRune(char)
			userInputIdx++
			continue
		}

		toColor := greenColor
		// the typed char is correct
		if gs.UserInput[userInputIdx] != char {
			toColor = redColor
		}
		// the type char is incorrect
		b.WriteString(toColor)
		b.WriteRune(char)
		userInputIdx++
	}
	fmt.Print(b.String())
}

func (gs *GameState) RunGameLoop() {
	reader := os.Stdin
	for gs.Position < len(gs.TargetText) {
		var charArr [4]byte
		n, err := reader.Read(charArr[:])
		if err != nil || n != 1 {
			log.Printf("reading buffer: %v %d", err, n)
			return
		}

		// 3 is for ctrl+c and 27 for esc key
		if charArr[0] == 3 || charArr[0] == 27  {
			return
		}

		gs.ProcessInput(charArr[0])
		gs.Render()
	}
}

func (gs *GameState) Reset(width int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.Position = 0
	gs.StartTime = time.Now()
	clear(gs.UserInput)
	gs.Width = width
	gs.wrappedText = getWrappedText(gs.TargetText, width)

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
