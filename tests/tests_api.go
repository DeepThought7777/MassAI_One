package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/DeepThought7777/MassAI/codebase"
)

const baseURL = "http://127.0.0.1:8080/v1"

func main() {
	entityID := "eeb344f4-32ad-42a1-b678-59b1f6caf185"
	//base64String := "SGVsbG8gd29ybGQ="

	// Test /register endpoint
	testEndpoint(buildGetURL(baseURL, "register", entityID, "ABC", "512"), 409, codebase.ERR_INPUTCOUNT_NOT_VALID)
	testEndpoint(buildGetURL(baseURL, "register", entityID, "256", "ABC"), 409, codebase.ERR_OUTPUTCOUNT_NOT_VALID)
	testEndpoint(buildGetURL(baseURL, "register", entityID, "256", "512"), 200, codebase.INFO_ENTITY_REGISTERED)
	testEndpoint(buildGetURL(baseURL, "register", entityID, "256", "512"), 409, codebase.ERR_ALREADY_REGISTERED)
	testEndpoint(buildGetURL(baseURL, "register", entityID, "256", "512"), 409, codebase.ERR_ALREADY_REGISTERED)

	// Test /connect endpoint
	//connectURL := fmt.Sprintf("%s/connect?entityId=%s", baseURL, entityID)
	//testEndpoint(connectURL)

	// Test /send-inputs endpoint
	//sendInputsURL := fmt.Sprintf("%s/send-inputs?entityId=%s&inputsBase64=%s", baseURL, entityID, base64String)
	//testEndpoint(sendInputsURL)

	// Test /get-outputs endpoint
	//getOutputsURL := fmt.Sprintf("%s/get-outputs?entityId=%s", baseURL, entityID)
	//testEndpoint(getOutputsURL)

	// Test /disconnect endpoint
	//disconnectURL := fmt.Sprintf("%s/disconnect?entityId=%s", baseURL, entityID)
	//testEndpoint(disconnectURL)
}

func buildGetURL(baseURL, endpoint, entityID, bytesIn, bytesOut string) string {
	return fmt.Sprintf("%s/%s?entityId=%s&bytesInput=%s&bytesOutput=%s",
		baseURL, endpoint, entityID, bytesIn, bytesOut)
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
	fmt.Printf("Endpoint: %s  ", url)
	if response.StatusCode != statusCode {
		fmt.Println("ERROR: statuscode incorrect")
		ok = false
	}

	if string(body) != message {
		fmt.Println("ERROR: message in body incorrect")
		ok = false
	}
	if ok {
		fmt.Println("OK")
	}
}

/*
func testSendInputsEndpoint(url string) {
	// Example JSON payload
	payload := map[string]string{"base64Input": "SGVsbG8gd29ybGQ="}
	jsonPayload, _ := json.Marshal(payload)

	response, err := http.Put(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	fmt.Printf("Endpoint: %s\nStatus Code: %d\nResponse: %s\n\n", url, response.StatusCode, body)
}
*/
