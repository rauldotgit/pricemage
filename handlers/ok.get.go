package handlers

import (
	"net/http"

	"github.com/rauldotgit/wowauction/utils"
)

// classic handler structure with http.ResponseWriter and a pointer to the request
func HandleOK(rw http.ResponseWriter, r *http.Request) {
	utils.SendResponseJSON(rw, 200, struct{}{})
}
