package utils

import (
	"log"
	"net/http"
)

func LocalExec() {

	client := &http.Client{}
	var tsmAuth TSMAuthResponse = fetchTSMAuth(client)

	log.Println("TSM auth response:")
	log.Println(tsmAuth.AccessToken)
}
