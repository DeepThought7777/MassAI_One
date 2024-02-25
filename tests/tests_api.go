package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/DeepThought7777/MassAI/codebase"
)

const (
	baseURL             = "http://127.0.0.1:8080/v1"
	bytesNonDigit       = "ABC"
	bytesNegativeDigit  = "-1"
	bytesZeroDigit      = "0"
	bytesValidLength    = "8"
	bytesTooShortLength = "7"
	bytesTooLongLength  = "9"
)

func main() {
	unknownEntityID := "UNKNOWN_ID"
	registeredEntityID := "REGISTERED_ID"
	connectedEntityID := "CONNECTED_ID"

	fmt.Println("\n>>> SETUP")
	testEndpoint(buildRegisterURL(baseURL, registeredEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(buildRegisterURL(baseURL, connectedEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(buildConnectURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityConnected)

	fmt.Println("\n>>> TEST REGISTER / UNREGISTER")
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesNonDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesNegativeDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesZeroDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesNonDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesNegativeDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesZeroDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(buildUnregisterURL(baseURL, unknownEntityID), http.StatusOK, codebase.InfoEntityUnregistered)
	testEndpoint(buildUnregisterURL(baseURL, unknownEntityID), http.StatusConflict, codebase.ErrorNotRegistered)
	testEndpoint(buildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(buildUnregisterURL(baseURL, unknownEntityID), http.StatusOK, codebase.InfoEntityUnregistered)

	fmt.Println("\n>>> TEST CONNECT / DISCONNECT")
	testEndpoint(buildConnectURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityConnected)
	testEndpoint(buildDisconnectURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityDisconnected)
	testEndpoint(buildDisconnectURL(baseURL, registeredEntityID), http.StatusConflict, codebase.ErrorEntityNotConnected)
	testEndpoint(buildConnectURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorEntityAlreadyConnected)

	fmt.Println("\n>>> TEST SEND_INPUTS / GET_OUTPUTS")
	testEndpoint(buildSendInputsURL(baseURL, unknownEntityID, bytesValidLength), http.StatusNotFound, codebase.ErrorNotRegistered)
	testEndpoint(buildGetOutputsURL(baseURL, unknownEntityID), http.StatusNotFound, codebase.ErrorNotRegistered)
	testEndpoint(buildSendInputsURL(baseURL, registeredEntityID, bytesValidLength), http.StatusConflict, codebase.ErrorEntityNotConnected)
	testEndpoint(buildGetOutputsURL(baseURL, registeredEntityID), http.StatusConflict, codebase.ErrorEntityNotConnected)

	testEndpoint(buildSendInputsURL(baseURL, connectedEntityID, bytesTooShortLength), http.StatusConflict, codebase.ErrorLengthInvalid)
	testEndpoint(buildGetOutputsURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorLengthInvalid)

	testEndpoint(buildSendInputsURL(baseURL, connectedEntityID, bytesTooLongLength), http.StatusConflict, codebase.ErrorLengthInvalid)
	testEndpoint(buildGetOutputsURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorLengthInvalid)

	testEndpoint(buildSendInputsURL(baseURL, connectedEntityID, bytesValidLength), http.StatusOK, codebase.InfoInputDataSent)
	testEndpoint(buildGetOutputsURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoOutputDataValid)

	fmt.Println("\n>>> TEARDOWN")
	testEndpoint(buildDisconnectURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityDisconnected)
	testEndpoint(buildUnregisterURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityUnregistered)
	testEndpoint(buildUnregisterURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityUnregistered)

	codebase.DisplayAndOptionallyExit("Press Enter to exit...", true)
}

func testEndpoint(url string, statusCode int, message string) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	ok := true
	fmt.Printf("Endpoint: %s\n", url)
	if response.StatusCode != statusCode {
		fmt.Printf("ERROR: statuscode incorrect, was [%d], should be [%d]\n", response.StatusCode, statusCode)
		ok = false
	}

	if string(body) != message {
		fmt.Printf("ERROR: message body incorrect, was [%s], should be [%s]\n", string(body), message)
		ok = false
	}
	if ok {
		fmt.Printf("OK: message body [%s]\n", string(body))
	}
}

func buildRegisterURL(baseURL, entityID, bytesIn, bytesOut string) string {
	return fmt.Sprintf("%s/register?entityId=%s&bytesInput=%s&bytesOutput=%s", baseURL, entityID, bytesIn, bytesOut)
}

func buildUnregisterURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/unregister?entityId=%s&", baseURL, entityID)
}

func buildConnectURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/connect?entityId=%s&", baseURL, entityID)
}

func buildDisconnectURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/disconnect?entityId=%s&", baseURL, entityID)
}

func buildSendInputsURL(baseURL, entityID, stringLength string) string {
	length, err := strconv.Atoi(stringLength)
	if err != nil {
		fmt.Println(">>> LENGTH STRING INVALID")
		return fmt.Sprintf("%s/send_inputs?entityId=%s", baseURL, entityID)
	}

	inputsByteSlice, err := codebase.RandomBytes(length)
	if err != nil {
		fmt.Println(">>> CANNOT GENERATE RANDOM BYTES")
		inputsByteSlice = []byte("")
	}

	inputsBase64 := codebase.ByteSliceToBase64URL(inputsByteSlice)
	return fmt.Sprintf("%s/send_inputs?entityId=%s&inputsBase64=%s", baseURL, entityID, inputsBase64)
}

func buildGetOutputsURL(baseURL, entityID string) string {
	return fmt.Sprintf("%s/get_outputs?entityId=%s", baseURL, entityID)
}
