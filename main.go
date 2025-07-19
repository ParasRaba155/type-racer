package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

    go handleResize(gs)

	gs.RunGameLoop(fd)
	gs.ShowGameResult()
}

func handleResize(gs *GameState) {
	resize := make(chan os.Signal, 1)
	signal.Notify(resize, syscall.SIGWINCH)

	for {
		<-resize
		gs.Reset()
	}
}

