package mind

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/DeepThought7777/MassAI/codebase"
)

// Brain contains all information that relates to the brain as a whole.
type Brain struct {
	BrainName     string            `json:"brain_name"`   // Unique name to be used for the brain
	StoragePath   string            `json:"storage_path"` // The folder the data needs to be put in
	Linkups       map[string]Linkup `json:"linkups"`      // The set of Linkups for this brain
	neurons       map[string]*Neuron
	chargingQueue []*Neuron
	firingQueue   []*Neuron
}

// NewBrain creates a new Brain structure, and returns it (not a reference).
func NewBrain(basePath, name string) Brain {
	newGuid := codebase.NewGUID()
	if name != "" {
		newGuid = name
	}
	newBrain := Brain{
		BrainName:   newGuid,
		StoragePath: filepath.Join(basePath, newGuid),
		neurons:     make(map[string]*Neuron, 0),
		Linkups:     make(map[string]Linkup, 0),
	}
	err := codebase.CreateFolder(newBrain.StoragePath)
	if err != nil {
		log.Fatalf("Error creating folder: %v", err)
	}

	return newBrain
}

// LoadedSize returns the size of the map containing the Neurons loaded In memory.
func (b *Brain) LoadedSize() int {
	return len(b.neurons)
}

func (b *Brain) AddNeuron(neuron *Neuron) {
	b.neurons[neuron.Name] = neuron
}

func (b *Brain) DeleteNeuron(guid string) {
	delete(b.neurons, guid)
}

// GetNeuron returns a reference to a Neuron if it is In the loaded map, or a nil reference otherwise.
func (b *Brain) GetNeuron(ID string) *Neuron {
	neuron := b.GetNeuronIfLoaded(ID)
	if neuron != nil {
		return neuron
	}

	neuron, err := b.LoadNeuron(ID)
	if err != nil {
		return nil
	}
	return neuron
}

// GetNeuronIfLoaded returns a reference to a Neuron if it is In the loaded map, or a nil reference otherwise.
func (b *Brain) GetNeuronIfLoaded(ID string) *Neuron {
	neuron, ok := b.neurons[ID]
	if ok {
		return neuron
	}
	return nil
}

func (b *Brain) GetStorageFolder(baseFolder string) string {
	return filepath.Join(baseFolder, b.BrainName)
}

func (b *Brain) StoreBrain() error {
	if b == nil || b.BrainName == "" {
		return nil
	}

	spec := filepath.Join(b.StoragePath, b.BrainName+".mai")
	file, err := os.Create(spec)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(b)
	if err != nil {
		err = file.Close()
		return err
	}

	err = file.Close()
	return nil
}

func LoadBrain(storagePath, brainName string) (*Brain, error) {
	spec := filepath.Join(storagePath, brainName+".mai")
	file, err := os.Open(spec)
	if err != nil {
		return nil, err
	}

	var brain Brain
	err = json.NewDecoder(file).Decode(&brain)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	brain.neurons = make(map[string]*Neuron, 0)

	_ = file.Close()
	return &brain, nil
}

func (b *Brain) StoreNeuron(neuron *Neuron) error {
	if neuron == nil || neuron.Name == "" {
		return nil
	}

	neuron.Dirty = false

	spec := filepath.Join(b.StoragePath, neuron.Name+".json")
	file, err := os.Create(spec)
	if err != nil {
		return err
	}

	err = json.NewEncoder(file).Encode(neuron)
	if err != nil {
		err = file.Close()
		return err
	}

	err = file.Close()
	return nil
}

func (b *Brain) LoadNeuron(neuronName string) (*Neuron, error) {
	spec := filepath.Join(b.StoragePath, neuronName+".json")
	data, err := os.ReadFile(spec)
	if err != nil {
		return nil, err
	}

	var neuron Neuron
	err = json.Unmarshal(data, &neuron)
	if err != nil {
		return nil, err
	}

	b.neurons[neuron.Name] = &neuron
	return &neuron, nil
}

func (b *Brain) LoadAllLinkupNeurons() {
	for _, linkup := range b.Linkups {
		b.LoadLinkupNeurons(linkup.Name)
	}
}

func (b *Brain) LoadLinkupNeurons(linkupName string) {
	link := b.Linkups[linkupName]
	for _, nerve := range link.Nerves {
		_, _ = b.LoadNeuron(nerve)
	}
}

func (b *Brain) QueueForCharging(neuron *Neuron) {
	b.chargingQueue = append(b.chargingQueue, neuron)
}

func (b *Brain) UnQueueFromCharging() *Neuron {
	if len(b.chargingQueue) == 0 {
		return nil
	}
	n := b.chargingQueue[0]
	b.chargingQueue = b.chargingQueue[1:]
	return n
}

func (b *Brain) QueueForFiring(neuron *Neuron) {
	b.firingQueue = append(b.firingQueue, neuron)
}

func (b *Brain) UnQueueFromFiring() *Neuron {
	if len(b.firingQueue) == 0 {
		return nil
	}
	n := b.firingQueue[0]
	b.firingQueue = b.firingQueue[1:]
	return n
}
