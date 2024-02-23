package mind

import (
	"bufio"
	"fmt"
	"github.com/DeepThought7777/MassAI/codebase"
	"os"
	"sync"
)

// Linkup contains all information that relates to a single Linkup to an input/output device.
type LinkupConsole struct {
	Linkup // the connected Linkup module
}

// NewLinkupConsole returns a new Linkup properly initialized
func NewLinkupConsole(brain *Brain, name string) LinkupConsole {
	if name == "" {
		name = codebase.NewGUID()
	}
	newLinkupConsole := LinkupConsole{
		Linkup: NewLinkup(brain, 32, name),
	}

	return newLinkupConsole
}

func (c *LinkupConsole) StartLinkup() {
	// Channel to communicate between the two goroutines
	signalChan := make(chan []bool)

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine to convert runes from stdin to signals and send to neural network
	go func() {
		defer wg.Done()
		defer close(signalChan)
		reader := bufio.NewReader(os.Stdin)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				return
			}
			signals, ok := codebase.RuneToSignals(r)
			// Simulating sending signals to neural network
			if ok {
				signalChan <- signals
			}
		}
	}()

	// Goroutine to receive signals from the neural network and convert them into runes
	go func() {
		defer wg.Done()
		for signals := range signalChan {
			// Simulating receiving signals from neural network
			reconstructedRune, ok := codebase.SignalsToRune(signals)
			if ok {
				fmt.Print(string(reconstructedRune))
			}
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}
