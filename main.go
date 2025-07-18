package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	greetingMessage = "Ready For the Race!!!!!"

	// ANSI escape codes for colors
	resetColor = "\033[0m"
	redColor   = "\033[31m"
	greenColor = "\033[32m"
	cyanColor  = "\033[36m"
	grayColor  = "\033[37m"

	// ANSI escape codes for cursor/screen manipulation
	carriageReturn      = "\r"
	carriageNewLine     = "\r\n"
	hideCursor          = "\033[?25l"
	showCursor          = "\033[?25h"
	deleteTillNewLine   = "\033[K"
	underLineText       = "\033[4m"
	clearScreen         = "\033[2J"
	moveCursorToHome    = "\033[H"
	moveCursorUpOneLine = "\033[F"

	// backspace char
	backSpaceChar = 127
)

// setupTerminalRawMode puts the terminal into raw mode and returns the original state.
func setupTerminalRawMode(fd int) (*term.State, error) {
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to make terminal raw: %w", err)
	}
	fmt.Print(hideCursor) // Hide cursor immediately after raw mode is set
	return oldState, nil
}

// restoreTerminalMode restores the terminal to its original state.
func restoreTerminalMode(fd int, oldState *term.State) {
	fmt.Print(showCursor) // Show cursor before restoring
	if err := term.Restore(fd, oldState); err != nil {
		log.Printf("error restoring terminal: %v", err) // Log but don't panic on defer
	}
}

// clearAndResetCursor clears the screen and moves the cursor to home.
func clearAndResetCursor() {
	fmt.Print(clearScreen)
	fmt.Print(moveCursorToHome)
}

// printToTerminal prints text with optional color.
func printToTerminal(text string, color string) {
	if color == "" {
		fmt.Print(text)
		return
	}
	fmt.Print(color)
	fmt.Print(text)
	fmt.Print(resetColor)
}

// moveCursorUpAndClearLine moves cursor up one line and clears it.
func moveCursorUpAndClearLine() {
	fmt.Print(moveCursorUpOneLine)
	fmt.Print(deleteTillNewLine)
}

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
func (gs *GameState) ProcessInput(char rune) {
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
	reader := bufio.NewReader(os.Stdin)
	for gs.Position < len(gs.TargetText) {
		char, _, err := reader.ReadRune()
		if err != nil {
			log.Printf("reading buffer: %v", err)
			return
		}

		gs.ProcessInput(char)

		width, _, err := term.GetSize(int(fd))
		if err != nil {
			log.Printf("reading terminal size: %v", err)
		}
		gs.Render(width)
	}
}

func (gs *GameState) ShowGameResult() {
	fmt.Print(resetColor)
	fmt.Print(carriageNewLine)
	duration := time.Since(gs.StartTime).Seconds()
	stats := GetStats([]rune(gs.TargetText), gs.UserInput, duration)
	fmt.Printf("result: %s", stats)
	fmt.Print(carriageNewLine)
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

func initializeGame() *GameState {
	text := GetRandomText()
	gs := NewGameState(text)

	clearAndResetCursor()
	printToTerminal(greetingMessage+carriageNewLine, grayColor)
	printToTerminal(gs.TargetText, cyanColor)
	fmt.Print(carriageReturn)

	return gs
}

func main() {
	logFile, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)

	fd := os.Stdin.Fd()
	oldState, err := setupTerminalRawMode(int(fd))
	if err != nil {
		panic(err)
	}
	defer restoreTerminalMode(int(fd), oldState)

	gs := initializeGame()

	gs.RunGameLoop(fd)
	gs.ShowGameResult()
}
