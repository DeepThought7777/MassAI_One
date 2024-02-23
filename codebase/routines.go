package codebase

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
)

// displayAndOptionallyExit() displays an error message,
// waits for enter to be pressed, then optionally exits the program.
func displayAndOptionallyExit(errorMessage string, exit bool) {
	fmt.Println(errorMessage)
	fmt.Println(">>> Press the [ENTER] key to end the program <<<")
	_, err := fmt.Scanln()
	if !exit || err != nil {
		return
	}
	os.Exit(-1)
}

func CreateFolder(folderPath string) error {
	err := os.MkdirAll(folderPath, 0755) // Create the folder with permissions set to 0755 (readable/writeable by owner and group, readable by others)
	if err != nil {
		return fmt.Errorf("failed to create folder: %v", err)
	}
	fmt.Printf("Folder created successfully: %s\n", folderPath)
	return nil
}

// RuneToSignals translates a single rune into 32 separate parallel signals
// represented as a slice of booleans
func RuneToSignals(r rune) ([]bool, bool) {
	signals := make([]bool, 32)
	// Check if the rune is valid
	_, ok := isValidRune(r)
	if !ok {
		// If not valid, return signals for NULL rune
		return signals, ok
	}

	// If valid, proceed with the conversion
	for i := 0; i < 32; i++ {
		signals[31-i] = (r>>i)&1 == 1
	}
	return signals, ok
}

// SignalsToRune reconstructs a single rune from 32 separate parallel signals
// represented as a slice of booleans
func SignalsToRune(signals []bool) (rune, bool) {
	var r rune
	for i, signal := range signals {
		if signal {
			r |= 1 << uint(31-i)
		}
	}
	return isValidRune(r)
}

// isValidRune checks if the given rune is valid
func isValidRune(r rune) (rune, bool) {
	if r >= 0 && r <= 0x10FFFF {
		return r, true
	}
	// Return NULL rune if invalid
	return '\x00', false
}

func NewGUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80

	returnGuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	fmt.Println(returnGuid)
	return returnGuid
}
