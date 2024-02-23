package mind

import (
	"github.com/DeepThought7777/MassAI/codebase"
)

const FiringLimit = 7777
const ResidualCharge = -111

// Neuron contains all information that relates to a single neuron In the brain.
type Neuron struct {
	Name   string         `json:"brain_name"` // Unique ID
	In     map[string]int `json:"in"`         // IDs of incoming Neurons
	Out    map[string]int `json:"out"`        // IDs of outgoing Neurons
	Charge int            `json:"charge"`     // the charge level
	Dirty  bool           `json:"dirty"`      // is the neuron not synced to disk?
}

// NewNeuron returns a new neuron properly initialized
func NewNeuron(brain Brain) *Neuron {
	newNeuron := &Neuron{
		Name:   codebase.NewGUID(),
		In:     make(map[string]int, 0),
		Out:    make(map[string]int, 0),
		Charge: 0,
		Dirty:  true,
	}

	brain.neurons[newNeuron.Name] = newNeuron
	return newNeuron
}

// AddNeuronIn adds a neuron to the
func (n *Neuron) AddNeuronIn(neuron *Neuron) {
	n.In[neuron.Name] = 0
	neuron.Out[n.Name] = 0
}

func (n *Neuron) AddNeuronOut(neuron *Neuron) {
	n.Out[neuron.Name] = 0
	neuron.In[n.Name] = 0
}

func (n *Neuron) GetFirstIn() string {
	if n == nil || len(n.In) == 0 {
		return ""
	}
	for key, _ := range n.In {
		return key
	}
	return ""
}

func (n *Neuron) GetFirstOut() string {
	if n == nil || len(n.Out) == 0 {
		return ""
	}
	for key, _ := range n.Out {
		return key
	}
	return ""
}

func (n *Neuron) Accumulate(brain *Brain) {
	for neuronID, value := range n.In {
		_ = brain.GetNeuron(neuronID)
		n.Charge += value
		n.In[neuronID] = 0
	}
	if n.Charge > FiringLimit {
		brain.QueueForFiring(n)
	}
}

func (n *Neuron) Distribute(b *Brain) {
	for neuronID, value := range n.Out {
		_ = b.GetNeuron(neuronID)
		n.Charge += value
		n.In[neuronID] = 0
	}
}
