package wapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"mime/multipart"

	"github.com/rauldotgit/pricemage/utils"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Sub         string `json:"sub"`
}

func makeMultiGrantBuffer() (bytes.Buffer, error) {
	// writing multi part request body (multi part is specific format)
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	writer.SetBoundary("----boundary")
	err := writer.WriteField("grant_type", "client_credentials")
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error in makeMultiGrantBuffer: %w", err)
	}

	closeErr := writer.Close()
	if closeErr != nil {
		return bytes.Buffer{}, fmt.Errorf("error in makeMultiGrantBuffer: %w", closeErr)
	}

	return b, nil
}

// TODO: implement proper error handling
func GetAuth(client *http.Client) (AuthResponse, error) {

	multiGrantBuffer, err := makeMultiGrantBuffer()
	if err != nil {
		return AuthResponse{}, err
	}

	reqBody := bytes.NewReader(multiGrantBuffer.Bytes())

	resJSON, err := utils.Fetch("https://oauth.battle.net/token", "POST", reqBody, client, nil)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error in GetAuth: %w", err)
	}

	var res AuthResponse

	marshErr := json.Unmarshal(resJSON, &res)
	if marshErr != nil {
		return AuthResponse{}, fmt.Errorf("error in GetAuth: %w", marshErr)
	}

	return res, nil
}
