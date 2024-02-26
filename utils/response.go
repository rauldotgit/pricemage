package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseJSON(rw http.ResponseWriter, status int, payload interface{}) {
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
