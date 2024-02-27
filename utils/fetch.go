package utils

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

func fetch(url string, method string, reqBody io.Reader, client *http.Client) ([]byte, error) {

	var reqMethod string
	errMsg := "Error in FetchTSMAuth: "

	switch method {
	case http.MethodGet:
		reqMethod = http.MethodGet
	case http.MethodPost:
		reqMethod = http.MethodPost
	case http.MethodDelete:
		reqMethod = http.MethodDelete
	case http.MethodPut:
		reqMethod = http.MethodPut
	case http.MethodOptions:
		reqMethod = http.MethodOptions
	default:
		return []byte{}, errors.New(errMsg + "invalid method string.")
	}

	req, err := http.NewRequest(reqMethod, url, reqBody)
	if err != nil {
		return []byte{}, errors.New(errMsg + "failed to create request.")
	}

	// set request header
	req.Header.Set("Content-Type", "application/json")

	// make request
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, errors.New(errMsg + "failed to make request.")
	}

	// close response body i.e. signal to server finish
	defer res.Body.Close()

	log.Println(res)

	// check status code
	if res.StatusCode > 299 {
		return []byte{}, errors.New(errMsg + "unexpected status code - " + strconv.Itoa(res.StatusCode))
	}

	// read response body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, errors.New(errMsg + "failed to read response body.")
	}

	return resBody, nil
}
