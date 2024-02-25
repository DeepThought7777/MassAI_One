package api

import (
	"net/http"
	"strconv"

	"github.com/DeepThought7777/MassAI/codebase"
	"github.com/DeepThought7777/MassAI/mind"
)

var entities = make(map[string]*mind.Entity)

// RegisterEntity handles the registering of an entity to the MassAI mind.
func RegisterEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")
	bytesInput := params.Get("bytesInput")
	bytesOutput := params.Get("bytesOutput")

	if _, exists := entities[entityID]; exists {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorAlreadyRegistered)
		return
	}

	inputs, err := strconv.Atoi(bytesInput)
	if err != nil || inputs < 1 {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorValueNotPositive)
		return
	}

	outputs, err := strconv.Atoi(bytesOutput)
	if err != nil || outputs < 1 {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorValueNotPositive)
		return
	}

	entities[entityID] = mind.NewEntity(entityID, inputs, outputs)
	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoEntityRegistered)
}

// UnregisterEntity handles the unregistering of an entity to the MassAI mind.
// This is akin to amputating a limb, because unregistering just chops off the nerves for inputs and outputs,
// and leaves the Neurons it was connected to intact.
func UnregisterEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	if _, exists := entities[entityID]; !exists {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorNotRegistered)
		return
	}

	delete(entities, entityID)
	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoEntityUnregistered)
}

func ConnectEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		codebase.WriteToBody(w, codebase.ErrorNotRegistered)
		return
	}

	if entity.IsConnected {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorEntityAlreadyConnected)
		return
	}

	entity.IsConnected = true

	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoEntityConnected)
}

func DisconnectEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		codebase.WriteToBody(w, codebase.ErrorNotRegistered)
		return
	}

	if !entity.IsConnected {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorEntityNotConnected)
		return
	}

	entity.IsConnected = false

	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoEntityDisconnected)
}

func SendInputs(w http.ResponseWriter, r *http.Request) {
	entityID := r.URL.Query().Get("entityId")
	inputsBase64 := r.URL.Query().Get("inputsBase64")

	if entityID == "" || inputsBase64 == "" {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorEntityNotConnected)
		return
	}

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		codebase.WriteToBody(w, codebase.ErrorNotRegistered)
		return
	}

	if !entity.IsConnected {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorEntityNotConnected)
		return
	}

	inputs, err := codebase.Base64ToByteSlice(inputsBase64)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorCodingFailed)
		return
	}

	if len(inputs) != entity.BytesInput {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorLengthInvalid)
		return
	}

	entity.Base64Input = inputsBase64
	// also copy to output for now, to test...
	entity.Base64Output = inputsBase64

	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoInputDataSent)
}

func GetOutputs(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		codebase.WriteToBody(w, codebase.ErrorNotRegistered)
		return
	}

	if !entity.IsConnected {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorEntityNotConnected)
		return
	}

	outputs, err := codebase.Base64ToByteSlice(entity.Base64Output)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorCodingFailed)
		return
	}

	if len(outputs) != entity.BytesOutput {
		w.WriteHeader(http.StatusConflict)
		codebase.WriteToBody(w, codebase.ErrorLengthInvalid)
		return
	}

	w.WriteHeader(http.StatusOK)
	codebase.WriteToBody(w, codebase.InfoOutputDataValid)
}
