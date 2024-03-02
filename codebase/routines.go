package codebase

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func DisplayAndOptionallyExit(errorMessage string, exit bool) {
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

func RandomBytes(numBytes int) ([]byte, error) {
	if numBytes < 0 {
		return nil, fmt.Errorf("numBytes cannot be negative")
	}
	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ByteSliceToBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64ToByteSlice(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

func WriteToBody(w http.ResponseWriter, message string) {
	_, err := w.Write([]byte(message))
	if err != nil {
		fmt.Printf("ERROR: Couldn't write Body: %s\n", message)
	}
}

func BuildRegisterURL(baseURL, entityID, bytesIn, bytesOut string) string {
	return fmt.Sprintf("%s/register?entityId=%s&bytesInput=%s&bytesOutput=%s", baseURL, entityID, bytesIn, bytesOut)
}

func BuildUnregisterURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/unregister?entityId=%s&", baseURL, entityID)
}

func BuildConnectURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/connect?entityId=%s&", baseURL, entityID)
}

func BuildDisconnectURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/disconnect?entityId=%s&", baseURL, entityID)
}

func BuildSendInputsURL(baseURL, entityID, stringLength string) string {
	length, err := strconv.Atoi(stringLength)
	if err != nil {
		fmt.Println(">>> LENGTH STRING INVALID")
		return fmt.Sprintf("%s/send_inputs?entityId=%s", baseURL, entityID)
	}

	inputsByteSlice, err := RandomBytes(length)
	if err != nil {
		fmt.Println(">>> CANNOT GENERATE RANDOM BYTES")
		inputsByteSlice = []byte("")
	}

	inputsBase64 := ByteSliceToBase64URL(inputsByteSlice)
	return fmt.Sprintf("%s/send_inputs?entityId=%s&inputsBase64=%s", baseURL, entityID, inputsBase64)
}

func BuildGetOutputsURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/get_outputs?entityId=%s", baseURL, entityID)
}

func SendRequest(url string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {

	}
	fmt.Printf("OK: message body [%s]\n", string(body))
}
