package mind

import (
	"github.com/DeepThought7777/MassAI/codebase"
)

// Linkup contains all information that relates to a single Linkup to an input/output device.
type Linkup struct {
	Name   string   `json:"linkupname"` // Unique ID
	Nerves []string `json:"nerve_ids"`  // The pre-assigned set of Neurons
	Dirty  bool     `json:"dirty"`      // is the Linkup not synced to disk?

	nerves []*Neuron // The loaded set of Neurons (not stored here)
}

// NewLinkup returns a new Linkup properly initialized
func NewLinkup(brain *Brain, nervecount int, name string) Linkup {
	if name == "" {
		name = codebase.NewGUID()
	}
	newLinkup := Linkup{
		Name:   name,
		Nerves: make([]string, 0),
		nerves: make([]*Neuron, 0),
		Dirty:  true,
	}
	for i := 0; i < nervecount; i++ {
		neuron := NewNeuron(*brain)
		newLinkup.nerves = append(newLinkup.nerves, neuron)
		newLinkup.Nerves = append(newLinkup.Nerves, neuron.Name)
		brain.StoreNeuron(neuron)
	}

	brain.Linkups[newLinkup.Name] = newLinkup
	return newLinkup
}

func (l *Linkup) Sense(brain *Brain) {
	for _, neuron := range l.nerves {
		neuron.Accumulate(brain)
	}
}

func Act(brain *Brain) {
	//for _, neuron := range b.nerves {
	//	neuron.Distribute(brain)
	//}
}
