package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DeepThought7777/MassAI/codebase"
	"github.com/DeepThought7777/MassAI/mind"
)

var entities = make(map[string]*mind.Entity)

// registerEntity handles the registering of an entity to the MassAI mind.
func RegisterEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")
	bytesInput := params.Get("bytesInput")
	bytesOutput := params.Get("bytesOutput")

	if _, exists := entities[entityID]; exists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(codebase.ERR_ALREADY_REGISTERED))
		return
	}

	inputs, err := strconv.Atoi(bytesInput)
	if err != nil || inputs < 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(codebase.ERR_INPUTCOUNT_NOT_VALID))
		return
	}

	outputs, err := strconv.Atoi(bytesOutput)
	if err != nil || inputs < 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(codebase.ERR_OUTPUTCOUNT_NOT_VALID))
		return
	}

	entities[entityID] = mind.NewEntity(entityID, inputs, outputs)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(codebase.INFO_ENTITY_REGISTERED))
}

func ConnectEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found - Entity not registered")
		return
	}

	entity.IsConnected = true

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func DisconnectEntity(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found - Entity not registered")
		return
	}

	if !entity.IsConnected {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, "Conflict - Entity not connected")
		return
	}

	entity.IsConnected = false

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func SendInputs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	entityID := r.URL.Query().Get("entityId")
	inputsBase64 := r.URL.Query().Get("inputsBase64")

	// Validate required parameters
	if entityID == "" || inputsBase64 == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Decode BASE64-encoded input bytes
	inputs, err := base64.StdEncoding.DecodeString(inputsBase64)
	if err != nil {
		http.Error(w, "Invalid BASE64 encoding", http.StatusBadRequest)
		return
	}

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found - Entity not registered")
		return
	}

	// Check BASE64 string length requirements
	if len(inputs) != entity.BytesInput {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprint(w, "Precondition Failed - BASE64 string length requirements not met")
		return
	}

	entity.Base64Input = inputsBase64

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func GetOutputs(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	entityID := params.Get("entityId")

	entity, exists := entities[entityID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found - Entity not registered")
		return
	}

	if !entity.IsConnected {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprint(w, "Precondition Failed - Entity not connected")
		return
	}

	// Check BASE64 string length requirements
	if len(entity.Base64Output) != entity.BytesOutput {
		w.WriteHeader(http.StatusPreconditionFailed)
		fmt.Fprint(w, "Precondition Failed - BASE64 string length requirements not met")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, entity.Base64Output)
}
