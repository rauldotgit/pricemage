package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendErrorJSON(rw http.ResponseWriter, status int, errMsg string) {
	if status > 499 {
		log.Println("Responding with 5XX error:", errMsg)
	}

	type errResponse struct {
		// tells the marshal / unmarshal how to process in json
		Error string `json:"error"`
	}

	// send response with error msg 
	SendResponseJSON(rw, status, errResponse{
		Error: errMsg,
	})
}

func SendResponseJSON(rw http.ResponseWriter, status int, payload interface{}) {
	// will attempt to load our data to json
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error in responseJSON: failed to marshal: %v", payload)
	}

	// add content type, and status to header
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(200)

	// write the response
	rw.Write(data)
}
