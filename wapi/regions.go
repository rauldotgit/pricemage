package wapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/rauldotgit/pricemage/utils"
)

type Region struct {
	RegionID     int
	Name         string
	Prefix       string
	GMTOffset    int
	GameVersion  string
	LastModified int64
}

type Metadata struct {
	TotalItems int
}

type WOWRegions struct {
	Items    []Region
	Metadata Metadata
}

func GetRegions(client *http.Client, authToken *string) (WOWRegions, error) {
	url := "https://realm-api.tradeskillmaster.com/regions"

	if authToken == nil {
		return WOWRegions{}, errors.New("error in GetRegions: No access token specified")
	}

	res, err := utils.Fetch(url, "GET", nil, client, authToken)
	if err != nil {
		return WOWRegions{}, fmt.Errorf("error in GetRegions: %w", err)
	}

	var regions WOWRegions

	umErr := json.Unmarshal(res, &regions)
	if umErr != nil {
		log.Println("Error in region unmarshal")
		return WOWRegions{}, fmt.Errorf("error in GetRegions: %w", umErr)
	}

	return regions, nil
}
