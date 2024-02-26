package handlers

import (
	"net/http"
)

// classic handler structure with http.ResponseWriter and a pointer to the request
func handleOk(rw http.ResponseWriter, r *http.Request) {
	return
}
