package main

import (
	"fmt"
	"github.com/DeepThought7777/MassAI/codebase"
	"io"
	"net/http"
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
	testEndpoint(codebase.BuildRegisterURL(baseURL, registeredEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(codebase.BuildRegisterURL(baseURL, connectedEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(codebase.BuildConnectURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityConnected)

	fmt.Println("\n>>> TEST REGISTER / UNREGISTER")
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesNonDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesNegativeDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesZeroDigit, bytesValidLength), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesNonDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesNegativeDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesZeroDigit), http.StatusConflict, codebase.ErrorValueNotPositive)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(codebase.BuildUnregisterURL(baseURL, unknownEntityID), http.StatusOK, codebase.InfoEntityUnregistered)
	testEndpoint(codebase.BuildRegisterURL(baseURL, unknownEntityID, bytesValidLength, bytesValidLength), http.StatusOK, codebase.InfoEntityRegistered)
	testEndpoint(codebase.BuildUnregisterURL(baseURL, unknownEntityID), http.StatusOK, codebase.InfoEntityUnregistered)

	fmt.Println("\n>>> TEST CONNECT / DISCONNECT")
	testEndpoint(codebase.BuildConnectURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityConnected)
	testEndpoint(codebase.BuildDisconnectURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityDisconnected)
	testEndpoint(codebase.BuildDisconnectURL(baseURL, registeredEntityID), http.StatusConflict, codebase.ErrorEntityNotConnected)
	testEndpoint(codebase.BuildConnectURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorEntityAlreadyConnected)

	fmt.Println("\n>>> TEST SEND_INPUTS / GET_OUTPUTS")
	testEndpoint(codebase.BuildSendInputsURL(baseURL, unknownEntityID, bytesValidLength), http.StatusNotFound, codebase.ErrorNotRegistered)
	testEndpoint(codebase.BuildGetOutputsURL(baseURL, unknownEntityID), http.StatusNotFound, codebase.ErrorNotRegistered)
	testEndpoint(codebase.BuildSendInputsURL(baseURL, registeredEntityID, bytesValidLength), http.StatusConflict, codebase.ErrorEntityNotConnected)
	testEndpoint(codebase.BuildGetOutputsURL(baseURL, registeredEntityID), http.StatusConflict, codebase.ErrorEntityNotConnected)

	testEndpoint(codebase.BuildSendInputsURL(baseURL, connectedEntityID, bytesTooShortLength), http.StatusConflict, codebase.ErrorLengthInvalid)
	testEndpoint(codebase.BuildGetOutputsURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorLengthInvalid)

	testEndpoint(codebase.BuildSendInputsURL(baseURL, connectedEntityID, bytesTooLongLength), http.StatusConflict, codebase.ErrorLengthInvalid)
	testEndpoint(codebase.BuildGetOutputsURL(baseURL, connectedEntityID), http.StatusConflict, codebase.ErrorLengthInvalid)

	testEndpoint(codebase.BuildSendInputsURL(baseURL, connectedEntityID, bytesValidLength), http.StatusOK, codebase.InfoInputDataSent)
	testEndpoint(codebase.BuildGetOutputsURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoOutputDataValid)

	fmt.Println("\n>>> TEARDOWN")
	testEndpoint(codebase.BuildDisconnectURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityDisconnected)
	testEndpoint(codebase.BuildUnregisterURL(baseURL, connectedEntityID), http.StatusOK, codebase.InfoEntityUnregistered)
	testEndpoint(codebase.BuildUnregisterURL(baseURL, registeredEntityID), http.StatusOK, codebase.InfoEntityUnregistered)

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
