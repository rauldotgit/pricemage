package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Fetch(url string, method string, reqBody io.Reader, client *http.Client, authToken *string) ([]byte, error) {

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return []byte{}, fmt.Errorf("error in fetch: %w", err)
	}

	if authToken != nil {
		req.Header.Set("Authorization", "Bearer "+*authToken)
	} else {
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", "----boundary"))
		req.SetBasicAuth(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	}

	// make request
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error in Fetch: %w", err)
	}

	// close response body i.e. signal to server finish
	defer res.Body.Close()

	// check status code
	if res.StatusCode > 299 {
		log.Println(res)
		return []byte{}, errors.New("error in Fetch: unexpected status code - " + strconv.Itoa(res.StatusCode))
	}

	// read response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error in Fetch: %w", err)
	}

	return resBody, nil
}
