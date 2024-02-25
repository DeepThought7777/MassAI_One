package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/DeepThought7777/MassAI/api"
	"github.com/DeepThought7777/MassAI/codebase"
	"github.com/DeepThought7777/MassAI/mind"
	"github.com/gorilla/mux"
)

func Storer(wg *sync.WaitGroup, stopChan chan string) {
	defer wg.Done()
	for {
		select {
		case _ = <-stopChan:
			fmt.Println("Storer received stop signal")
			return
		default:
			// Do the storing work here
			fmt.Println("Storer is working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func Loader(wg *sync.WaitGroup, stopChan chan string) {
	defer wg.Done()
	for {
		select {
		case _ = <-stopChan:
			fmt.Println("Loader received stop signal")
			return
		default:
			// Do the loading work here
			fmt.Println("Loader is working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func Cleaner(wg *sync.WaitGroup, stopChan chan string) {
	defer wg.Done()
	for {
		select {
		case stopMsg := <-stopChan:
			switch stopMsg {
			case "STOP":
				fmt.Println("Cleaner received stop signal")
				return
			case "exit":
				fmt.Println("Cleaner received exit signal")
				return
			case "pause":
				fmt.Println("Cleaner received pause signal")
				// Handle pause logic here
			default:
				fmt.Println("Cleaner received unknown signal:", stopMsg)
			}
			return
		default:
			// Do the cleaning work here
			fmt.Println("Cleaner is working...")
			time.Sleep(1 * time.Second)
		}
	}
}

func parseInt(s string) int {
	value, err := fmt.Sscan(s)
	if err != nil {
		return 0
	}
	return value
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/register", api.RegisterEntity).Methods("GET")
	router.HandleFunc("/v1/unregister", api.UnregisterEntity).Methods("GET")
	router.HandleFunc("/v1/connect", api.ConnectEntity).Methods("GET")
	router.HandleFunc("/v1/disconnect", api.DisconnectEntity).Methods("GET")
	router.HandleFunc("/v1/send_inputs", api.SendInputs).Methods("GET")
	router.HandleFunc("/v1/get_outputs", api.GetOutputs).Methods("GET")

	http.Handle("/", router)
	codebase.NewGUID()

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

/*
func main() {
	// Define command-line parameters
	path := flag.String("path", "", "Folder to put stored files in.")
	name := flag.String("name", "", "Name for the brain file, except extension")
	test := flag.Bool("test", false, "Test parameter (optional)")

	// Parse command-line parameters
	flag.Parse()

	// Validate required parameters
	if *path == "" {
		fmt.Println("Usage: ./your_program --basepath <value> --name <value> [--test]")
		return
	}

	// Print the values
	fmt.Printf("Base Path: %s\n", *path)
	fmt.Printf("Name: %s\n", *name)

	// Check if the --test flag is provided
	if *test {
		fmt.Println("EXECUTING TEST!!!")
		selfTest(*path, *name)
		os.Exit(0)
	}

	os.Exit(-1)

		addr := ":8080"

		filename := os.Args[1]
		fmt.Println("Filename:", filename)

		http.HandleFunc("/create", createHandler)
		http.HandleFunc("/startup", startupHandler)
		http.HandleFunc("/shutdown", shutdownHandler)
		http.HandleFunc("/connect", connectHandler)
		http.HandleFunc("/disconnect", disconnectHandler)

		go func() {
			ListenAndServe(addr)
		}()

		var wg sync.WaitGroup
		stopLoaderCh := make(chan string)
		stopStorerCh := make(chan string)
		stopCleanerCh := make(chan string)

		wg.Add(3)
		go Loader(&wg, stopLoaderCh)
		go Storer(&wg, stopStorerCh)
		go Cleaner(&wg, stopCleanerCh)

		// Send a stop signal to the workers after 5 seconds
		time.Sleep(5 * time.Second)
		stopLoaderCh <- "STOP"
		stopStorerCh <- "STOP"
		stopCleanerCh <- "STOP"

		Think()
		wg.Wait()
}
*/

func selfTest(path, name string) {
	brain := mind.NewBrain(path, name)
	// NEW and STORE brain, executed only once for now...
	/*
		// Add a Linkup to it with 10 nerve endings
		_ = mind.NewLinkup(&brain, 10, "TestLinkup")
		err := brain.StoreBrain()
		if err != nil {
			os.Exit(-2)
		}
		err = brain.StoreBrain()
		if err != nil {
			os.Exit(-2)
		}
	*/

	loadedBrain, err := mind.LoadBrain(brain.StoragePath, brain.BrainName)
	if err != nil {
		os.Exit(-2)
	}

	loadedBrain.LoadAllLinkupNeurons()

	fmt.Printf("Brain size:     %d\n", loadedBrain.LoadedSize())
	fmt.Printf("Loaded brain:   %v\n", loadedBrain)

	/*
		// Build a Brain of some Neurons
		// with each neuron connected to the next one
		firstNeuron := mind.NewNeuron(brain)
		firstNeuronGuid := firstNeuron.Name
		neuron := firstNeuron
		start := time.Now()
		for i := 0; i < 99; i++ {
			firstNeuronGuid = neuron.Name
			lastNeuron := neuron
			neuron = mind.NewNeuron(brain)
			neuron.AddNeuronIn(lastNeuron)
		}
		elapsed := time.Since(start)
		fmt.Printf("Brain size:     %d\n", brain.LoadedSize())
		fmt.Printf("Creation time:  %s\n", elapsed)

		// Test for lookup of a Neuron via its ID()
		// map access is FAST!!!
		start = time.Now()
		foundNeuron := brain.GetNeuronIfLoaded(firstNeuron.Name)
		for {
			if foundNeuron == nil {
				break
			}
			foundNeuron = brain.GetNeuronIfLoaded(foundNeuron.GetFirstOut())
			// fmt.Printf(".")
		}
		elapsed = time.Since(start)
		fmt.Printf("Lookup time:    %s\n", elapsed)

		// Test for storing of Neurons
		start = time.Now()
		foundNeuron = brain.GetNeuronIfLoaded(firstNeuron.Name)
		for {
			if foundNeuron == nil {
				break
			}
			nextNeuronGuid := foundNeuron.GetFirstOut()
			brain.DeleteNeuron(foundNeuron.Name)
			foundNeuron = brain.GetNeuronIfLoaded(nextNeuronGuid)
			err := brain.StoreNeuron(foundNeuron)
			if err != nil {

			}
		}
		elapsed = time.Since(start)
		fmt.Printf("Store time:     %s\n", elapsed)
		fmt.Printf("Brain size:     %d\n", brain.LoadedSize())

		// Test for loading of Neurons
		start = time.Now()
		nextNeuronGuid := firstNeuronGuid
		for {
			newNeuron, _ := brain.LoadNeuron(nextNeuronGuid)
			if newNeuron == nil {
				break
			}
			nextNeuronGuid = newNeuron.GetFirstIn()
			brain.AddNeuron(newNeuron)
		}
		elapsed = time.Since(start)
		fmt.Printf("Load time:      %s\n", elapsed)
		fmt.Printf("Brain size:     %d\n", brain.LoadedSize())
	*/
}
