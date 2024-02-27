package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type TSMAuthRequest struct {
	ClientID  string `json:"client_id"`
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`
	Token     string `json:"token"`
}

type TSMAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func fetchTSMAuth(client *http.Client) TSMAuthResponse {
	url := "https://auth.tradeskillmaster.com/oauth2/token"

	// create request body
	reqData := TSMAuthRequest{
		ClientID:  "c260f00d-1071-409a-992f-dda2e5498536",
		GrantType: "api_token",
		Scope:     "app:realm-api app:pricing-api",
		Token:     os.Getenv("TSM_API_KEY"),
	}

	// marshal the req body
	reqJSON, marshErr := json.Marshal(reqData)
	if marshErr != nil {
		log.Println(marshErr)
		return TSMAuthResponse{}
	}

	// create an io reader for json
	reqBody := bytes.NewReader(reqJSON)

	resJSON, fetchErr := fetch(url, "POST", reqBody, client)
	if fetchErr != nil {
		log.Println(fetchErr)
		return TSMAuthResponse{}
	}

	// create unmarshalling struct
	var authRes TSMAuthResponse

	// unmarshal into struct
	err := json.Unmarshal(resJSON, &authRes)
	if err != nil {
		log.Println(err)
		return TSMAuthResponse{}
	}

	return authRes
}

func refreshTSMAuth() {
	return
}
